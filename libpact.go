package main

import "C"

import (
	"encoding/json"

	"github.com/dius/libpact/api"
	"github.com/dius/libpact/consumer"
	"github.com/dius/libpact/pactfile"
)

var mux api.Mux

func init() {
	pact := pactfile.New()
	consumer := &consumer.Mux{
		Pact: pact,
	}
	mux = api.Mux(consumer.HandlerByName)
}

//export call
func call(cmethod *C.char, cjstring *C.char) *C.char {
	method := C.GoString(cmethod)
	jstring := C.GoString(cjstring)
	resp := mux.HandleCall(method, jstring)
	b, err := json.Marshal(resp)
	if err != nil {
		return C.CString(`{"status":500,error:"marshall error"}`)
	}
	return C.CString(string(b))
}

func main() {}
