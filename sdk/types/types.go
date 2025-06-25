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

type AgentCard struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	URL                string             `json:"url,omitempty"`
	Skills             []AgentSkill       `json:"skills"`
	DefaultInputModes  []string           `json:"default_input_modes,omitempty"`
	DefaultOutputModes []string           `json:"default_output_modes,omitempty"`
	Provider           *AgentProvider     `json:"provider,omitempty"`
	Capabilities       *AgentCapabilities `json:"capabilities,omitempty"`
	Version            string             `json:"version"`
	IconUrl            string             `json:"icon_url,omitempty"`
}

type AgentProvider struct {
	Organization string `json:"organization"`
	URL          string `json:"url,omitempty"`
}

type AgentSkill struct {
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
	Id          string   `json:"id,omitempty"`
	InputModes  []string `json:"input_modes,omitempty"`
	Name        string   `json:"name,omitempty"`
	OutputModes []string `json:"output_modes,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type AgentCapabilities struct {
	Extensions             []AgentExtension `json:"extensions,omitempty"`
	StateTransitionHistory bool             `json:"stateTransitionHistory,omitempty"`
	Streaming              bool             `json:"streaming,omitempty"`
}

type AgentExtension struct {
	URL         string         `json:"url"`
	Description string         `json:"description,omitempty"`
	Params      map[string]any `json:"params,omitempty"`
	Required    bool           `json:"required,omitempty"`
}

type AuthorizationCodeOAuthFlow struct {
	AuthorizationUrl string            `json:"authorization_url,omitempty"`
	RefreshUrl       string            `json:"refresh_url,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
	TokenUrl         string            `json:"token_url,omitempty"`
}

type SendMessageRequest struct {
	Id     string            `json:"id,omitempty"`
	Method string            `json:"method,omitempty"`
	Params *MessageSendParam `json:"params,omitempty"`
}

type MessageSendParam struct {
	Configuration *MessageSendConfiguration `json:"configuration,omitempty"`
	Message       *Message                  `json:"message,omitempty"`
	Metadata      map[string]any            `json:"metadata,omitempty"`
}

type MessageSendConfiguration struct {
	AcceptedOutputModes    []string                `json:"accepted_output_modes,omitempty"`
	Blocking               bool                    `json:"blocking,omitempty"`
	HistoryLength          int                     `json:"history_length,omitempty"`
	PushNotificationConfig *PushNotificationConfig `json:"push_notification_config,omitempty"`
}

type SendMessageResponse struct {
	Id     string `json:"id,omitempty"`
	Result Event  `json:"result,omitempty"`
}

type GetTaskRequest struct {
	Id      string          `json:"id,omitempty"`
	JSONRPC string          `json:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  TaskQueryParams `json:"params,omitempty"`
}

type GetTaskSuccessResponse struct {
	Id      string `json:"id,omitempty"`
	JSONRPC string `json:"jsonrpc,omitempty"`
	Result  Task   `json:"result,omitempty"`
}

type TaskQueryParams struct {
	Id            string         `json:"id"`
	HistoryLength int            `json:"history_length,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

type CancelTaskRequest struct {
	Id     string       `json:"id,omitempty"`
	Params TaskIdParams `json:"params,omitempty"`
}

type SetTaskPushNotificationConfigRequest struct {
	Id      string                     `json:"id,omitempty"`
	JSONRPC string                     `json:"jsonrpc,omitempty"`
	Method  string                     `json:"method,omitempty"`
	Params  TaskPushNotificationConfig `json:"params,omitempty"`
}

type SetTaskPushNotificationConfigSuccessResponse struct {
	Id      string                     `json:"id,omitempty"`
	JSONRPC string                     `json:"jsonrpc,omitempty"`
	Result  TaskPushNotificationConfig `json:"result,omitempty"`
}

type GetTaskPushNotificationConfigRequest struct {
	Id      string       `json:"id,omitempty"`
	JSONRPC string       `json:"jsonrpc,omitempty"`
	Method  string       `json:"method,omitempty"`
	Params  TaskIdParams `json:"params,omitempty"`
}

type GetTaskPushNotificationConfigSuccessResponse struct {
	Id      string                     `json:"id,omitempty"`
	JSONRPC string                     `json:"jsonrpc"`
	Result  TaskPushNotificationConfig `json:"result"`
}

type TaskResubscriptionRequest struct {
	Id      string       `json:"id,omitempty"`
	JSONRPC string       `json:"jsonrpc,omitempty"`
	Method  string       `json:"method"`
	Params  TaskIdParams `json:"params,omitempty"`
}
