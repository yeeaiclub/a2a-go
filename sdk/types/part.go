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

type Part interface {
	GetKind() string
	GetMetadata() map[string]any
}

type FileContent interface {
	GetMimeType() string
	GetName() string
	GetKind() string
}

// DataPart Represents a structured data segment within a message part.
type DataPart struct {
	Data     map[string]any `json:"data,omitempty"`
	Kind     string         `json:"kind"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (d *DataPart) GetKind() string {
	return "data"
}

func (d *DataPart) GetMetadata() map[string]any {
	return d.Metadata
}

type FilePart struct {
	File     FileContent    `json:"file"`
	Kind     string         `json:"kind,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (f *FilePart) GetKind() string {
	return "file"
}

func (f *FilePart) GetMetadata() map[string]any {
	return f.Metadata
}

func UnmarshalFileContent(data json.RawMessage) (FileContent, error) {
	var probe struct {
		Bytes *string `json:"bytes,omitempty"`
		Url   *string `json:"url,omitempty"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("failed to probe file content: %w", err)
	}

	switch {
	case probe.Bytes != nil:
		var withBytes FileWithBytes
		if err := json.Unmarshal(data, &withBytes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileWithBytes: %w", err)
		}
		return &withBytes, nil
	case probe.Url != nil:
		var withUrl FileWithUrl
		if err := json.Unmarshal(data, &withUrl); err != nil {
			return nil, fmt.Errorf("failed to unmarshal FileWithUrl: %w", err)
		}
		return &withUrl, nil
	default:
		return nil, fmt.Errorf("unknown file content type")
	}
}

func (f *FilePart) UnmarshalJSON(data []byte) error {
	aux := &struct {
		File     json.RawMessage `json:"file"`
		Kind     string          `json:"kind,omitempty"`
		Metadata map[string]any  `json:"metadata,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	f.Kind = aux.Kind
	f.Metadata = aux.Metadata

	fileContent, err := UnmarshalFileContent(aux.File)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file content: %w", err)
	}

	f.File = fileContent
	return nil
}

type FileBase struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
}

type FileWithBytes struct {
	FileBase
	Bytes string `json:"bytes,omitempty"`
}

func (fb *FileWithBytes) GetMimeType() string {
	return fb.MimeType
}

func (fb *FileWithBytes) GetName() string {
	return fb.Name
}

func (fb *FileWithBytes) GetKind() string {
	return "bytes"
}

type FileWithUrl struct {
	FileBase
	Url string `json:"url,omitempty"`
}

func (fu *FileWithUrl) GetMimeType() string {
	return fu.MimeType
}

func (fu *FileWithUrl) GetName() string {
	return fu.Name
}

func (fu *FileWithUrl) GetKind() string {
	return "url"
}

type TextPart struct {
	Kind     string         `json:"kind"`
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (t *TextPart) GetKind() string {
	return "text"
}

func (t *TextPart) GetMetadata() map[string]any {
	return t.Metadata
}
