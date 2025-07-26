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

// Part type constants
const (
	PartTypeText = "text"
	PartTypeFile = "file"
	PartTypeData = "data"
)

type Part interface {
	GetKind() string
	GetMetadata() map[string]any
}

type FileContent struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
	Bytes    string `json:"bytes,omitempty"`
	Url      string `json:"url,omitempty"`
}

// DataPart Represents a structured data segment within a message part.
type DataPart struct {
	Data     map[string]any `json:"data,omitempty"`
	Kind     string         `json:"kind"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (d *DataPart) GetKind() string {
	return PartTypeData
}

func (d *DataPart) GetMetadata() map[string]any {
	return d.Metadata
}

type FilePart struct {
	File     FileContent    `json:"file"`
	Kind     string         `json:"kind,omitempty"`
	Bytes    string         `json:"bytes,omitempty"`
	Url      string         `json:"url,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (f *FilePart) GetKind() string {
	return PartTypeFile
}

func (f *FilePart) GetMetadata() map[string]any {
	return f.Metadata
}

type TextPart struct {
	Kind     string         `json:"kind"`
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (t *TextPart) GetKind() string {
	return PartTypeText
}

func (t *TextPart) GetMetadata() map[string]any {
	return t.Metadata
}
