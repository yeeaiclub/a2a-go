package client

import (
	"github.com/yeeaiclub/a2a-go/sdk/client/middleware"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// BeforeArgs holds information passed to the interceptor before a method call.
type BeforeArgs struct {
	// Input is the request payload for the method call.
	Input any

	// Method is the name of the method being called.
	Method string

	// AgentCard is the agent card associated with the call.
	AgentCard *types.AgentCard

	// Context carries per-call state such as HTTP request and security config.
	Context *middleware.CallContext

	// EarlyReturn, if set by Before, short-circuits the call and uses this value as the result.
	EarlyReturn any
}

// AfterArgs holds information passed to the interceptor after a method call completes.
type AfterArgs struct {
	// Result is the response from the method call.
	Result any

	// Method is the name of the method that was called.
	Method string

	// AgentCard is the agent card associated with the call.
	AgentCard *types.AgentCard

	// Context carries per-call state such as HTTP request and security config.
	Context *middleware.CallContext

	// EarlyReturn indicates whether the call was short-circuited by a Before interceptor.
	EarlyReturn bool
}

// ClientCallInterceptor defines the interface for client-side call interceptors.
// Interceptors can inspect and modify requests before they are sent,
// which is ideal for concerns like authentication, logging, or tracing.
type ClientCallInterceptor interface {
	// Before is invoked before a transport method call.
	// Return an error to prevent the call from proceeding.
	Before(args *BeforeArgs) error

	// After is invoked after a transport method call completes.
	// Return an error to signal a failure in the interceptor itself.
	After(args *AfterArgs) error
}
