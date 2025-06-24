package errs

import "errors"

var (
	QueueEmpty           = errors.New("queue is empty")
	QueueClosed          = errors.New("queue is closed")
	UnSupportedOperation = errors.New("this operation is not supported")
	TaskNotFound         = errors.New("task not found")
)
