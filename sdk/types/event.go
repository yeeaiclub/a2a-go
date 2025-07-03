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

import (
	"encoding/json"
	"fmt"
	"io"
)

type Event interface {
	GetTaskId() string
	GetContextId() string
	EventType() string
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

func (t *TaskArtifactUpdateEvent) Done() bool {
	return true
}

func (t *TaskArtifactUpdateEvent) GetContextId() string {
	return t.ContextId
}

func (t *TaskArtifactUpdateEvent) GetTaskId() string {
	return t.TaskId
}

func (t *TaskArtifactUpdateEvent) EventType() string {
	return "artifact_update"
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

func (t *TaskStatusUpdateEvent) EventType() string {
	return "status_update"
}

// StreamEvent the wrap of many events
type StreamEvent struct {
	Err    error
	Event  Event
	Closed bool
}

func (s *StreamEvent) Done() bool {
	return s.Event.Done() || s.Closed
}

func (s *StreamEvent) GetContextId() string {
	return s.Event.GetContextId()
}

func (s *StreamEvent) GetTaskId() string {
	return s.Event.GetTaskId()
}

func (s *StreamEvent) EventType() string {
	return s.Event.EventType()
}

func (s *StreamEvent) MarshalTo(w io.Writer, id string) error {
	if s.Err != nil {
		data, _ := json.Marshal(InternalError())
		if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
			return err
		}
		return nil
	}
	successResp := JSONRPCSuccessResponse(id, s.Event)
	data, err := json.Marshal(successResp)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
		return err
	}
	return nil
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
