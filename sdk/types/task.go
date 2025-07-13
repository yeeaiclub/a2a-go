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
)

type TaskState string

const (
	SUBMITTED     TaskState = "submitted"
	WORKING       TaskState = "working"
	InputRequired TaskState = "required"
	COMPLETED     TaskState = "completed"
	CANCELED      TaskState = "canceled"
	FAILED        TaskState = "failed"
	REJECTED      TaskState = "rejected"
	AUTH_REQUIRED TaskState = "auth_required"
	UNKNOWN       TaskState = "unknown"
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

func (a *Artifact) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ArtifactId  string            `json:"artifact_id,omitempty"`
		Description string            `json:"description,omitempty"`
		Extensions  []string          `json:"extension,omitempty"`
		Metadata    map[string]any    `json:"metadata,omitempty"`
		Name        string            `json:"name,omitempty"`
		Parts       []json.RawMessage `json:"parts,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	a.ArtifactId = aux.ArtifactId
	a.Description = aux.Description
	a.Extensions = aux.Extensions
	a.Metadata = aux.Metadata
	a.Name = aux.Name

	// Unmarshal parts
	a.Parts = make([]Part, len(aux.Parts))
	for i, partData := range aux.Parts {
		part, err := UnmarshalPart(partData)
		if err != nil {
			return fmt.Errorf("failed to unmarshal part %d: %w", i, err)
		}
		a.Parts[i] = part
	}

	return nil
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

func (t *Task) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Id        string            `json:"id"`
		ContextId string            `json:"context_id"`
		History   []*Message        `json:"history,omitempty"`
		Kind      string            `json:"kind,omitempty"`
		Status    TaskStatus        `json:"task_status,omitempty"`
		Metadata  map[string]any    `json:"metadata,omitempty"`
		Artifacts []json.RawMessage `json:"artifacts,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	t.Id = aux.Id
	t.ContextId = aux.ContextId
	t.History = aux.History
	t.Kind = aux.Kind
	t.Status = aux.Status
	t.Metadata = aux.Metadata

	if aux.Artifacts == nil {
		t.Artifacts = []Artifact{}
	} else {
		// Unmarshal artifacts
		t.Artifacts = make([]Artifact, len(aux.Artifacts))
		for i, artifactData := range aux.Artifacts {
			var artifact Artifact
			if err := json.Unmarshal(artifactData, &artifact); err != nil {
				return fmt.Errorf("failed to unmarshal artifact %d: %w", i, err)
			}
			t.Artifacts[i] = artifact
		}
	}

	return nil
}
