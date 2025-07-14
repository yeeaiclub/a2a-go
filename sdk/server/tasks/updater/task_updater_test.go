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

package updater

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		taskId    string
		contextId string
		state     types.TaskState
		final     bool
	}{
		{"working status", "tid", "cid", types.WORKING, false},
		{"completed status", "tid", "cid", types.COMPLETED, true},
		{"failed status", "tid", "cid", types.FAILED, true},
		{"rejected status", "tid", "cid", types.REJECTED, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(10)
			defer queue.Close()
			updater := NewTaskUpdater(queue, tc.taskId, tc.contextId)
			updater.UpdateStatus(tc.state)
			ch := queue.Subscribe(context.Background())
			e := <-ch
			statusEvent, ok := e.Event.(*types.TaskStatusUpdateEvent)
			assert.True(t, ok)
			assert.Equal(t, tc.taskId, statusEvent.TaskId)
			assert.Equal(t, tc.contextId, statusEvent.ContextId)
			assert.Equal(t, tc.state, statusEvent.Status.State)
			if tc.final {
				assert.True(t, statusEvent.Final)
			} else {
				assert.False(t, statusEvent.Final)
			}
			assert.NotEmpty(t, statusEvent.Status.TimeStamp)
		})
	}
}

func TestComplete_Failed_Reject(t *testing.T) {
	tests := []struct {
		name  string
		fn    func(updater *TaskUpdater)
		state types.TaskState
		final bool
	}{
		{"complete", func(u *TaskUpdater) { u.Complete() }, types.COMPLETED, true},
		{"failed", func(u *TaskUpdater) { u.Failed() }, types.FAILED, true},
		{"reject", func(u *TaskUpdater) { u.Reject() }, types.REJECTED, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(10)
			updater := NewTaskUpdater(queue, "tid", "cid")
			tc.fn(updater)
			ch := queue.Subscribe(context.Background())
			e := <-ch
			statusEvent, ok := e.Event.(*types.TaskStatusUpdateEvent)
			assert.True(t, ok)
			assert.Equal(t, tc.state, statusEvent.Status.State)
			assert.True(t, statusEvent.Final)
			queue.Close()
		})
	}
}

func TestAddArtifact(t *testing.T) {
	t.Run("add artifact", func(t *testing.T) {
		queue := event.NewQueue(10)
		defer queue.Close()
		updater := NewTaskUpdater(queue, "tid", "cid")
		parts := []types.Part{&types.TextPart{Kind: "text", Text: "hello"}}
		updater.AddArtifact(parts, WithName("artifact1"))
		ch := queue.Subscribe(context.Background())
		e := <-ch
		artifactEvent, ok := e.Event.(*types.TaskArtifactUpdateEvent)
		assert.True(t, ok)
		assert.Equal(t, "tid", artifactEvent.TaskId)
		assert.Equal(t, "cid", artifactEvent.ContextId)
		assert.NotNil(t, artifactEvent.Artifact)
		assert.Equal(t, "artifact1", artifactEvent.Artifact.Name)
		assert.Len(t, artifactEvent.Artifact.Parts, 1)
	})
}

func TestNewAgentMessage(t *testing.T) {
	t.Run("new agent message", func(t *testing.T) {
		updater := NewTaskUpdater(nil, "tid", "cid")
		parts := []types.Part{&types.TextPart{Kind: "text", Text: "msg"}}
		msg := updater.NewAgentMessage(parts, WithMetadata(map[string]any{"foo": "bar"}))
		assert.Equal(t, types.Agent, msg.Role)
		assert.Equal(t, "tid", msg.TaskID)
		assert.Equal(t, "cid", msg.ContextID)
		assert.Len(t, msg.Parts, 1)
		assert.Equal(t, "bar", msg.Metadata["foo"])
		assert.NotEmpty(t, msg.MessageID)
	})
}

func TestTimestampOption(t *testing.T) {
	t.Run("timestamp option", func(t *testing.T) {
		queue := event.NewQueue(10)
		defer queue.Close()
		updater := NewTaskUpdater(queue, "tid", "cid")
		ts := time.Now().Add(-time.Hour).Format(time.RFC3339)
		updater.UpdateStatus(types.SUBMITTED, WithTimestamp(ts))
		ch := queue.Subscribe(context.Background())
		e := <-ch
		statusEvent, ok := e.Event.(*types.TaskStatusUpdateEvent)
		assert.True(t, ok)
		assert.Equal(t, ts, statusEvent.Status.TimeStamp)
	})
}
