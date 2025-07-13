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

package execution

import (
	"github.com/google/uuid"
	"github.com/yeeaiclub/a2a-go/internal/errs"
	"github.com/yeeaiclub/a2a-go/sdk/server"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type RequestContext struct {
	TaskId      string
	ContextId   string
	Params      types.MessageSendParam
	Task        *types.Task
	CallContext *server.CallContext
}

func NewRequestContext(options ...RequestContextOption) (*RequestContext, error) {
	reqContext := &RequestContext{}

	for _, opt := range options {
		opt.Option(reqContext)
	}

	tid := reqContext.TaskId
	cid := reqContext.ContextId
	if tid == "" && reqContext.Params.Message != nil && reqContext.Params.Message.TaskID != "" {
		tid = reqContext.Params.Message.TaskID
	}
	if tid == "" {
		tid = uuid.New().String()
	}
	if cid == "" && reqContext.Params.Message != nil && reqContext.Params.Message.ContextID != "" {
		cid = reqContext.Params.Message.ContextID
	}
	if cid == "" {
		cid = uuid.New().String()
	}

	if reqContext.Params.Message != nil {
		reqContext.Params.Message.TaskID = tid
		reqContext.Params.Message.ContextID = cid
	}

	if reqContext.Task != nil {
		if reqContext.Task.Id != "" && reqContext.Task.Id != tid {
			return nil, errs.ErrBadTaskId
		}
		if reqContext.Task.ContextId != "" && reqContext.Task.ContextId != cid {
			return nil, errs.ErrBadTaskId
		}
	}

	reqContext.TaskId = tid
	reqContext.ContextId = cid

	return reqContext, nil
}

type RequestContextOption interface {
	Option(ctx *RequestContext)
}

type RequestContextOptionFunc func(ctx *RequestContext)

func (fn RequestContextOptionFunc) Option(ctx *RequestContext) {
	fn(ctx)
}

func WithTaskId(taskId string) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.TaskId = taskId
	})
}

func WithContextId(contextId string) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.ContextId = contextId
	})
}

func WithParams(params types.MessageSendParam) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.Params = params
	})
}

func WithTask(task *types.Task) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.Task = task
	})
}

func WithServerContext(callContext *server.CallContext) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.CallContext = callContext
	})
}
