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
