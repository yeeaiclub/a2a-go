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

package aggregator

import (
	"context"

	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type FullConsumer struct {
	manager *manager.TaskManager
}

func NewFullConsumer(manager *manager.TaskManager) *FullConsumer {
	return &FullConsumer{manager: manager}
}

func (c *FullConsumer) Consume(ctx context.Context, queue *event.Queue) (types.Event, error) {
	for e := range queue.Subscribe(ctx) {
		switch e.Type {
		case types.EventCanceled:
			return nil, ctx.Err()
		case types.EventError:
			return nil, e.Err
		case types.EventClosed:
			return nil, nil
		case types.EventDone:
			if _, err := c.manager.Process(ctx, e.Event); err != nil {
				return nil, err
			}
			return c.manager.GetTask(ctx)
		case types.EventData:
			if msg, ok := e.Event.(*types.Message); ok {
				return msg, nil
			}
			if _, err := c.manager.Process(ctx, e.Event); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
