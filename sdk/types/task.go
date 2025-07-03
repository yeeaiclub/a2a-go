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

type TaskState string

const (
	SUBMITTED     = "submitted"
	WORKING       = "working"
	InputRequired = "required"
	COMPLETED     = "completed"
	CANCELED      = "canceled"
	FAILED        = "failed"
	REJECTED      = "rejected"
	AUTH_REQUIRED = "auth_required"
	UNKNOWN       = "unknown"
)

type TaskStatus struct {
	Message   *Message  `json:"message,omitempty"`
	State     TaskState `json:"state"`
	TimeStamp string    `json:"time_stamp,omitempty"`
}

type Artifact struct {
	ArtifactId  string         `json:"artifact_id,omitempty"`
	Description string         `json:"description,omitempty"`
	Extensions  []string       `json:"extension,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Name        string         `json:"name,omitempty"`
	Parts       []Part         `json:"parts,omitempty"`
}

type Task struct {
	Id        string         `json:"id"`
	ContextId string         `json:"context_id"`
	History   []*Message     `json:"history,omitempty"`
	Kind      string         `json:"kind,omitempty"`
	Status    TaskStatus     `json:"task_status,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Artifacts []Artifact     `json:"artifacts,omitempty"`
}

func (t *Task) Done() bool {
	return t.Status.State == COMPLETED ||
		t.Status.State == CANCELED ||
		t.Status.State == FAILED ||
		t.Status.State == REJECTED ||
		t.Status.State == UNKNOWN ||
		t.Status.State == InputRequired
}

func (t *Task) GetContextId() string {
	return t.ContextId
}

func (t *Task) GetTaskId() string {
	return t.Id
}

func (t *Task) EventType() string {
	return "task"
}
