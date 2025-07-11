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
	"time"

	"github.com/yeeaiclub/a2a-go/sdk/auth"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// callContextPool is a pool for reusing CallContext objects
var callContextPool = sync.Pool{
	New: func() any {
		return &CallContext{
			state:           make(map[string]any),
			securitySchemes: make(map[string]types.SecurityScheme),
		}
	},
}

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

	// state for storing arbitrary key-value pairs
	state map[string]any

	// Mutex for thread-safe access to shared fields
	mu sync.RWMutex

	// Context for cancellation and timeout
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCallContext creates a new CallContext with the given context
// If no context is provided, it uses context.Background()
func NewCallContext(ctx context.Context) *CallContext {
	if ctx == nil {
		ctx = context.Background()
	}

	// Get a context from the pool
	callCtx := callContextPool.Get().(*CallContext)

	// Create a new context with cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	callCtx.ctx = cancelCtx
	callCtx.cancel = cancel

	return callCtx
}

// NewCallContextWithRequest creates a new CallContext with an HTTP request
// It uses the request's context as the base context for proper lifecycle management
func NewCallContextWithRequest(req *http.Request) *CallContext {
	callCtx := NewCallContext(req.Context())
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
	if c.state == nil {
		c.state = make(map[string]any)
	}
	c.state[key] = value
}

// Get retrieves a value from the metadata
func (c *CallContext) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.state == nil {
		return nil
	}
	return c.state[key]
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

// Deadline returns the time when work done on behalf of this context
// should be canceled, or ok==false if no deadline is set
func (c *CallContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Err returns a non-nil error value after Done is closed
func (c *CallContext) Err() error {
	return c.ctx.Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. This delegates to the underlying context.
func (c *CallContext) Value(key any) any {
	return c.ctx.Value(key)
}

// reset clears all fields and prepares the context for reuse
func (c *CallContext) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cancel the current context if it exists
	if c.cancel != nil {
		c.cancel()
	}

	// Clear all fields
	c.User = nil
	c.request = nil
	c.security = types.SecurityRequirement{}
	c.securitySchemes = make(map[string]types.SecurityScheme)
	c.state = make(map[string]any)
	c.ctx = nil
	c.cancel = nil
}

// Release returns the CallContext to the pool for reuse
// This should be called when the context is no longer needed
func (c *CallContext) Release() {
	c.reset()
	callContextPool.Put(c)
}
