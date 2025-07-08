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

package types

type SecurityScheme interface {
	GetType() string
}

// APIKeySecurityScheme API key security scheme
type APIKeySecurityScheme struct {
	Description string `json:"description,omitempty"`
	In          In     `json:"in,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

func (a APIKeySecurityScheme) GetType() string {
	return APIKEY
}

// HTTPAuthSecurityScheme HTTP Authentication security scheme.
type HTTPAuthSecurityScheme struct {
	Scheme      string `json:"scheme"`
	Description string `json:"description,omitempty"`
	BaseFormat  string `json:"base_format,omitempty"`
	Type        string `json:"type,omitempty"`
}

func (h HTTPAuthSecurityScheme) GetType() string {
	return HTTP
}

// OAuth2SecurityScheme OAuth Security scheme configuration
type OAuth2SecurityScheme struct {
	Description string `json:"description"`
	Flows       any    `json:"flows"`
	Type        string `json:"type"`
}

func (h OAuth2SecurityScheme) GetType() string {
	return OAUTH2
}

type OpenIdConnectSecurityScheme struct {
	Description      string `json:"description,omitempty"`
	OpenIdConnectUrl string `json:"open_id_connect_url,omitempty"`
	Type             string `json:"type,omitempty"`
}

func (h OpenIdConnectSecurityScheme) GetType() string {
	return OPENIDConnect
}

type OAuthFlows struct {
	AuthorizationCode AuthorizationCodeOAuthFlow `json:"authorization_code,omitempty"`
	ClientCredentials ClientCredentialsOAuthFlow `json:"client_credentials,omitempty"`
}

type AuthorizationCodeOAuthFlow struct {
	AuthorizationUrl string            `json:"authorization_url,omitempty"`
	RefreshUrl       string            `json:"refresh_url,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
	TokenUrl         string            `json:"token_url,omitempty"`
}

type ImplicitOAuthFlow struct {
	AuthorizationUrl string            `json:"authorization_url,omitempty"`
	RefreshUrl       string            `json:"refresh_url,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// ClientCredentialsOAuthFlow Configuration details for a supported OAuth Flow
type ClientCredentialsOAuthFlow struct {
	TokenUrl   string            `json:"token_url"`
	RefreshUrl string            `json:"refresh_url,omitempty"`
	Scopes     map[string]string `json:"scopes,omitempty"`
}

// In the location
type In string

const (
	InCookie In = "cookie"
	InHeader In = "header"
	InQuery  In = "query"
)

const (
	APIKEY        = "api_key"
	HTTP          = "http"
	OAUTH2        = "oauth2"
	OPENIDConnect = "openIdConnect"
)
