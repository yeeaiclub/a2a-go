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

	"github.com/yeeaiclub/a2a-go/internal/jsonx"
)

type Role string

const (
	Agent Role = "agent"
	User  Role = "user"
)

// Message Represents a single message exchanged between user and agent
type Message struct {
	Role             Role           `json:"role"`
	TaskID           string         `json:"task_id,omitempty"`
	ContextID        string         `json:"context_id,omitempty"`
	Extensions       []string       `json:"extensions,omitempty"`
	Kind             string         `json:"kind,omitempty"`
	MessageID        string         `json:"message_id,omitempty"`
	ReferenceTaskIDs []string       `json:"referenceTaskIds,omitempty"`
	Parts            []Part         `json:"parts,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

func (m *Message) GetTaskId() string {
	return m.TaskID
}

func (m *Message) GetContextId() string {
	return m.ContextID
}

func (m *Message) Type() string {
	return EventTypeMessage
}

func (m *Message) Done() bool {
	return true
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message // avoid recursion
	aux := &struct {
		Parts []json.RawMessage `json:"parts,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	kindMap := map[string]func() Part{
		"text": func() Part { return &TextPart{} },
		"data": func() Part { return &DataPart{} },
		"file": func() Part { return &FilePart{} },
	}

	parts, err := jsonx.UnmarshalSliceByKind(aux.Parts, kindMap)
	if err != nil {
		return err
	}

	m.Parts = parts
	return nil
}
