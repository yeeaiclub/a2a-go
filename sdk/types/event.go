package types

type Event interface {
	Done() bool
}

type StreamEvent struct {
	Err   error
	Event Event
}

func (s *StreamEvent) Done() bool {
	return s.Event.Done()
}
