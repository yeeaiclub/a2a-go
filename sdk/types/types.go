package types

const Version = "2.0"

type In string

const (
	InCookie In = "cookie"
	InHeader    = "header"
	InQuery     = "query"
)

// APIKeySecurityScheme API key security scheme
type APIKeySecurityScheme struct {
	Description string `json:"description,omitempty"`
	In          In     `json:"in,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type AgentExtension struct {
	Description string         `json:"description,omitempty"`
	Params      map[string]any `json:"params,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Url         string         `json:"url,omitempty"`
}

type AgentProvider struct {
	Organization string `json:"organization,omitempty"`
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

type AuthorizationCodeOAuthFlow struct {
	AuthorizationUrl string            `json:"authorization_url,omitempty"`
	RefreshUrl       string            `json:"refresh_url,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
	TokenUrl         string            `json:"token_url,omitempty"`
}

type DataPart struct {
	Data     map[string]any `json:"data,omitempty"`
	Kind     string         `json:"kind,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type FileBase struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
}

type FileWithBytes struct {
	Bytes    string `json:"bytes,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
}

type FileWithUrl struct {
	MimeType string `json:"mime_type,omitempty"`
	Name     string `json:"name,omitempty"`
	Url      string `json:"url,omitempty"`
}

type HTTPAuthSecurityScheme struct {
}

type ClientCredentialsOAuthFlow struct {
	RefreshUrl string            `json:"refresh_url,omitempty"`
	Scopes     map[string]string `json:"scopes,omitempty"`
	TokenUrl   string            `json:"token_url,omitempty"`
}

type AgentCard struct {
	Name               string            `json:"name,omitempty"`
	DefaultInputModes  []string          `json:"default_input_modes,omitempty"`
	DefaultOutputModes []string          `json:"default_output_modes,omitempty"`
	Description        string            `json:"description,omitempty"`
	IconUrl            string            `json:"icon_url,omitempty"`
	URL                string            `json:"url,omitempty"`
	Provider           AgentProvider     `json:"provider,omitempty"`
	Skills             []AgentSkill      `json:"skills,omitempty"`
	Capabilities       AgentCapabilities `json:"capabilities,omitempty"`
	Version            string            `json:"version,omitempty"`
}

type AgentCapabilities struct {
	Extensions             []AgentExtension `json:"extensions,omitempty"`
	StateTransitionHistory bool             `json:"stateTransitionHistory,omitempty"`
	Streaming              bool             `json:"streaming,omitempty"`
}

type Message struct {
	ContextID        string         `json:"context_id,omitempty"`
	Extensions       []string       `json:"extensions,omitempty"`
	Kind             string         `json:"kind,omitempty"`
	MessageID        string         `json:"message_id,omitempty"`
	ReferenceTaskIDs []string       `json:"referenceTaskIds,omitempty"`
	Role             string         `json:"role,omitempty"`
	TaskID           string         `json:"taskId,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

type Task struct {
	Id        string         `json:"id,omitempty"`
	ContextId string         `json:"context_id,omitempty"`
	History   []Message      `json:"history,omitempty"`
	Kind      string         `json:"kind,omitempty"`
	Status    TaskStatues    `json:"task_status,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type TaskStatues struct {
	Message   Message    `json:"id,omitempty"`
	State     TaskStatue `json:"state,omitempty"`
	TimeStamp string     `json:"time_stamp,omitempty"`
}

type TaskStatue string

const (
	SUBMITTED      = "submitted"
	WORKING        = "working"
	INPUT_REQUIRED = "required"
	COMPLETED      = "completed"
	CANCELED       = "canceled"
	FAILED         = "failed"
	UNKNOWN        = "unknown"
)

type SendMessageRequest struct {
	Id     string           `json:"id,omitempty"`
	Method string           `json:"method,omitempty"`
	Params MessageSendParam `json:"params,omitempty"`
}

type MessageSendParam struct {
	Configuration MessageSendConfiguration `json:"configuration,omitempty"`
	Message       Message                  `json:"message,omitempty"`
	Metadata      map[string]any           `json:"metadata,omitempty"`
}

type MessageSendConfiguration struct {
	AcceptedOutputModes    []string               `json:"accepted_output_modes,omitempty"`
	Blocking               bool                   `json:"blocking,omitempty"`
	HistoryLength          int                    `json:"history_length,omitempty"`
	PushNotificationConfig PushNotificationConfig `json:"push_notification_config,omitempty"`
}

type PushNotificationAuthenticationInfo struct {
	Credentials string   `json:"credentials,omitempty"`
	Schemes     []string `json:"schemes,omitempty"`
}

type PushNotificationConfig struct {
	Id             string                             `json:"id,omitempty"`
	URL            string                             `json:"url,omitempty"`
	Authentication PushNotificationAuthenticationInfo `json:"authentication,omitempty"`
}

type TaskPushNotificationConfig struct {
	TaskId string `json:"task_id,omitempty"`
}

type SendMessageResponse struct {
	Id string `json:"id,omitempty"`
}

type TaskIdParams struct {
	Id string `json:"id,omitempty"`
}

type GetTaskRequest struct {
	Id      string       `json:"id,omitempty"`
	JSONRPC string       `json:"jsonrpc,omitempty"`
	Params  TaskIdParams `json:"params,omitempty"`
}

type GetTaskResponse struct {
}

type CancelTaskRequest struct {
	Id string
}

type CancelTaskResponse struct {
}

type SendStreamingMessageRequest struct {
}

type TaskStatusUpdateEvent struct {
	ContextId string `json:"context_id,omitempty"`
	Final     bool   `json:"final,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Metadata  map[string]any
	Status    TaskStatue
}

type TaskQueryParam struct {
	HistoryLength int    `json:"history_length,omitempty"`
	Id            string `json:"id,omitempty"`
}

type TaskArtifactUpdateEvent struct {
	TaskId    string `json:"task_id,omitempty"`
	ContextId string `json:"context_id,omitempty"`
	Kind      string `json:"kind,omitempty"`
}

type Artifact struct {
}

type SetTaskPushNotificationConfigRequest struct {
	Id string
}

type SetTaskPushNotificationConfigSuccessResponse struct {
	Id     string `json:"id,omitempty"`
	Result TaskPushNotificationConfig
}

type GetTaskPushNotificationConfigRequest struct {
	Id string `json:"id,omitempty"`
}

type TaskPushNotificationConfigRequest struct {
	Id     string       `json:"id,omitempty"`
	Params TaskIdParams `json:"params,omitempty"`
}

type TaskResubscriptionRequest struct {
	Id     string       `json:"id,omitempty"`
	Params TaskIdParams `json:"params,omitempty"`
}

type JSONRPCErrorResponse[T any] struct {
	Error   Result[T] `json:"error,omitempty"`
	Id      string    `json:"id,omitempty"`
	JSONRPC string    `json:"jsonrpc,omitempty"`
}

type Result[T any] struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
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
