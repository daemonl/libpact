/*
Package api provides methods for creating an api server, served either over
HTTP or as a shared library, through FFI functions from other applications.

To achieve this, it provides a subset of HTTP, and a subset of what's possible with FFI.

fn(Request) -> Response is the basic interface for implementations
*/
package api

import (
	"fmt"
	"net/http"
)

// Mux wraps a HandlerFunc as either a HTTP Server of FFI Server
type Mux func(call string) HandlerFunc

// ServeHTTP implements a http server for a Mux
func (mux Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	mux(fmt.Sprintf("%s %s", req.Method, req.URL.Path)).
		ServeHTTP(rw, req)
}

// HandleCall implements a FFI Server for a Mux
func (mux Mux) HandleCall(call string, data string) FFIResponse {
	return mux(call).ServeFFI(data)
}
