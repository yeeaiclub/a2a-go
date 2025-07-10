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
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// CallContext a context passed with each client call
type CallContext struct {
	mu              sync.Mutex
	state           map[string]any
	request         *http.Request
	Security        types.SecurityRequirement
	SecuritySchemes map[string]types.SecurityScheme
}

func (c *CallContext) SetRequest(req *http.Request) {
	c.request = req
}

func (c *CallContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == nil {
		c.state = make(map[string]any)
	}
	c.state[key] = value
}

func (c *CallContext) Get(key string) any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state[key]
}

func (c *CallContext) Request() *http.Request {
	return c.request
}

func (c *CallContext) GetSecurityRequirement() types.SecurityRequirement {
	return c.Security
}

func (c *CallContext) GetSecuritySchemes(key string) types.SecurityScheme {
	return c.SecuritySchemes[key]
}

func (c *CallContext) SetSecurityConfig(security types.SecurityRequirement, schemes map[string]types.SecurityScheme) {
	c.Security = security
	c.SecuritySchemes = schemes
}

func NewCallContext(size uint) *CallContext {
	return &CallContext{state: make(map[string]any, size)}
}
