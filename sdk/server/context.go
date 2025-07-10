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

package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/auth"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// CallContext represents the context for a single API call
// It contains user information, request data, security configuration, and metadata
type CallContext struct {
	// User information
	User auth.User

	// Request information
	request *http.Request

	// Security configuration
	security        types.SecurityRequirement
	securitySchemes map[string]types.SecurityScheme

	// Metadata for storing arbitrary key-value pairs
	Metadata map[string]any

	// Mutex for thread-safe access to shared fields
	mu sync.RWMutex

	// Context for cancellation and timeout
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCallContext creates a new CallContext with default values
func NewCallContext() *CallContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &CallContext{
		Metadata:        make(map[string]any),
		securitySchemes: make(map[string]types.SecurityScheme),
		ctx:             ctx,
		cancel:          cancel,
	}
}

// NewCallContextWithRequest creates a new CallContext with an HTTP request
func NewCallContextWithRequest(req *http.Request) *CallContext {
	callCtx := NewCallContext()
	callCtx.SetRequest(req)
	return callCtx
}

// SetRequest sets the HTTP request for this context
func (c *CallContext) SetRequest(req *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.request = req
}

// Request returns the HTTP request associated with this context
func (c *CallContext) Request() *http.Request {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.request
}

// SetUser sets the user information for this context
func (c *CallContext) SetUser(user auth.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.User = user
}

// GetUser returns the user information
func (c *CallContext) GetUser() auth.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.User
}

// Set stores a key-value pair in the metadata
func (c *CallContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Metadata == nil {
		c.Metadata = make(map[string]any)
	}
	c.Metadata[key] = value
}

// Get retrieves a value from the metadata
func (c *CallContext) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Metadata == nil {
		return nil
	}
	return c.Metadata[key]
}

// GetSecurityRequirement returns the security requirement for this context
func (c *CallContext) GetSecurityRequirement() types.SecurityRequirement {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.security
}

// GetSecuritySchemes returns the security scheme for a given key
func (c *CallContext) GetSecuritySchemes(key string) types.SecurityScheme {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.securitySchemes[key]
}

// SetSecurityConfig sets the security configuration for this context
func (c *CallContext) SetSecurityConfig(security types.SecurityRequirement, schemes map[string]types.SecurityScheme) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.security = security
	c.securitySchemes = schemes
}

// GetAllSecuritySchemes returns all security schemes
func (c *CallContext) GetAllSecuritySchemes() map[string]types.SecurityScheme {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to avoid external modification
	schemes := make(map[string]types.SecurityScheme, len(c.securitySchemes))
	for k, v := range c.securitySchemes {
		schemes[k] = v
	}
	return schemes
}

// Context returns the underlying context.Context
func (c *CallContext) Context() context.Context {
	return c.ctx
}

// Cancel cancels the context
func (c *CallContext) Cancel() {
	c.cancel()
}

// Done returns a channel that's closed when the context is cancelled
func (c *CallContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *CallContext) Err() error {
	return c.ctx.Err()
}
