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

package types

// APIKeySecurityScheme API key security scheme
type APIKeySecurityScheme struct {
	Description string `json:"description,omitempty"`
	In          In     `json:"in,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

// HTTPAuthSecurityScheme HTTP Authentication security scheme.
type HTTPAuthSecurityScheme struct {
	Scheme      string `json:"scheme"`
	Description string `json:"description,omitempty"`
	BaseFormat  string `json:"base_format,omitempty"`
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
