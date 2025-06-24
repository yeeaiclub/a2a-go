package tasks

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

// PushNotifier interface to store, retrieve push notification for tasks
// and send push notifications.
type PushNotifier interface {
	// SetInfo set of update the push notification configuration for a task
	SetInfo(ctx context.Context, taskId string, config *types.PushNotificationConfig) error
	// GetInfo retrieves the push
	GetInfo(ctx context.Context, taskId string) (*types.PushNotificationConfig, error)
	// Delete the push notification configuration for a task
	Delete(ctx context.Context, taskId string) error
	// SendNotification sends a push notification containing the latest task state
	SendNotification(task *types.Task) error
}
