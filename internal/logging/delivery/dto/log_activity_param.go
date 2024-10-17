package dto

import "time"

type LogActivityParam struct {
	LogId       string    `json:"log_id"`
	Endpoint    string    `json:"endpoint"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"creadted_at"`
}
