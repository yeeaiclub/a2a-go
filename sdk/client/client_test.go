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
	"github.com/yeeaiclub/a2a-go/sdk/client/middleware"
	"github.com/yeeaiclub/a2a-go/sdk/types"
	"github.com/yeeaiclub/a2a-go/sdk/web"
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
			task: types.Task{Id: "123", Artifacts: []types.Artifact{}},
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
			client := NewClient(http.DefaultClient, server.URL)
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
			want:   types.Task{Id: "123", Artifacts: []types.Artifact{}},
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

			client := NewClient(http.DefaultClient, server.URL)
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
			want:   types.Task{Id: "123", Artifacts: []types.Artifact{}},
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

			client := NewClient(http.DefaultClient, server.URL)
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
				Message: &types.Message{Role: types.User, Kind: types.EventTypeMessage},
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
						Id:   "123",
						Kind: types.EventTypeTask,
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

			client := NewClient(http.DefaultClient, server.URL)
			eventChan := make(chan types.Event)
			errChan := make(chan error, 1)

			go func() {
				errChan <- client.SendMessageStream(tc.params, eventChan)
				close(eventChan)
			}()

			var events []*types.Task
			for event := range eventChan {
				task, ok := event.(*types.Task)
				require.True(t, ok)
				events = append(events, task)
			}
			err := <-errChan
			require.NoError(t, err)
			<-done
			assert.Equal(t, string(events[0].Status.State), string(types.COMPLETED))
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name            string
		setupClient     func() *A2AClient
		setupContext    func() web.Context
		expectedHeaders map[string]string
		expectedError   bool
	}{
		{
			name: "no middleware",
			setupClient: func() *A2AClient {
				return NewClient(&http.Client{}, "http://example.com")
			},
			setupContext: func() web.Context {
				ctx := middleware.NewCallContext(1)
				req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				ctx.SetRequest(req)
				return ctx
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "with authentication middleware",
			setupClient: func() *A2AClient {
				client := NewClient(&http.Client{}, "http://example.com")
				credential := middleware.NewInMemoryContextCredentials()
				credential.SetCredentials("session1", "apiKey", "api-token-123")

				client.Use(middleware.Intercept(credential))
				client.card = &types.AgentCard{
					Security: types.SecurityRequirement{
						{"apiKey": []string{}},
					},
					SecuritySchemes: map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					},
				}
				return client
			},
			setupContext: func() web.Context {
				ctx := middleware.NewCallContext(1)
				ctx.SetSecurityConfig(
					types.SecurityRequirement{
						{"apiKey": []string{}},
					},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					})
				ctx.Set("sessionId", "session1")
				req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				ctx.SetRequest(req)
				return ctx
			},
			expectedHeaders: map[string]string{
				"X-API-Key": "api-token-123",
			},
			expectedError: false,
		},
		{
			name: "middleware error",
			setupClient: func() *A2AClient {
				client := NewClient(&http.Client{}, "http://example.com")
				errorMiddleware := func(next web.HandlerFunc) web.HandlerFunc {
					return func(ctx web.Context) error {
						return assert.AnError
					}
				}
				client.Use(errorMiddleware)

				return client
			},
			setupContext: func() web.Context {
				ctx := middleware.NewCallContext(1)
				req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				ctx.SetRequest(req)
				return ctx
			},
			expectedHeaders: map[string]string{},
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupClient()
			ctx := tt.setupContext()

			err := client.apply(ctx)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			for expectedHeader, expectedValue := range tt.expectedHeaders {
				actualValue := ctx.Request().Header.Get(expectedHeader)
				assert.Equal(t, expectedValue, actualValue,
					"Header %s should be set to %s, but got %s",
					expectedHeader, expectedValue, actualValue)
			}

			if len(tt.expectedHeaders) == 0 {
				assert.Empty(t, ctx.Request().Header)
			}
		})
	}
}

func TestA2AClient_Use(t *testing.T) {
	t.Run("test use", func(t *testing.T) {
		client := NewClient(&http.Client{}, "http://example.com")

		assert.Empty(t, client.middlewares)

		middleware1 := func(next web.HandlerFunc) web.HandlerFunc {
			return func(ctx web.Context) error {
				return next(ctx)
			}
		}
		client.Use(middleware1)
		assert.Len(t, client.middlewares, 1)

		middleware2 := func(next web.HandlerFunc) web.HandlerFunc {
			return func(ctx web.Context) error {
				return next(ctx)
			}
		}
		middleware3 := func(next web.HandlerFunc) web.HandlerFunc {
			return func(ctx web.Context) error {
				return next(ctx)
			}
		}
		client.Use(middleware2, middleware3)
		assert.Len(t, client.middlewares, 3)
	})
}
