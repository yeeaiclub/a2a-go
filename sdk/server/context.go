package server

import "github.com/yumosx/a2a-go/sdk/auth"

type ServerCallContext struct {
	User     auth.User
	Metadata map[string]any
}
