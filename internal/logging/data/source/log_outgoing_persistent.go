package source

import (
	LogModel "e-commerce/internal/logging/data/model"
	Library "e-commerce/library"
)

type LogOutgoing interface {
	Insert(param *LogModel.LogOutgoing) error
}

type LogOutgoingImpl struct {
	library Library.Library
}

func NewLogOutgoingPersistent(library Library.Library) LogOutgoing {
	return &LogOutgoingImpl{
		library: library,
	}
}

func (s *LogOutgoingImpl) Insert(param *LogModel.LogOutgoing) error {
	var err error
	// err := s.infinitium.GetConnection().Exec(`
	// 	INSERT INTO t_outgoing_log (
	// 		trace_id,
	// 		backend_system,
	// 		service_name,
	// 		request_payload,
	// 		response_payload,
	// 		request_date,
	// 		response_date,
	// 		status_code,
	// 		ms_name,
	// 		server_node,
	// 		source_node,
	// 		response_time
	// ) VALUES (
	// 		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	// )
	// `, param.TraceId, param.BackendSystem, param.ServiceName, param.RequestPayload, param.ResponsePayload, param.RequestDate, param.ResponseDate, param.StatusCode, param.MsName, param.ServerNode, param.SourceNode, param.ResponseTime).Error

	// if err == nil {
	// 	return nil
	// }

	return err
}
