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

package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/sdk/types"
	"github.com/yeeaiclub/a2a-go/sdk/web"
)

type MockCredential struct {
	credentials map[string]string // schemeName -> token
	errors      map[string]error  // schemeName -> error
}

func NewMockCredential() *MockCredential {
	return &MockCredential{
		credentials: make(map[string]string),
		errors:      make(map[string]error),
	}
}

func (m *MockCredential) GetCredentials(securitySchemeName string, context web.Context) (string, error) {
	if err, exists := m.errors[securitySchemeName]; exists {
		return "", err
	}
	if token, exists := m.credentials[securitySchemeName]; exists {
		return token, nil
	}
	return "", nil
}

func (m *MockCredential) SetCredentials(schemeName, token string) {
	m.credentials[schemeName] = token
}

func (m *MockCredential) SetError(schemeName string, err error) {
	m.errors[schemeName] = err
}

func TestIntercept(t *testing.T) {
	tests := []struct {
		name             string
		setupContext     func() *CallContext
		setupCredentials func() *MockCredential
		expectedHeaders  map[string]string
		expectedError    bool
	}{
		{
			name: "no security requirements",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.SetSecurityConfig(types.SecurityRequirement{}, map[string]types.SecurityScheme{})
				return ctx
			},
			setupCredentials: NewMockCredential,
			expectedHeaders:  map[string]string{},
			expectedError:    false,
		},
		{
			name: "API key security scheme - header",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "api-token-123")
				return cred
			},
			expectedHeaders: map[string]string{
				"X-API-Key": "api-token-123",
			},
			expectedError: false,
		},
		{
			name: "API key security scheme - query (should not set header)",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InQuery,
							Name: "api_key",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "api-token-123")
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "HTTP Bearer security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"bearer": []string{}}},
					map[string]types.SecurityScheme{
						"bearer": types.HTTPAuthSecurityScheme{
							Type:   types.HTTP,
							Scheme: "Bearer",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("bearer", "bearer-token-456")
				return cred
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer bearer-token-456",
			},
			expectedError: false,
		},
		{
			name: "HTTP Basic security scheme (should not set header)",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"basic": []string{}}},
					map[string]types.SecurityScheme{
						"basic": types.HTTPAuthSecurityScheme{
							Type:   types.HTTP,
							Scheme: "Basic",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("basic", "basic-token-789")
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "OAuth2 security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"oauth2": []string{}}},
					map[string]types.SecurityScheme{
						"oauth2": types.OAuth2SecurityScheme{
							Type: types.OAUTH2,
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("oauth2", "oauth-token-abc")
				return cred
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer oauth-token-abc",
			},
			expectedError: false,
		},
		{
			name: "OpenID Connect security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"openid": []string{}}},
					map[string]types.SecurityScheme{
						"openid": types.OpenIdConnectSecurityScheme{
							Type: types.OPENIDConnect,
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("openid", "openid-token-def")
				return cred
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer openid-token-def",
			},
			expectedError: false,
		},
		{
			name: "multiple security requirements",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{
						{"apiKey": []string{}},
						{"bearer": []string{}},
					},
					map[string]types.SecurityScheme{
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
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "api-token-123")
				cred.SetCredentials("bearer", "bearer-token-456")
				return cred
			},
			expectedHeaders: map[string]string{
				"X-API-Key":     "api-token-123",
				"Authorization": "Bearer bearer-token-456",
			},
			expectedError: false,
		},
		{
			name: "credential not found",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					},
				)
				return ctx
			},
			setupCredentials: NewMockCredential,
			expectedHeaders:  map[string]string{},
			expectedError:    false,
		},
		{
			name: "credential error",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetError("apiKey", assert.AnError)
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "security scheme not found",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "api-token-123")
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "nil security scheme",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": nil,
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "api-token-123")
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "empty credentials",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"apiKey": []string{}}},
					map[string]types.SecurityScheme{
						"apiKey": types.APIKeySecurityScheme{
							Type: types.APIKEY,
							In:   types.InHeader,
							Name: "X-API-Key",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("apiKey", "")
				return cred
			},
			expectedHeaders: map[string]string{},
			expectedError:   false,
		},
		{
			name: "HTTP Bearer case insensitive",
			setupContext: func() *CallContext {
				ctx := NewCallContext(1)
				ctx.Set("sessionId", "session1")
				ctx.SetSecurityConfig(
					types.SecurityRequirement{{"bearer": []string{}}},
					map[string]types.SecurityScheme{
						"bearer": types.HTTPAuthSecurityScheme{
							Type:   types.HTTP,
							Scheme: "bearer",
						},
					},
				)
				return ctx
			},
			setupCredentials: func() *MockCredential {
				cred := NewMockCredential()
				cred.SetCredentials("bearer", "bearer-token-456")
				return cred
			},
			expectedHeaders: map[string]string{
				"Authorization": "Bearer bearer-token-456",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			credential := tt.setupCredentials()

			request, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
			require.NoError(t, err)
			ctx.SetRequest(request)

			middleware := Intercept(credential)
			handler := middleware(func(ctx web.Context) error {
				return nil
			})

			err = handler(ctx)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

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
