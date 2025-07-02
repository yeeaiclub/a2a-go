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
	"github.com/yeeaiclub/a2a-go/sdk/server"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type RequestContext struct {
	TaskId      string
	ContextId   string
	Params      types.MessageSendParam
	Task        *types.Task
	CallContext server.ServerCallContext
}

func NewRequestContext(options ...RequestContextOption) *RequestContext {
	reqContext := &RequestContext{}

	for _, opt := range options {
		opt.Option(reqContext)
	}
	if reqContext.TaskId == "" {
		reqContext.TaskId = uuid.New().String()
	}
	return reqContext
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

func WithServerContext(callContext server.ServerCallContext) RequestContextOption {
	return RequestContextOptionFunc(func(ctx *RequestContext) {
		ctx.CallContext = callContext
	})
}
