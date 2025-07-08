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

package event

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestEnqueue(t *testing.T) {
	testcases := []struct {
		name   string
		before func(queue *Queue)
		want   []types.StreamEvent
	}{
		{
			name: "enqueue and done",
			before: func(queue *Queue) {
				queue.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				queue.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: true})
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false}, Type: types.EventData},
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: true}, Type: types.EventDone},
			},
		},
		{
			name: "enqueue and error",
			before: func(queue *Queue) {
				queue.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				queue.EnqueueError(errors.New("error"))
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false}, Type: types.EventData},
				{Err: errors.New("error"), Type: types.EventError},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			tc.before(queue)
			queue.Close()
			list := make([]types.StreamEvent, 0)
			for e := range queue.ch {
				list = append(list, e)
			}
			assert.ElementsMatch(t, tc.want, list)
		})
	}
}

func TestSubscribe(t *testing.T) {
	testcases := []struct {
		name   string
		before func(queue *Queue)
		want   []types.StreamEvent
	}{
		{
			name: "subscribe",
			before: func(queue *Queue) {
				queue.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				queue.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: true})
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false}, Type: types.EventData},
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: true}, Type: types.EventDone},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			defer queue.Close()
			tc.before(queue)
			events := queue.Subscribe(context.Background())
			list := make([]types.StreamEvent, 0)
			for ev := range events {
				list = append(list, ev)
			}
			assert.ElementsMatch(t, tc.want, list)
		})
	}
}
