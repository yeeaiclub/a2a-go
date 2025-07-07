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

package types

type AgentCard struct {
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	URL                string                    `json:"url,omitempty"`
	Skills             []AgentSkill              `json:"skills"`
	DefaultInputModes  []string                  `json:"default_input_modes,omitempty"`
	DefaultOutputModes []string                  `json:"default_output_modes,omitempty"`
	Provider           *AgentProvider            `json:"provider,omitempty"`
	Capabilities       *AgentCapabilities        `json:"capabilities,omitempty"`
	Version            string                    `json:"version"`
	IconUrl            string                    `json:"icon_url,omitempty"`
	Security           []map[string][]string     `json:"security,omitempty"`
	SecuritySchemes    map[string]SecurityScheme `json:"security_schemes,omitempty"`
}

type AgentProvider struct {
	Organization string `json:"organization"`
	URL          string `json:"url,omitempty"`
}

type AgentSkill struct {
	Id          string   `json:"id,omitempty"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
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

type SendMessageRequest struct {
	Id     string           `json:"id,omitempty"`
	Method string           `json:"method"`
	Params MessageSendParam `json:"params"`
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

type PushNotificationConfig struct {
	Id             string                              `json:"id,omitempty"`
	URL            string                              `json:"url,omitempty"`
	Authentication *PushNotificationAuthenticationInfo `json:"authentication,omitempty"`
}

type PushNotificationAuthenticationInfo struct {
	Credentials string   `json:"credentials,omitempty"`
	Schemes     []string `json:"schemes,omitempty"`
}

type SendStreamingMessageRequest struct {
	Id      string           `json:"id,omitempty"`
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  MessageSendParam `json:"params"`
}

type GetTaskRequest struct {
	Id      string          `json:"id,omitempty"`
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  TaskQueryParams `json:"params"`
}

type TaskQueryParams struct {
	Id            string         `json:"id"`
	HistoryLength int            `json:"history_length,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

type CancelTaskRequest struct {
	Id     string       `json:"id,omitempty"`
	Method string       `json:"method"`
	Params TaskIdParams `json:"params"`
}

type TaskIdParams struct {
	Id       string         `json:"id,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type SetTaskPushNotificationConfigRequest struct {
	Id      string                     `json:"id,omitempty"`
	JSONRPC string                     `json:"jsonrpc"`
	Method  string                     `json:"method"`
	Params  TaskPushNotificationConfig `json:"params"`
}

type TaskPushNotificationConfig struct {
	TaskId string                  `json:"task_id,omitempty"`
	Config *PushNotificationConfig `json:"config,omitempty"`
}

type GetTaskPushNotificationConfigRequest struct {
	Id      string       `json:"id,omitempty"`
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  TaskIdParams `json:"params"`
}

type TaskResubscriptionRequest struct {
	Id      string       `json:"id,omitempty"`
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  TaskIdParams `json:"params"`
}
