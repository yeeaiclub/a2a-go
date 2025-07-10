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

package jsonx

import (
	"encoding/json"
	"fmt"
)

type KindUnmarshaler interface {
	GetKind() string
}

func UnmarshalByKind[T KindUnmarshaler](data []byte, kindMap map[string]func() T) (T, error) {
	var zero T
	var kindHolder struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &kindHolder); err != nil {
		return zero, fmt.Errorf("failed to unmarshal %s: %w", kindHolder.Kind, err)
	}

	constructor, exists := kindMap[kindHolder.Kind]
	if !exists {
		return zero, fmt.Errorf("unknown kind: %q", kindHolder.Kind)
	}
	result := constructor()
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal %s: %w", kindHolder.Kind, err)
	}
	return result, nil
}

func UnmarshalSliceByKind[T KindUnmarshaler](data []json.RawMessage, kindMap map[string]func() T) ([]T, error) {
	results := make([]T, 0, len(data))
	for i, raw := range data {
		result, err := UnmarshalByKind(raw, kindMap)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item %d: %w", i, err)
		}
		results = append(results, result)
	}
	return results, nil
}
