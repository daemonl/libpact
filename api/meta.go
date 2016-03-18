package api

// HandlerFunc is api's equivalent of net.http.Handler
type HandlerFunc func(Request) (Response, error)

// Request contains the minimum information for an API call, can be implemented
// as an FFI Function call, or represent a HTTP Handler
type Request interface {
	ReadBodyInto(interface{}) error
}

// Response contains the minimum information for an API response, can be
// an FFI call's return, or HTTP Response
type Response interface {
	GetEncodable() interface{}
	StatusCode() int
}

// ObjectResponse just writes an object
type ObjectResponse struct {
	// A reader to write as the response body
	Object interface{}

	// Status uses http status codes, even for FFI, and should be mapped to
	// errors if required
	Status int
}

// StatusCode for implementation of Response
func (resp *ObjectResponse) StatusCode() int {
	return resp.Status
}

// GetEncodable returns the gooey insides
func (resp *ObjectResponse) GetEncodable() interface{} {
	return resp.Object
}

// BuildStringResponse is a basic response which has only a string as the body. e.g. "OK"
func BuildStringResponse(status int, msg string) Response {
	return &ObjectResponse{
		Object: map[string]interface{}{"msg": msg},
		Status: status,
	}
}

// BuildObjectResponse is a basic response with an interface body.
func BuildObjectResponse(status int, msg interface{}) Response {
	return &ObjectResponse{
		Object: msg,
		Status: status,
	}
}

// NotFound implements HandlerFunc where no handlers are
// found
func NotFound(req Request) (Response, error) {
	return BuildStringResponse(404, "No Such Call"), nil
}
