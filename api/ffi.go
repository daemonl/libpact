package api

import "encoding/json"

type FFIRequest string

func (r *FFIRequest) ReadBodyInto(val interface{}) error {
	return json.Unmarshal([]byte(*r), val)
}

type FFIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"msg,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

func HandleFFI(req string, fn HandlerFunc) FFIResponse {

	r := FFIRequest(req)
	resp, err := fn(&r)
	if err != nil {
		return FFIResponse{
			Status:  500,
			Message: err.Error(),
		}
	}
	return FFIResponse{
		Status: resp.StatusCode(),
		Body:   resp.GetEncodable(),
	}
}
