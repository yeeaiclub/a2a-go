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

	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestConsumeAll(t *testing.T) {
	testcases := []struct {
		name   string
		before func(q *Queue)
		want   []types.StreamEvent
	}{
		{
			name: "consume single task",
			before: func(q *Queue) {
				q.Enqueue(&types.Task{Id: "1", Status: types.TaskStatus{State: types.FAILED}})
			},
			want: []types.StreamEvent{
				{Event: &types.Task{Id: "1", Status: types.TaskStatus{State: types.FAILED}}},
			},
		},
		{
			name: "consume multiple tasks",
			before: func(q *Queue) {
				q.Enqueue(&types.Task{Id: "1"})
				q.Enqueue(&types.Task{Id: "2", Status: types.TaskStatus{State: types.FAILED}})
			},
			want: []types.StreamEvent{
				{Event: &types.Task{Id: "1"}},
				{Event: &types.Task{Id: "2", Status: types.TaskStatus{State: types.FAILED}}},
			},
		},
		{
			name: "consume multiple task update event",
			before: func(q *Queue) {
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1"})
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", Final: true})
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1"}},
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", Final: true}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			defer queue.Close()
			tc.before(queue)

			consumer := NewConsumer(queue, nil)
			events := consumer.ConsumeAll(context.Background(), 2)

			var received []types.StreamEvent
			for event := range events {
				require.NoError(t, event.Err)
				received = append(received, event)
			}
			require.ElementsMatch(t, tc.want, received)
		})
	}
}
