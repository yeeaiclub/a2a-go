// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

const Version = "2.0"

type Role string

const (
	Agent Role = "agent"
	User  Role = "user"
)

type PushNotificationConfig struct {
	Id             string                             `json:"id,omitempty"`
	URL            string                             `json:"url,omitempty"`
	Authentication PushNotificationAuthenticationInfo `json:"authentication,omitempty"`
}

type PushNotificationAuthenticationInfo struct {
	Credentials string   `json:"credentials,omitempty"`
	Schemes     []string `json:"schemes,omitempty"`
}

type TaskPushNotificationConfig struct {
	TaskId string                  `json:"task_id,omitempty"`
	Config *PushNotificationConfig `json:"config"`
}

type TaskIdParams struct {
	Id string `json:"id,omitempty"`
}

type CancelTaskResponse struct {
	Id      string `json:"id,omitempty"`
	JSONRPC string `json:"jsonrpc,omitempty"`
}

type SendStreamingMessageRequest struct {
}

type Artifact struct {
	ArtifactId  string         `json:"artifact_id,omitempty"`
	Description string         `json:"description,omitempty"`
	Extensions  []string       `json:"extension,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Name        string         `json:"name,omitempty"`
	Parts       []Part         `json:"parts,omitempty"`
}

type Result[T any] struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

type JSONRPCErrorResponse[T any] struct {
	Error   Result[T] `json:"error,omitempty"`
	Id      string    `json:"id,omitempty"`
	JSONRPC string    `json:"jsonrpc,omitempty"`
}

func InternalError() Result[string] {
	return Result[string]{
		Code:    -32603,
		Message: "Internal error",
		Data:    "",
	}
}

func ContentTypeNotSupportedError[T any](data T) Result[T] {
	return Result[T]{
		Code:    -32005,
		Message: "Incompatible content types",
		Data:    data,
	}
}

func JSONParseError[T any](data T) Result[T] {
	return Result[T]{
		Code:    -32700,
		Data:    data,
		Message: "Invalid JSON payload",
	}
}

func MethodNotFoundError[T any](data T) Result[T] {
	return Result[T]{
		Code:    -32601,
		Data:    data,
		Message: "Method not found",
	}
}
