package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageUnmarshalJSONWithTextPart(t *testing.T) {
	t.Run("unmarshalJSON with text part", func(t *testing.T) {
		jsonData := `
{
  "role": "user",
  "task_id": "123",
  "context_id": "ctx456",
  "parts": [
    {
      "kind": "text",
      "text": "Hello, world!"
    }
  ],
  "metadata": {
    "source": "web"
  }
}
`
		var msg Message
		err := json.Unmarshal([]byte(jsonData), &msg)
		require.NoError(t, err)

		assert.Equal(t, User, msg.Role)
		assert.Equal(t, "123", msg.TaskID)
		assert.Equal(t, "ctx456", msg.ContextID)
		assert.Len(t, msg.Parts, 1)
		assert.IsType(t, &TextPart{}, msg.Parts[0])
		textPart := msg.Parts[0].(*TextPart)
		assert.Equal(t, "Hello, world!", textPart.Text)
		assert.Equal(t, "web", msg.Metadata["source"])
		assert.Equal(t, "ctx456", msg.GetContextId())
		assert.Equal(t, "123", msg.GetTaskId())
	})
}

func TestMessageUnmarshalJSONWithDataPart(t *testing.T) {
	t.Run("MessageUnmarshalJSONWithDataPart", func(t *testing.T) {
		jsonData := `
{
  "role": "agent",
  "parts": [
    {
      "kind": "data",
      "data": {"1": "a"}
    }
  ]
}
`
		var msg Message
		err := json.Unmarshal([]byte(jsonData), &msg)
		require.NoError(t, err)

		assert.Equal(t, Agent, msg.Role)
		assert.Len(t, msg.Parts, 1)
		assert.IsType(t, &DataPart{}, msg.Parts[0])

		dataPart := msg.Parts[0].(*DataPart)
		assert.Equal(t, map[string]any{"1": "a"}, dataPart.Data)
	})
}

func TestMessageUnmarshalJSONWithFilePart(t *testing.T) {
	t.Run("MessageUnmarshalJSONWithFilePart", func(t *testing.T) {
		jsonData := `
{
  "role": "user",
  "parts": [
    {
      "kind": "file",
      "file": {
        "bytes": "base64encodeddata"
      },
      "metadata": {
        "name": "test.txt"
      }
    }
  ]
}
`
		var msg Message
		err := json.Unmarshal([]byte(jsonData), &msg)
		require.NoError(t, err)

		assert.Len(t, msg.Parts, 1)
		filePart := msg.Parts[0].(*FilePart)
		assert.IsType(t, &FileWithBytes{}, filePart.File)
		bytesVal := filePart.File.(*FileWithBytes)
		assert.Equal(t, "base64encodeddata", bytesVal.Bytes)
		assert.Equal(t, "test.txt", filePart.Metadata["name"])
	})
}
