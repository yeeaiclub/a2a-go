// Copyright 2025 yumosx
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

package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestNewRequestContext(t *testing.T) {
	testcases := []struct {
		name      string
		taskId    string
		contextId string
		params    types.MessageSendParam
		tasks     *types.Task
		want      *RequestContext
	}{
		{
			name:      "new request context",
			taskId:    "1",
			contextId: "2",
			params:    types.MessageSendParam{},
			tasks:     &types.Task{Id: "1", ContextId: "2"},
			want:      &RequestContext{ContextId: "2", TaskId: "1", Params: types.MessageSendParam{}, Task: &types.Task{Id: "1", ContextId: "2"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			context := NewRequestContext(
				WithContextId(tc.contextId),
				WithTaskId(tc.taskId),
				WithTask(tc.tasks),
				WithParams(tc.params),
			)
			assert.Equal(t, context, tc.want)
		})
	}
}
