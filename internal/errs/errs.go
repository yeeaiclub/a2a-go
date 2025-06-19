package errs

import "errors"

var (
	QueueEmpty  = errors.New("queue is empty")
	QueueClosed = errors.New("queue is closed")
)

var (
	AgentNoResponse = errors.New("agent no response")
)
