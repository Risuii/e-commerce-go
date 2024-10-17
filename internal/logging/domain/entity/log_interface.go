package entity

import (
	LogModel "e-commerce/internal/logging/data/model"
)

type LogInterface struct {
	TraceId         string
	ServerNode      string
	ServiceName     string
	RequestPayload  string
	RequestDate     string
	ResponsePayload string
	ResponseDate    string
	SourceName      string
	MsName          string
}

func (e LogInterface) ToModel() *LogModel.LogInterface {
	result := LogModel.LogInterface{
		TraceId:         e.TraceId,
		ServerNode:      e.ServerNode,
		ServiceName:     e.ServiceName,
		RequestPayload:  e.RequestPayload,
		RequestDate:     e.RequestDate,
		ResponsePayload: e.ResponsePayload,
		ResponseDate:    e.ResponseDate,
		SourceName:      e.SourceName,
		MsName:          e.MsName,
	}

	return &result
}
