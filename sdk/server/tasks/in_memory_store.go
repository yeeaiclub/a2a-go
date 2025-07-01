// Copyright 2025 yumosx
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

package tasks

import (
	"context"
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/types"
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
	delete(s.tasks, taskID)
	return nil
}
