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

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/types"
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
				require.NoError(t, err)
				assert.Equal(t, req.Method, types.MethodMessageSend)
				resp := types.JSONRPCResponse{
					JSONRPC: types.Version,
					Result:  types.Task{Id: "123"},
				}
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(resp)
				require.NoError(t, err)
			}))
			defer server.Close()
		})
	}
}
