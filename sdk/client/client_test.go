// Copyright 2025 yeeaiclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestSendMessage(t *testing.T) {
	testcases := []struct {
		name   string
		params types.MessageSendParam
		task   types.Task
	}{
		{
			name: "send message",
			params: types.MessageSendParam{
				Message: &types.Message{
					TaskID: "123",
				},
			},
			task: types.Task{Id: "123"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req types.JSONRPCRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, types.MethodMessageSend, req.Method)
				resp := types.JSONRPCResponse{
					JSONRPC: types.Version,
					Result:  types.Task{Id: "123"},
				}
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(resp)
				assert.NoError(t, err)
			}))
			defer server.Close()
			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			message, err := client.SendMessage(tc.params)
			require.NoError(t, err)
			task, err := types.MapTo[types.Task](message.Result)
			require.NoError(t, err)
			assert.Equal(t, task, tc.task)
		})
	}
}

func TestGetTask(t *testing.T) {
	testcases := []struct {
		name   string
		params types.TaskQueryParams
		want   types.Task
	}{
		{
			name:   "test get task",
			params: types.TaskQueryParams{Id: "1"},
			want:   types.Task{Id: "123"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				var req types.JSONRPCRequest
				err := json.NewDecoder(request.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, types.MethodTasksGet, req.Method)
				resp := types.JSONRPCResponse{
					JSONRPC: types.Version,
					Result:  types.Task{Id: "123"},
				}
				writer.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(writer).Encode(resp)
				assert.NoError(t, err)
			}))
			defer server.Close()

			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			resp, err := client.GetTask(tc.params)
			require.NoError(t, err)
			task, err := types.MapTo[types.Task](resp.Result)
			require.NoError(t, err)
			assert.Equal(t, tc.want, task)
		})
	}
}

func TestCancelTask(t *testing.T) {
	testcases := []struct {
		name   string
		params types.TaskIdParams
		want   types.Task
	}{
		{
			name:   "cancel task",
			params: types.TaskIdParams{Id: "123"},
			want:   types.Task{Id: "123"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				var req types.JSONRPCRequest
				err := json.NewDecoder(request.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, types.MethodTasksCancel, req.Method)
				resp := types.JSONRPCResponse{
					JSONRPC: types.Version,
					Result:  types.Task{Id: "123"},
				}
				writer.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(writer).Encode(resp)
				assert.NoError(t, err)
			}))
			defer server.Close()

			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			resp, err := client.CancelTask(tc.params)
			require.NoError(t, err)

			task, err := types.MapTo[types.Task](resp.Result)
			require.NoError(t, err)
			assert.Equal(t, tc.want, task)
		})
	}
}

func TestMessageStream(t *testing.T) {
	testcases := []struct {
		name   string
		params types.MessageSendParam
	}{
		{
			name: "message stream",
			params: types.MessageSendParam{
				Message: &types.Message{Role: types.User},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				var req types.JSONRPCRequest
				err := json.NewDecoder(request.Body).Decode(&req)
				assert.NoError(t, err)
				assert.Equal(t, types.MethodMessageStream, req.Method)

				writer.Header().Set("Content-Type", "application/json")
				writer.(http.Flusher).Flush()
				events := []*types.Task{
					{
						Id: "123",
						Status: types.TaskStatus{
							State: types.COMPLETED,
						},
					},
				}
				for _, event := range events {
					resp := types.JSONRPCResponse{
						JSONRPC: types.Version,
						Result:  event,
					}
					err = json.NewEncoder(writer).Encode(resp)
					assert.NoError(t, err)
					writer.(http.Flusher).Flush()
				}
				close(done)
			}))

			defer server.Close()

			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			eventChan := make(chan any)
			errChan := make(chan error, 1)

			go func() {
				errChan <- client.SendMessageStream(tc.params, eventChan)
				close(eventChan)
			}()

			var events []types.Task
			for event := range eventChan {
				rawMsg, ok := event.(json.RawMessage)
				require.True(t, ok)
				var task types.Task
				err = json.Unmarshal(rawMsg, &task)
				require.NoError(t, err)
				events = append(events, task)
			}
			err = <-errChan
			require.NoError(t, err)
			<-done
			assert.Equal(t, string(events[0].Status.State), string(types.COMPLETED))
		})
	}
}
