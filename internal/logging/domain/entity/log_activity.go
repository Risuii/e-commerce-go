package entity

import (
	LogModel "e-commerce/internal/logging/data/model"
)

type LogActivity struct {
	TraceID         string
	Endpoint        string
	Path            string
	Description     string
	CreatedAt       string
	RequestPayload  string
	ResponsePayload string
}

func (e LogActivity) ToModel() *LogModel.LogActivity {
	result := LogModel.LogActivity{
		TraceID:         e.TraceID,
		Endpoint:        e.Endpoint,
		Path:            e.Path,
		Description:     e.Description,
		CreatedAt:       e.CreatedAt,
		RequestPayload:  e.RequestPayload,
		ResponsePayload: e.ResponsePayload,
	}

	return &result
}
