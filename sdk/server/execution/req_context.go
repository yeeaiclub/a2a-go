package execution

import "github.com/yumosx/a2a-go/sdk/types"

type RequestContext struct {
	TaskId    string
	ContextId string
	Params    types.MessageSendParam
	Tasks     []types.Task
}
