package main

import "C"

import (
	"encoding/json"

	"github.com/dius/libpact/api"
	"github.com/dius/libpact/mock"
	"github.com/dius/libpact/pactfile"
)

var Pact *pactfile.Root
var Consumer *api.Consumer

func init() {
	Pact = &pactfile.Root{
		Interactions: []pactfile.Interaction{},
	}
	Consumer = &api.Consumer{
		Pact: Pact,
	}
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

	switch method {
	case "mock":
		bind := ""
		json.Unmarshal([]byte(jstring), &bind)
		go mock.Serve(bind, Pact)
		return api.FFIResponse{
			Status: 200,
			Body:   "OK",
		}
	}

	return Consumer.HandleCall(method, jstring)
}

func main() {}
