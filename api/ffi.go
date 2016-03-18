package api

import "encoding/json"

// FFIRequest implements api.Request for ffi requests
type FFIRequest string

// ReadBodyInto decodes the JSON body
func (r *FFIRequest) ReadBodyInto(val interface{}) error {
	return json.Unmarshal([]byte(*r), val)
}

// FFIResponse gets written back to the caller
type FFIResponse struct {
	Status int         `json:"status"`
	Error  string      `json:"error,omitempty"`
	Body   interface{} `json:"body,omitempty"`
}

// ServeFFI implements an FFI Handler for a generic handler func
func (fn HandlerFunc) ServeFFI(req string) FFIResponse {
	r := FFIRequest(req)
	resp, err := fn(&r)
	if err != nil {
		return FFIResponse{
			Status: 500,
			Error:  err.Error(),
		}
	}
	return FFIResponse{
		Status: resp.StatusCode(),
		Body:   resp.GetEncodable(),
	}
}
