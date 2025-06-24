package types

type Part interface {
	GetPart() string
}

// DataPart Represents a structured data segment within a message part.
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

type TextPart struct {
	Kind     string         `json:"kind"`
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (t *TextPart) GetPart() string {
	return "text"
}
