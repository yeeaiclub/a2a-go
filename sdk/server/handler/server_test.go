// Copyright 2025 yumosx
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

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/types"
)

var mockAgentCard = types.AgentCard{
	Name:        "test agent",
	Description: "a test agent for test",
	URL:         "http://localhost:41241/",
	Version:     "1.0",
	Capabilities: &types.AgentCapabilities{
		Streaming:              true,
		StateTransitionHistory: true,
	},
	Skills: []types.AgentSkill{
		{
			Id:          "test-skills",
			Name:        "test skill",
			Description: "a test skill for unit test",
			InputModes:  []string{"text/plain"},
		},
	},
}

func TestHandleMessageSend(t *testing.T) {
	testcases := []struct {
		name   string
		before func(store tasks.TaskStore)
		params types.MessageSendParam
		want   types.Task
	}{
		{
			name: "test message send",
			params: types.MessageSendParam{
				Message: &types.Message{
					TaskID: "1",
					Role:   types.User,
				},
			},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: types.Task{Id: "1", ContextId: "2", History: []*types.Message{{TaskID: "1", ContextID: "2"}}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := newExecutor()
			handler := NewDefaultHandler(store, executor, WithQueueManger(QueueManger{}))
			server := NewServer(mockAgentCard, handler, "/")
			request := types.JSONRPCRequest{
				Id:     "1",
				Method: types.MethodMessageSend,
				Params: tc.params,
			}
			req, err := json.Marshal(request)
			require.NoError(t, err)
			newReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(req))
			w := httptest.NewRecorder()
			server.ServeHTTP(w, newReq)

			var resp types.JSONRPCResponse
			err = json.NewDecoder(w.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, resp.JSONRPC, types.Version)
			task, err := types.MapTo[types.Task](resp.Result)
			require.NoError(t, err)
			assert.Equal(t, task.Id, tc.want.Id)
		})
	}
}

func TestHandleMessageSendStream(t *testing.T) {
	testcases := []struct {
		name   string
		params types.MessageSendParam
		before func(store tasks.TaskStore)
		want   types.JSONRPCResponse
	}{
		{
			name: "test message send stream",
			params: types.MessageSendParam{
				Message: &types.Message{
					TaskID: "1",
					Role:   types.User,
				},
			},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: types.JSONRPCResponse{
				Id:      "1",
				JSONRPC: types.Version,
				Result:  types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true, Status: types.TaskStatus{State: types.COMPLETED}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := newExecutor()
			handler := NewDefaultHandler(store, executor, WithQueueManger(QueueManger{}))
			server := NewServer(mockAgentCard, handler, "/")
			request := types.JSONRPCRequest{
				Id:     "1",
				Method: types.MethodMessageStream,
				Params: tc.params,
			}
			req, err := json.Marshal(request)
			require.NoError(t, err)
			newReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(req))
			newReq.Header.Set("Content-Type", "application/json")
			newReq.Header.Set("Accept", "text/event-stream")
			w := httptest.NewRecorder()
			server.ServeHTTP(w, newReq)
			response := strings.Split(strings.TrimSpace(w.Body.String()), "\n")
			assert.True(t, len(response) == 1)
			var resp types.JSONRPCResponse

			err = json.Unmarshal([]byte(response[0]), &resp)
			require.NoError(t, err)
			assert.Equal(t, resp.Id, tc.want.Id)
			assert.Equal(t, resp.JSONRPC, tc.want.JSONRPC)

			value, err := types.MapTo[types.TaskStatusUpdateEvent](resp.Result)
			require.NoError(t, err)
			assert.NotEmpty(t, value.Status.TimeStamp)
			value.Status.TimeStamp = ""
			assert.Equal(t, value, tc.want.Result)
		})
	}
}

func TestHandleGetTask(t *testing.T) {
	testcases := []struct {
		name   string
		params types.TaskQueryParams
		before func(store tasks.TaskStore)
		want   types.Task
	}{
		{
			name: "test get task",
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			params: types.TaskQueryParams{
				Id: "1",
			},
			want: types.Task{
				Id:        "1",
				ContextId: "2",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := newExecutor()
			handler := NewDefaultHandler(store, executor, WithQueueManger(QueueManger{}))
			server := NewServer(mockAgentCard, handler, "/")

			request := types.JSONRPCRequest{
				Id:     "1",
				Method: types.MethodTasksGet,
				Params: tc.params,
			}
			body, err := json.Marshal(request)
			require.NoError(t, err)
			newReq := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, newReq)

			var resp types.JSONRPCResponse
			err = json.NewDecoder(recorder.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, resp.JSONRPC, types.Version)
			task, err := types.MapTo[types.Task](resp.Result)
			require.NoError(t, err)
			assert.Equal(t, task, tc.want)
		})
	}
}

func TestGetCard(t *testing.T) {
	testcases := []struct {
		name string
		card types.AgentCard
	}{
		{
			name: "get card",
			card: mockAgentCard,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			executor := newExecutor()
			handler := NewDefaultHandler(store, executor, WithQueueManger(QueueManger{}))

			server := NewServer(mockAgentCard, handler, "/", WithAgentCardPath("/card"))

			req := httptest.NewRequest("GET", "/card", nil)
			w := httptest.NewRecorder()

			server.handleGetAgentCard(w, req)

			body := w.Body.String()
			require.Equal(t, http.StatusOK, w.Code)
			var card types.AgentCard
			err := json.Unmarshal([]byte(body), &card)
			require.NoError(t, err)
			assert.Equal(t, card, tc.card)
		})
	}
}
