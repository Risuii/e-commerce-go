package dto

type LogActivityParam struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Path    string      `json:"path,omitempty"`
	Status  bool        `json:"status,omitempty"`
	TraceID string      `json:"traceID,omitempty"`
}
