package source

import (
	LogModel "e-commerce/internal/logging/data/model"
	Library "e-commerce/library"
)

type LogInterface interface {
	Insert(param *LogModel.LogInterface) error
}

type LogInterfaceImpl struct {
	library Library.Library
}

func NewLogInterfacePersistent(library Library.Library) LogInterface {
	return &LogInterfaceImpl{
		library: library,
	}
}

func (s *LogInterfaceImpl) Insert(param *LogModel.LogInterface) error {
	var err error
	// err := s.infinitium.GetConnection().Exec(`
	// 	INSERT INTO t_interface_log (
	// 		trace_id,
	// 		server_node,
	// 		service_name,
	// 		request_payload,
	// 		response_payload,
	// 		request_date,
	// 		response_date,
	// 		source_name,
	// 		ms_name
	// 	)
	// 	VALUES (
	// 		?, ?, ?, ?, ?, ?, ?, ?, ?
	// 	)
	// `, param.TraceId, param.ServerNode, param.ServiceName, param.RequestPayload, param.ResponsePayload, param.RequestDate, param.ResponseDate, param.SourceName, param.MsName).Error

	// if err == nil {
	// 	return nil
	// }

	return err
}
