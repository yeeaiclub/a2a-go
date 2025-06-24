package tasks

import (
	"context"
	"sync"

	"github.com/yumosx/a2a-go/sdk/types"
)

type InMemoryTaskStore struct {
	tasks map[string]*types.Task
	mu    sync.Mutex
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{
		tasks: make(map[string]*types.Task),
	}
}

func (s *InMemoryTaskStore) Save(ctx context.Context, task *types.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.GetTaskId()] = task
	return nil
}

func (s *InMemoryTaskStore) Get(ctx context.Context, taskID string) (*types.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, exists := s.tasks[taskID]; exists {
		return task, nil
	}
	return nil, nil
}

func (s *InMemoryTaskStore) Delete(ctx context.Context, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; exists {
		delete(s.tasks, taskID)
	}
	return nil
}
