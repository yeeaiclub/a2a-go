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
