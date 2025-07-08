package client

// CallContext a context passed with each client call
type CallContext struct {
	State map[string]any
}

func NewCallContext(size uint) *CallContext {
	return &CallContext{State: make(map[string]any, size)}
}
