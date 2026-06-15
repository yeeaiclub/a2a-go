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


package handler

import (
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/aggregator"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
)

// HandlerOption allows customizing DefaultHandler via functional options.
type HandlerOption interface {
	Option(d *DefaultHandler)
}

// HandlerOptionFunc is a function type for HandlerOption.
type HandlerOptionFunc func(d *DefaultHandler)

func (fn HandlerOptionFunc) Option(d *DefaultHandler) {
	fn(d)
}

// WithTaskManager sets a custom TaskManager for the handler.
func WithTaskManager(taskManger *manager.TaskManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.manager = taskManger
	})
}

// WithQueueManager sets a custom QueueManager for the handler.
func WithQueueManager(queueManger event.QueueManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.queueManger = queueManger
	})
}

// WithResultAggregator sets a custom ResultAggregator for the handler.
func WithResultAggregator(rg *aggregator.ResultAggregator) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.resultAggregator = rg
	})
}

// WithPushNotifier sets a custom PushNotifier for the handler.
func WithPushNotifier(pushNotifier tasks.PushNotifier) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.pushNotifier = pushNotifier
	})
}
