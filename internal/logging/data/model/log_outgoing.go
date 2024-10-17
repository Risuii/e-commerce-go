package model

type LogOutgoing struct {
	TraceId         string `db:"trace_id"`
	BackendSystem   string `db:"backend_system"`
	ServiceName     string `db:"service_name"`
	RequestPayload  string `db:"request_payload"`
	ResponsePayload string `db:"response_payload"`
	RequestDate     string `db:"request_date"`
	ResponseDate    string `db:"response_date"`
	StatusCode      string `db:"status_code"`
	MsName          string `db:"ms_name"`
	ServerNode      string `db:"server_node"`
	SourceNode      string `db:"source_node"`
	ResponseTime    string `db:"response_time"`
}
