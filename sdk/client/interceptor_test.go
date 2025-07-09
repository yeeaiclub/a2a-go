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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestNewInterceptor(t *testing.T) {
	credentials := NewInMemoryContextCredentials()
	interceptor := NewInterceptor(credentials)

	assert.NotNil(t, interceptor)
	assert.Equal(t, credentials, interceptor.CredentialService)
}

func TestIntercept(t *testing.T) {
	tests := []struct {
		name             string
		setupContext     *CallContext
		setupCredentials map[string]string // schemeName -> token
		agentCard        types.AgentCard
		expectedHeaders  map[string]string
	}{
		{
			name:             "no security requirements",
			setupContext:     NewCallContext(1),
			setupCredentials: map[string]string{},
			agentCard: types.AgentCard{
				Security:        types.SecurityRequirement{},
				SecuritySchemes: map[string]types.SecurityScheme{},
			},
			expectedHeaders: map[string]string{},
		},
		{
			name: "API key security scheme - header",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
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
			},
			expectedHeaders: map[string]string{
				"X-API-Key": "api-token-123",
			},
		},
		{
			name: "API key security scheme - query (should not set header)",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"apiKey": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"apiKey": types.APIKeySecurityScheme{
						Type: types.APIKEY,
						In:   types.InQuery,
						Name: "api_key",
					},
				},
			},
			expectedHeaders: map[string]string{},
		},
		{
			name: "HTTP Bearer security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"bearer": "bearer-token-456",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"bearer": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"bearer": types.HTTPAuthSecurityScheme{
						Type:   types.HTTP,
						Scheme: "Bearer",
					},
				},
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer bearer-token-456",
			},
		},
		{
			name: "HTTP Basic security scheme (should not set header)",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"basic": "basic-token-789",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"basic": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"basic": types.HTTPAuthSecurityScheme{
						Type:   types.HTTP,
						Scheme: "Basic",
					},
				},
			},
			expectedHeaders: map[string]string{},
		},
		{
			name: "OAuth2 security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"oauth2": "oauth-token-abc",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"oauth2": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"oauth2": types.OAuth2SecurityScheme{
						Type: types.OAUTH2,
					},
				},
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer oauth-token-abc",
			},
		},
		{
			name: "OpenID Connect security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"openid": "openid-token-def",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"openid": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"openid": types.OpenIdConnectSecurityScheme{
						Type: types.OPENIDConnect,
					},
				},
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer openid-token-def",
			},
		},
		{
			name: "multiple security requirements",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
				"bearer": "bearer-token-456",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"apiKey": []string{}},
					{"bearer": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"apiKey": types.APIKeySecurityScheme{
						Type: types.APIKEY,
						In:   types.InHeader,
						Name: "X-API-Key",
					},
					"bearer": types.HTTPAuthSecurityScheme{
						Type:   types.HTTP,
						Scheme: "Bearer",
					},
				},
			},
			expectedHeaders: map[string]string{
				"X-API-Key":     "api-token-123",
				"Authorization": "Bearer bearer-token-456",
			},
		},
		{
			name: "credential not found",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{},
			agentCard: types.AgentCard{
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
			},
			expectedHeaders: map[string]string{},
		},
		{
			name: "security scheme not found",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"apiKey": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{},
			},
			expectedHeaders: map[string]string{},
		},
		{
			name: "nil security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.State["sessionId"] = "session1"
				return ctx
			}(),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
				Security: types.SecurityRequirement{
					{"apiKey": []string{}},
				},
				SecuritySchemes: map[string]types.SecurityScheme{
					"apiKey": nil,
				},
			},
			expectedHeaders: map[string]string{},
		},
		{
			name:         "nil context",
			setupContext: nil,
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
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
			},
			expectedHeaders: map[string]string{},
		},
		{
			name:         "empty sessionId in context",
			setupContext: NewCallContext(1),
			setupCredentials: map[string]string{
				"apiKey": "api-token-123",
			},
			agentCard: types.AgentCard{
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
			},
			expectedHeaders: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credentials := NewInMemoryContextCredentials()
			for schemeName, token := range tt.setupCredentials {
				if tt.setupContext != nil {
					if sessionId, ok := tt.setupContext.State["sessionId"]; ok {
						if sessionIdStr, ok := sessionId.(string); ok {
							credentials.SetCredentials(sessionIdStr, schemeName, token)
						}
					}
				}
			}

			interceptor := NewInterceptor(credentials)

			request, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
			require.NoError(t, err)

			interceptor.Intercept(request, tt.setupContext, tt.agentCard)

			for expectedHeader, expectedValue := range tt.expectedHeaders {
				actualValue := request.Header.Get(expectedHeader)
				assert.Equal(t, expectedValue, actualValue,
					"Header %s should be set to %s, but got %s",
					expectedHeader, expectedValue, actualValue)
			}

			if len(tt.expectedHeaders) == 0 {
				assert.Empty(t, request.Header)
			}
		})
	}
}
