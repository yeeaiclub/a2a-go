package handler

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

type MessageHandler interface {
	OnMessageSend(ctx context.Context, params types.MessageSendParam) (types.Message, error)
	OnMessageSendStream(ctx context.Context, params types.MessageSendParam) (chan types.Event, error)
}

type TaskHandler interface {
	OnGetTask(ctx context.Context, params types.TaskQueryParam) (types.Task, error)
	OnCancelTask(ctx context.Context, params types.TaskIdParams) (types.Task, error)
	OnSetTaskPushNotificationConfig(ctx context.Context, params types.TaskPushNotificationConfig) (types.TaskPushNotificationConfig, error)
	OnGetTaskPushNotificationConfig(ctx context.Context, params types.TaskIdParams) (types.TaskPushNotificationConfig, error)
	OnResubscribeToTask(ctx context.Context, params types.TaskIdParams) (chan types.Event, error)
}

type Handler interface {
	MessageRequestHandler
	TaskHandler
}
