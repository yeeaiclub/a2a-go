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
	InHeader    = "header"
	InQuery     = "query"
)
