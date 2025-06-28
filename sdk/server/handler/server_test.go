package handler

import (
	"bytes"
	"context"
	"encoding/json"
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
				Result:  &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true, Status: types.TaskStatus{State: types.COMPLETED}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := NewExecutor()
			manger := QueueManger{}
			handler := NewDefaultHandler(store, executor, WithQueueManger(manger))
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
		})
	}
}
