package source

import (
	LogModel "e-commerce/internal/logging/data/model"
	Library "e-commerce/library"

	EcommercePackage "e-commerce/pkg/data_sources/e-commerce"
)

type LogActivity interface {
	Insert(param *LogModel.LogActivity) error
}

type LogActivityImpl struct {
	library   Library.Library
	ecommerce EcommercePackage.Ecommerce
}

func NewLogActivityPersistent(library Library.Library, ecommerce EcommercePackage.Ecommerce) LogActivity {
	return &LogActivityImpl{
		library:   library,
		ecommerce: ecommerce,
	}
}

func (s *LogActivityImpl) Insert(param *LogModel.LogActivity) error {

	err := s.ecommerce.GetConnection().Exec(`
	    INSERT INTO activity_log (
	        trace_id,
			endpoint,
	        path,
	        description,
	        created_at,
			request_payload,
			response_payload
	    )
	    VALUES (
	        ?, ?, ?, ?, ?, ?, ?
	    )
	`, param.TraceID, param.Endpoint, param.Path, param.Description, param.CreatedAt, param.RequestPayload, param.ResponsePayload).Error
	if err == nil {
		return nil
	}
	return err
}
