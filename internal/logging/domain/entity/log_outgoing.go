package entity

import (
	LogModel "e-commerce/internal/logging/data/model"
)

type LogOutgoing struct {
	TraceId         string
	BackendSystem   string
	ServiceName     string
	RequestPayload  string
	ResponsePayload string
	RequestDate     string
	ResponseDate    string
	StatusCode      string
	MsName          string
	ServerNode      string
	SourceNode      string
	ResponseTime    string
}

func (e LogOutgoing) ToModel() *LogModel.LogOutgoing {
	result := LogModel.LogOutgoing{
		TraceId:         e.TraceId,
		BackendSystem:   e.BackendSystem,
		ServiceName:     e.ServiceName,
		RequestPayload:  e.RequestPayload,
		ResponsePayload: e.ResponsePayload,
		RequestDate:     e.RequestDate,
		ResponseDate:    e.ResponseDate,
		StatusCode:      e.StatusCode,
		MsName:          e.MsName,
		ServerNode:      e.ServerNode,
		SourceNode:      e.SourceNode,
		ResponseTime:    e.ResponseTime,
	}

	return &result
}
