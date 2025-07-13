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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtifactUnmarshalJSON(t *testing.T) {
	t.Run("unmarshal artifact with text part", func(t *testing.T) {
		jsonData := `{
			"artifact_id": "test-artifact-1",
			"name": "test artifact",
			"description": "A test artifact",
			"parts": [
				{
					"kind": "text",
					"text": "This is a text part",
					"metadata": {"key": "value"}
				}
			]
		}`

		var artifact Artifact
		err := json.Unmarshal([]byte(jsonData), &artifact)
		require.NoError(t, err)

		assert.Equal(t, "test-artifact-1", artifact.ArtifactId)
		assert.Equal(t, "test artifact", artifact.Name)
		assert.Equal(t, "A test artifact", artifact.Description)
		assert.Len(t, artifact.Parts, 1)

		textPart, ok := artifact.Parts[0].(*TextPart)
		require.True(t, ok)
		assert.Equal(t, "text", textPart.Kind)
		assert.Equal(t, "This is a text part", textPart.Text)
		assert.Equal(t, map[string]any{"key": "value"}, textPart.Metadata)
	})

	t.Run("unmarshal artifact with file part", func(t *testing.T) {
		jsonData := `{
			"artifact_id": "test-artifact-2",
			"name": "file artifact",
			"parts": [
				{
					"kind": "document",
					"file": {
						"bytes": "base64encodeddata",
						"mime_type": "text/plain",
						"name": "test.txt"
					},
					"metadata": {"size": 123}
				}
			]
		}`

		var artifact Artifact
		err := json.Unmarshal([]byte(jsonData), &artifact)
		require.NoError(t, err)

		assert.Equal(t, "test-artifact-2", artifact.ArtifactId)
		assert.Equal(t, "file artifact", artifact.Name)
		assert.Len(t, artifact.Parts, 1)

		filePart, ok := artifact.Parts[0].(*FilePart)
		require.True(t, ok)
		assert.Equal(t, "document", filePart.Kind)
		assert.Equal(t, map[string]any{"size": float64(123)}, filePart.Metadata)

		fileWithBytes, ok := filePart.File.(*FileWithBytes)
		require.True(t, ok)
		assert.Equal(t, "base64encodeddata", fileWithBytes.Bytes)
		assert.Equal(t, "text/plain", fileWithBytes.MimeType)
		assert.Equal(t, "test.txt", fileWithBytes.Name)
	})

	t.Run("unmarshal artifact with data part", func(t *testing.T) {
		jsonData := `{
			"artifact_id": "test-artifact-3",
			"name": "data artifact",
			"parts": [
				{
					"kind": "data",
					"data": {"key1": "value1", "key2": 42},
					"metadata": {"type": "json"}
				}
			]
		}`

		var artifact Artifact
		err := json.Unmarshal([]byte(jsonData), &artifact)
		require.NoError(t, err)

		assert.Equal(t, "test-artifact-3", artifact.ArtifactId)
		assert.Equal(t, "data artifact", artifact.Name)
		assert.Len(t, artifact.Parts, 1)

		dataPart, ok := artifact.Parts[0].(*DataPart)
		require.True(t, ok)
		assert.Equal(t, "data", dataPart.Kind)
		assert.Equal(t, map[string]any{"key1": "value1", "key2": float64(42)}, dataPart.Data)
		assert.Equal(t, map[string]any{"type": "json"}, dataPart.Metadata)
	})

	t.Run("unmarshal artifact with multiple parts", func(t *testing.T) {
		jsonData := `{
			"artifact_id": "test-artifact-4",
			"name": "multi-part artifact",
			"parts": [
				{
					"kind": "text",
					"text": "First part"
				},
				{
					"kind": "data",
					"data": {"number": 123}
				}
			]
		}`

		var artifact Artifact
		err := json.Unmarshal([]byte(jsonData), &artifact)
		require.NoError(t, err)

		assert.Equal(t, "test-artifact-4", artifact.ArtifactId)
		assert.Equal(t, "multi-part artifact", artifact.Name)
		assert.Len(t, artifact.Parts, 2)

		// First part should be TextPart
		textPart, ok := artifact.Parts[0].(*TextPart)
		require.True(t, ok)
		assert.Equal(t, "text", textPart.Kind)
		assert.Equal(t, "First part", textPart.Text)

		// Second part should be DataPart
		dataPart, ok := artifact.Parts[1].(*DataPart)
		require.True(t, ok)
		assert.Equal(t, "data", dataPart.Kind)
		assert.Equal(t, map[string]any{"number": float64(123)}, dataPart.Data)
	})
}

func TestTaskUnmarshalJSON(t *testing.T) {
	t.Run("unmarshal task with artifacts", func(t *testing.T) {
		jsonData := `{
			"id": "task-123",
			"context_id": "context-456",
			"kind": "test_task",
			"task_status": {
				"state": "completed",
				"time_stamp": "2025-01-01T00:00:00Z"
			},
			"artifacts": [
				{
					"artifact_id": "artifact-1",
					"name": "test artifact",
					"parts": [
						{
							"kind": "text",
							"text": "This is a text part"
						}
					]
				},
				{
					"artifact_id": "artifact-2",
					"name": "file artifact",
					"parts": [
						{
							"kind": "document",
							"file": {
								"bytes": "base64data",
								"mime_type": "text/plain",
								"name": "test.txt"
							}
						}
					]
				}
			]
		}`

		var task Task
		err := json.Unmarshal([]byte(jsonData), &task)
		require.NoError(t, err)

		assert.Equal(t, "task-123", task.Id)
		assert.Equal(t, "context-456", task.ContextId)
		assert.Equal(t, "test_task", task.Kind)
		assert.Equal(t, COMPLETED, task.Status.State)
		assert.Equal(t, "2025-01-01T00:00:00Z", task.Status.TimeStamp)
		assert.Len(t, task.Artifacts, 2)

		// Check first artifact
		assert.Equal(t, "artifact-1", task.Artifacts[0].ArtifactId)
		assert.Equal(t, "test artifact", task.Artifacts[0].Name)
		assert.Len(t, task.Artifacts[0].Parts, 1)

		textPart, ok := task.Artifacts[0].Parts[0].(*TextPart)
		require.True(t, ok)
		assert.Equal(t, "text", textPart.Kind)
		assert.Equal(t, "This is a text part", textPart.Text)

		// Check second artifact
		assert.Equal(t, "artifact-2", task.Artifacts[1].ArtifactId)
		assert.Equal(t, "file artifact", task.Artifacts[1].Name)
		assert.Len(t, task.Artifacts[1].Parts, 1)

		filePart, ok := task.Artifacts[1].Parts[0].(*FilePart)
		require.True(t, ok)
		assert.Equal(t, "document", filePart.Kind)

		fileWithBytes, ok := filePart.File.(*FileWithBytes)
		require.True(t, ok)
		assert.Equal(t, "base64data", fileWithBytes.Bytes)
		assert.Equal(t, "text/plain", fileWithBytes.MimeType)
		assert.Equal(t, "test.txt", fileWithBytes.Name)
	})
}
