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

package tasks

import (
	"context"

	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// TaskStore Agent task store interface
type TaskStore interface {
	// Save or updates a task in the store.
	Save(ctx context.Context, task *types.Task) error
	// Get retrieves a task from the store by Id.
	Get(ctx context.Context, taskId string) (*types.Task, error)
	// Delete a tasks from the store by id
	Delete(ctx context.Context, id string) error
}
