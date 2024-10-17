package model

type LogInterface struct {
	TraceId         string `db:"trace_id"`
	ServerNode      string `db:"server_node"`
	ServiceName     string `db:"service_name"`
	RequestPayload  string `db:"request_payload"`
	RequestDate     string `db:"request_date"`
	ResponsePayload string `db:"response_payload"`
	ResponseDate    string `db:"response_date"`
	SourceName      string `db:"source_name"`
	MsName          string `db:"ms_name"`
}
