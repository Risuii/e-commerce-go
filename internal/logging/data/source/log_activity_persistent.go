package source

import (
	LogModel "e-commerce/internal/logging/data/model"
	Library "e-commerce/library"
)

type LogActivity interface {
	Insert(param *LogModel.LogActivity) error
}

type LogActivityImpl struct {
	library Library.Library
}

func NewLogActivityPersistent(library Library.Library) LogActivity {
	return &LogActivityImpl{
		library: library,
	}
}

func (s *LogActivityImpl) Insert(param *LogModel.LogActivity) error {
	var err error
	// err := s.infinitium.GetConnection().Exec(`
	//     INSERT INTO t_activity_log (
	//         trace_id,
	// 		endpoint,
	//         path,
	//         description,
	//         created_at,
	// 		request_payload
	//     )
	//     VALUES (
	//         ?, ?, ?, ?, ?, ?
	//     )
	// `, param.TraceId, param.Endpoint, param.Path, param.Description, param.CreatedAt, param.RequestPayload).Error
	// if err == nil {
	// 	return nil
	// }
	return err
}
