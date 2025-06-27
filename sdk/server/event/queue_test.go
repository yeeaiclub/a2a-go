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

package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yumosx/a2a-go/sdk/types"
)

func TestEnqueue(t *testing.T) {
	testcases := []struct {
		name  string
		input types.Event
		size  int
	}{
		{
			name:  "enqueue",
			size:  1,
			input: &types.TaskStatusUpdateEvent{TaskId: "1", Final: true},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(tc.size)
			defer queue.Close()
			queue.Enqueue(tc.input)
			//todo: refactor
			event := <-queue.queue
			assert.Equal(t, event, tc.input)
		})
	}
}

func TestDequeueNoWait(t *testing.T) {
	testcases := []struct {
		name   string
		before func(q *Queue)
		want   types.Event
	}{
		{
			name: "dequeue not wait",
			before: func(q *Queue) {
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", Final: true})
			},
			want: &types.TaskStatusUpdateEvent{
				TaskId: "1",
				Final:  true,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			defer queue.Close()
			tc.before(queue)
			wait := queue.DequeueNoWait(context.Background())
			assert.Equal(t, wait.Event, tc.want)
		})
	}
}

func TestDequeueWait(t *testing.T) {
	testcases := []struct {
		name   string
		before func(q *Queue)
		want   types.Event
	}{
		{
			name: "dequeue wait",
			before: func(q *Queue) {
				go func() {
					time.Sleep(2 * time.Second)
					q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", Final: true})
				}()
			},
			want: &types.TaskStatusUpdateEvent{TaskId: "1", Final: true},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			defer queue.Close()
			tc.before(queue)
			wait := queue.DequeueWait(context.Background())
			assert.Equal(t, wait.Event, tc.want)
		})
	}
}
