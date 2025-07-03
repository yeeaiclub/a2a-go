package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSONWithBytes(t *testing.T) {
	t.Run("unmarshal with bytes", func(t *testing.T) {
		jsonData := `
{
  "file": { "bytes": "base64encodeddata" },
  "kind": "document",
  "metadata": { "name": "test.txt" }
}
`
		var part FilePart
		err := json.Unmarshal([]byte(jsonData), &part)
		require.NoError(t, err)

		assert.Equal(t, "document", part.Kind)
		assert.Equal(t, map[string]any{"name": "test.txt"}, part.Metadata)
		bytesVal := part.File.(*FileWithBytes)
		assert.Equal(t, "base64encodeddata", bytesVal.Bytes)
		assert.Equal(t, "file", part.GetKind())
		assert.Equal(t, part.Metadata, part.GetMetadata())
	})
}

func TestUnmarshalJSONWithUrl(t *testing.T) {
	t.Run("unmarshal with url", func(t *testing.T) {
		jsonData := `
{
  "file": { "url": "https://example.com/file.txt " },
  "kind": "image",
  "metadata": { "size": 12345 }
}
`
		var part FilePart
		err := json.Unmarshal([]byte(jsonData), &part)
		require.NoError(t, err)

		assert.Equal(t, "image", part.Kind)
		assert.Equal(t, map[string]any{"size": float64(12345)}, part.Metadata)
		urlVal := part.File.(*FileWithUrl)
		assert.Equal(t, "https://example.com/file.txt ", urlVal.Url)
	})
}
