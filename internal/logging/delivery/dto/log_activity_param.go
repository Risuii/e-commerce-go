package dto

type LogActivityParam struct {
	TraceID string      `json:"traceID,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Path    string      `json:"path,omitempty"`
	Status  bool        `json:"status,omitempty"`
	Token   string      `json:"token,omitempty"`
}
