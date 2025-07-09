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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryContextCredentials(t *testing.T) {
	cred := NewInMemoryContextCredentials()
	assert.NotNil(t, cred)
	assert.NotNil(t, cred.store)
	assert.Empty(t, cred.store)
}

func TestGetCredentials(t *testing.T) {
	tests := []struct {
		name               string
		setupCredentials   map[string]string // "sessionId" -> id schemeName -> token
		context            *CallContext
		securitySchemeName string
		expectedToken      string
		expectedError      string
	}{
		{
			name:               "nil context",
			setupCredentials:   map[string]string{"session1:scheme1": "token1"},
			context:            nil,
			securitySchemeName: "scheme1",
			expectedToken:      "",
			expectedError:      "",
		},
		{
			name:               "no sessionId in context",
			setupCredentials:   map[string]string{"session1:scheme1": "token1"},
			context:            NewCallContext(1),
			securitySchemeName: "scheme1",
			expectedToken:      "",
			expectedError:      "",
		},
		{
			name:             "invalid sessionId type",
			setupCredentials: map[string]string{"session1:scheme1": "token1"},
			context: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = 123
				return ctx
			}(),
			securitySchemeName: "scheme1",
			expectedToken:      "",
			expectedError:      "type asset failed",
		},
		{
			name:             "credentials not found",
			setupCredentials: map[string]string{"session1:scheme1": "token1"},
			context: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			securitySchemeName: "scheme2",
			expectedToken:      "",
			expectedError:      "",
		},
		{
			name:             "successful retrieval",
			setupCredentials: map[string]string{"session1:scheme1": "token1"},
			context: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			securitySchemeName: "scheme1",
			expectedToken:      "token1",
			expectedError:      "",
		},
		{
			name: "multiple credentials for different sessions",
			setupCredentials: map[string]string{
				"session1:scheme1": "token1",
				"session1:scheme2": "token2",
				"session2:scheme1": "token3",
			},
			context: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			securitySchemeName: "scheme2",
			expectedToken:      "token2",
			expectedError:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := NewInMemoryContextCredentials()

			for key, token := range tt.setupCredentials {
				parts := strings.Split(key, ":")
				if len(parts) == 2 {
					sessionId := parts[0]
					schemeName := parts[1]
					cred.SetCredentials(sessionId, schemeName, token)
				}
			}

			result, err := cred.GetCredentials(tt.securitySchemeName, tt.context)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedToken, result)
		})
	}
}

func TestSetCredentials(t *testing.T) {
	tests := []struct {
		name        string
		credentials []struct {
			sessionId          string
			securitySchemeName string
			token              string
		}
		expectedStoreSize int
	}{
		{
			name: "single credential",
			credentials: []struct {
				sessionId          string
				securitySchemeName string
				token              string
			}{
				{"session1", "scheme1", "token1"},
			},
			expectedStoreSize: 1,
		},
		{
			name: "multiple credentials for same session",
			credentials: []struct {
				sessionId          string
				securitySchemeName string
				token              string
			}{
				{"session1", "scheme1", "token1"},
				{"session1", "scheme2", "token2"},
			},
			expectedStoreSize: 2,
		},
		{
			name: "multiple credentials for different sessions",
			credentials: []struct {
				sessionId          string
				securitySchemeName string
				token              string
			}{
				{"session1", "scheme1", "token1"},
				{"session2", "scheme1", "token2"},
				{"session1", "scheme2", "token3"},
			},
			expectedStoreSize: 3,
		},
		{
			name: "update existing credential",
			credentials: []struct {
				sessionId          string
				securitySchemeName string
				token              string
			}{
				{"session1", "scheme1", "token1"},
				{"session1", "scheme1", "new_token1"},
			},
			expectedStoreSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := NewInMemoryContextCredentials()

			for _, c := range tt.credentials {
				cred.SetCredentials(c.sessionId, c.securitySchemeName, c.token)
			}

			assert.Len(t, cred.store, tt.expectedStoreSize)

			if len(tt.credentials) > 0 {
				lastCred := tt.credentials[len(tt.credentials)-1]
				context := NewCallContext(1)
				context.State["sessionId"] = lastCred.sessionId

				result, err := cred.GetCredentials(lastCred.securitySchemeName, context)
				require.NoError(t, err)
				assert.Equal(t, lastCred.token, result)
			}
		})
	}
}

func TestUpdateCredentials(t *testing.T) {
	tests := []struct {
		name         string
		initialToken string
		updatedToken string
		sessionId    string
		schemeName   string
	}{
		{
			name:         "update token",
			initialToken: "old_token",
			updatedToken: "new_token",
			sessionId:    "session1",
			schemeName:   "scheme1",
		},
		{
			name:         "update to empty token",
			initialToken: "old_token",
			updatedToken: "",
			sessionId:    "session1",
			schemeName:   "scheme1",
		},
		{
			name:         "update to same token",
			initialToken: "token1",
			updatedToken: "token1",
			sessionId:    "session1",
			schemeName:   "scheme1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := NewInMemoryContextCredentials()
			context := NewCallContext(1)
			context.State["sessionId"] = tt.sessionId

			cred.SetCredentials(tt.sessionId, tt.schemeName, tt.initialToken)

			result, err := cred.GetCredentials(tt.schemeName, context)
			require.NoError(t, err)
			assert.Equal(t, tt.initialToken, result)

			cred.SetCredentials(tt.sessionId, tt.schemeName, tt.updatedToken)

			result, err = cred.GetCredentials(tt.schemeName, context)
			require.NoError(t, err)
			assert.Equal(t, tt.updatedToken, result)
		})
	}
}
