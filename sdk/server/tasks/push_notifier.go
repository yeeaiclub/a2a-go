package tasks

import "github.com/yumosx/a2a-go/sdk/types"

type PushNotifier interface {
	SetInfo(taskId string, config types.PushNotificationConfig)
	GetInfo(taskId string) types.PushNotificationConfig
	SendNotification(task types.Task)
}
