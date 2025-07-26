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

package aggregator

import (
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
)

// ResultAggregator is used to process the event streams from an AgentExecutor
type ResultAggregator struct {
	manager   *manager.TaskManager
	batchSize int
}

type ResultAggregatorOption interface {
	Option(rg *ResultAggregator)
}

func NewResultAggregator(taskManger *manager.TaskManager) *ResultAggregator {
	rg := &ResultAggregator{manager: taskManger, batchSize: 10}
	return rg
}

func (r *ResultAggregator) WithBatchSize(batchSize int) *ResultAggregator {
	r.batchSize = batchSize
	return r
}

func (r *ResultAggregator) BuildStreaming() *StreamingConsumer {
	return NewStreamingAggregator(r.manager, r.batchSize)
}

func (r *ResultAggregator) BuildFull() *FullConsumer {
	return NewFullConsumer(r.manager)
}

func (r *ResultAggregator) BuildInterruptible() *InterruptibleConsumer {
	return NewInterruptibleConsumer(r.manager)
}
