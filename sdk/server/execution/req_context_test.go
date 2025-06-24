package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yumosx/a2a-go/sdk/types"
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
