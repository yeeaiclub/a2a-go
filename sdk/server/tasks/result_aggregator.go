package tasks

import (
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/types"
)

type ResultAggregator struct {
	Manger  event.QueueManger
	Message types.Message
}

func NewResultAggregator(taskManger event.QueueManger) ResultAggregator {
	return ResultAggregator{Manger: taskManger}
}
