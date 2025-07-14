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

package types

import "encoding/json"

type Event interface {
	GetTaskId() string
	GetContextId() string
	Type() string
	Done() bool
}

// TaskArtifactUpdateEvent Send by server during sendStream and subscribe requests
type TaskArtifactUpdateEvent struct {
	TaskId    string         `json:"task_id"`
	ContextId string         `json:"context_id,omitempty"`
	Kind      string         `json:"kind,omitempty"`
	Artifact  *Artifact      `json:"artifact,omitempty"`
	Append    bool           `json:"append,omitempty"`
	LastChunk bool           `json:"last_chunk,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

func (t *TaskArtifactUpdateEvent) GetTaskId() string {
	return t.TaskId
}

func (t *TaskArtifactUpdateEvent) GetContextId() string {
	return t.ContextId
}

func (t *TaskArtifactUpdateEvent) Type() string {
	return EventTypeArtifactUpdate
}

func (t *TaskArtifactUpdateEvent) Done() bool {
	return false
}

// TaskStatusUpdateEvent Send by server during or subscribe requests
type TaskStatusUpdateEvent struct {
	TaskId    string         `json:"task_id,omitempty"`
	ContextId string         `json:"context_id,omitempty"`
	Final     bool           `json:"final,omitempty"`
	Kind      string         `json:"kind,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Status    TaskStatus     `json:"status,omitempty"`
}

func (t *TaskStatusUpdateEvent) Done() bool {
	return t.Final
}

func (t *TaskStatusUpdateEvent) GetContextId() string {
	return t.ContextId
}

func (t *TaskStatusUpdateEvent) GetTaskId() string {
	return t.TaskId
}

func (t *TaskStatusUpdateEvent) Type() string {
	return EventTypeStatusUpdate
}

type StreamEvent struct {
	Type  EventType
	Event Event
	Err   error
}

func (s *StreamEvent) EncodeJSONRPC(encoder *json.Encoder, id string) error {
	if s.Err != nil {
		return encoder.Encode(InternalError())
	}
	successResp := JSONRPCSuccessResponse(id, s.Event)
	err := encoder.Encode(successResp)
	if err != nil {
		return err
	}
	return nil
}

type EventType int

const (
	EventData EventType = iota
	EventError
	EventDone
	EventClosed
	EventCanceled
)

// Event type constants for all event implementations
const (
	EventTypeStatusUpdate   = "status_update"
	EventTypeArtifactUpdate = "artifact_update"
	EventTypeTask           = "task"
	EventTypeMessage        = "message"
)
