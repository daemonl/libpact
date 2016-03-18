package main

import "C"

import (
	"encoding/json"

	"github.com/dius/libpact/api"
	"github.com/dius/libpact/consumer"
	"github.com/dius/libpact/pactfile"
)

var Pact *pactfile.Root
var Consumer *consumer.Mux
var Mux api.Mux

func init() {
	Pact = pactfile.New()
	Consumer = &consumer.Mux{
		Pact: Pact,
	}
	Mux = api.Mux(Consumer.HandlerByName)
}

//export call
func call(cmethod *C.char, cjstring *C.char) *C.char {
	method := C.GoString(cmethod)
	jstring := C.GoString(cjstring)
	resp := callInternal(method, jstring)
	b, err := json.Marshal(resp)
	if err != nil {
		return C.CString(`{"status":500,error:"marshall error"}`)
	}
	return C.CString(string(b))
}

func callInternal(method string, jstring string) api.FFIResponse {
	if Pact == nil {
		return api.FFIResponse{
			Status: 500,
			Error:  "No Current Pactfile",
		}
	}

	return Mux.HandleCall(method, jstring)
}

func main() {}
