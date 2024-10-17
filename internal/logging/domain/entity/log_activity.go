package entity

import (
	LogModel "e-commerce/internal/logging/data/model"
)

type LogActivity struct {
	TraceId        string
	Endpoint       string
	Path           string
	Description    string
	CreatedAt      string
	RequestPayload string
}

func (e LogActivity) ToModel() *LogModel.LogActivity {
	result := LogModel.LogActivity{
		TraceId:        e.TraceId,
		Endpoint:       e.Endpoint,
		Path:           e.Path,
		Description:    e.Description,
		CreatedAt:      e.CreatedAt,
		RequestPayload: e.RequestPayload,
	}

	return &result
}
