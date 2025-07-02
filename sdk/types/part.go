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

type Part interface {
	GetPart() string
}

// DataPart Represents a structured data segment within a message part.
type DataPart struct {
	Data     map[string]any `json:"data,omitempty"`
	Kind     string         `json:"kind,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type FileBase struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
}

type FileWithBytes struct {
	Bytes    string `json:"bytes,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
}

type FileWithUrl struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
	Url      string `json:"url,omitempty"`
}

type TextPart struct {
	Kind     string         `json:"kind"`
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (t *TextPart) GetPart() string {
	return "text"
}
