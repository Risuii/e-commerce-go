package request_information

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	Constants "e-commerce/constants"
)

type RequestInformation struct {
	Path        string
	Method      string
	Header      []byte
	RequestBody []byte
	Params      []byte
}

func (r *RequestInformation) GetRequestInformation(c *gin.Context) map[string]interface{} {
	r.Path = c.Request.URL.Path
	r.Method = c.Request.Method
	contentType := c.ContentType()

	// GET HEADERS
	headerChannel := make(chan []byte)
	go func(c *gin.Context, resultChannel chan []byte) {
		var resultMap = make(map[string]interface{})
		for k, v := range c.Request.Header {
			resultMap[k] = v
			if len(v) == 1 {
				resultMap[k] = v[0]
			}
		}
		json, err := json.Marshal(resultMap)
		if err != nil {
			resultChannel <- nil
			close(resultChannel)
			return
		}
		resultChannel <- json
		close(resultChannel)
	}(c, headerChannel)

	// GET REQUEST BODY
	responseBodyChannel := make(chan []byte)
	go func(c *gin.Context, resultChannel chan []byte) {
		var resultMap map[string]interface{}
		var err error
		switch contentType {
		case Constants.MultiPartFormData:
			err = c.ShouldBind(&resultMap)
		case Constants.TextPlain:
			body, readErr := io.ReadAll(c.Request.Body)
			if readErr != nil {
				err = readErr
			} else {
				resultMap = map[string]interface{}{Constants.Body: string(body)}
			}
		default:
			err = c.ShouldBindBodyWith(&resultMap, binding.JSON)
		}

		if err != nil {
			resultChannel <- nil
			close(resultChannel)
			return
		}
		json, err := json.Marshal(resultMap)
		if err != nil {
			resultChannel <- nil
			close(resultChannel)
			return
		}
		resultChannel <- json
		close(resultChannel)
	}(c, responseBodyChannel)

	// GET PARAMS
	responseParamChannel := make(chan []byte)
	go func(c *gin.Context, resultChannel chan []byte) {
		var resultMap = make(map[string]interface{})
		for _, v := range c.Params {
			resultMap[v.Key] = v.Value
		}
		for k, v := range c.Request.URL.Query() {
			resultMap[k] = v
		}
		json, err := json.Marshal(resultMap)
		if err != nil {
			resultChannel <- nil
			close(resultChannel)
			return
		}
		resultChannel <- json
		close(resultChannel)
	}(c, responseParamChannel)

	r.Header = <-headerChannel
	r.RequestBody = <-responseBodyChannel
	r.Params = <-responseParamChannel

	result := map[string]interface{}{
		Constants.Path:        r.Path,
		Constants.Method:      r.Method,
		Constants.Header:      r.GetHeaderMap(),
		Constants.RequestBody: r.GetRequestBodyMap(),
		Constants.Params:      r.GetParamsMap(),
	}

	return result
}

func (r *RequestInformation) GetHeaderJSON() []byte {
	return r.Header
}

func (r *RequestInformation) GetHeaderMap() map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(r.Header, &result)
	return result
}

func (r *RequestInformation) GetHeader() string {
	return string(r.Header)
}

func (r *RequestInformation) GetRequestBodyJSON() []byte {
	return r.RequestBody
}

func (r *RequestInformation) GetRequestBodyMap() map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(r.RequestBody, &result)
	return result
}

func (r *RequestInformation) GetRequestBody() string {
	return string(r.RequestBody)
}

func (r *RequestInformation) GetParamsJSON() []byte {
	return r.Params
}

func (r *RequestInformation) GetParamsMap() map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(r.Params, &result)
	return result
}

func (r *RequestInformation) GetParams() string {
	return string(r.Params)
}
