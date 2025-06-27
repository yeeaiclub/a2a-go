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
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
)

// AgentExecutor Implementations of this interface contain the core logic of the agent,
// executing tasks based on requests and publishing updates to an event queue.
type AgentExecutor interface {
	// Execute the agent's logic for a given request context
	Execute(ctx context.Context, requestContext *RequestContext, queue *event.Queue) error
	// Cancel request the agent to cancel an ongoing task
	Cancel(ctx context.Context, requestContext *RequestContext, queue *event.Queue) error
}
