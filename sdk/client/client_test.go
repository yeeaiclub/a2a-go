// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/yumosx/a2a-go/sdk/types"
)

type ClientSuite struct {
	suite.Suite
}

func NewClientSuite() *ClientSuite {
	return &ClientSuite{}
}

func TestClient(t *testing.T) {
	suite.Run(t, NewClientSuite())
}

func (c *ClientSuite) TestSendMessage() {
	t := c.T()
	testcases := []struct {
		name string
		req  types.SendMessageRequest
		want types.SendMessageResponse
	}{
		{
			name: "send message",
			req: types.SendMessageRequest{
				Method: types.MessageSend,
			},
			want: types.SendMessageResponse{
				Id: "1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				err := json.NewEncoder(w).Encode(&tc.want)
				require.NoError(t, err)
			}))
			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			message, err := client.SendMessage(context.Background(), tc.req, nil)
			require.NoError(t, err)
			assert.Equal(t, message.Id, "1")
		})
	}
}

func (c *ClientSuite) TestGetTask() {
	t := c.T()
	testcases := []struct {
		name string
		req  types.GetTaskRequest
		want types.GetTaskSuccessResponse
	}{
		{
			name: "get task",
			req: types.GetTaskRequest{
				Id: "1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				err := json.NewEncoder(w).Encode(&tc.want)
				require.NoError(t, err)
			}))
			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			message, err := client.GetTask(context.Background(), tc.req, nil)
			require.NoError(t, err)
			assert.Equal(t, message, "1")
		})
	}
}

func (c *ClientSuite) CancelTask() {
	t := c.T()

	testcases := []struct {
		name string
		req  types.CancelTaskRequest
	}{
		{
			name: "cancel task",
			req:  types.CancelTaskRequest{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				err := json.NewEncoder(w).Encode(&tc.req)
				require.NoError(t, err)
			}))

			client, err := NewClient(http.DefaultClient, WithUrl(server.URL))
			require.NoError(t, err)
			message, err := client.CancelTask(context.Background(), tc.req, nil)
			require.NoError(t, err)
			assert.Equal(t, message, "1")
		})
	}
}
