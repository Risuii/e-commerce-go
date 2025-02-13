package model

type LogActivity struct {
	TraceID         string `db:"trace_id"`
	Endpoint        string `db:"endpoint"`
	Path            string `db:"path"`
	Description     string `db:"description"`
	CreatedAt       string `db:"created_at"`
	RequestPayload  string `db:"request_payload"`
	ResponsePayload string `db:"response_payload"`
}
