package model

type LogActivity struct {
	TraceId        string `db:"trace_id"`
	Endpoint       string `db:"endpoint"`
	Path           string `db:"path"`
	Description    string `db:"description"`
	CreatedAt      string `db:"created_at"`
	RequestPayload string `db:"request_payload"`
}
