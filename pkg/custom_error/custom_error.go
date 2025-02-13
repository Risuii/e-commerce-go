package custom_error

import (
	"errors"
	"net/http"

	Constants "e-commerce/constants"
	Library "e-commerce/library"
)

type CustomError struct {
	display        error
	plain          error
	path           string
	errValidations []map[string]interface{}
	library        Library.Library
	code           int
}

func New(
	display error,
	plain error,
	path string,
	library Library.Library,
) error {
	return &CustomError{
		display: display,
		plain:   plain,
		path:    path,
		library: library,
		code:    http.StatusInternalServerError,
	}
}

func (e *CustomError) Error() string {
	message := map[string]interface{}{
		"display": e.display.Error(),
		"plain":   e.plain.Error(),
		"path":    e.path,
	}

	result, err := e.library.JsonMarshal(message)
	if err != nil {
		return err.Error()
	}

	return string(result)
}

func (e *CustomError) GetDisplay() error {
	return e.display
}

func (e *CustomError) GetPlain() error {
	return e.plain
}

func (e *CustomError) GetPath() string {
	return e.path
}

func (e *CustomError) UnshiftPath(path string) error {
	e.path = path + " > " + e.path
	return e
}

func (e *CustomError) FromListMap(errs []map[string]interface{}) error {
	result, err := e.library.JsonMarshal(errs)
	if err != nil {
		return New(
			Constants.ErrFailedJSONMarshal,
			err,
			"CustomError:FromListMap",
			e.library,
		)
	}

	e.errValidations = errs
	e.plain = errors.New(string(result))
	return e
}

func (e *CustomError) GetErrorValidations() []map[string]interface{} {
	if len(e.errValidations) == 0 {
		return nil
	}
	return e.errValidations
}

func (e *CustomError) SetCode(code int) error {
	e.code = code

	return e
}

func (e *CustomError) GetCode() int {
	return e.code
}
