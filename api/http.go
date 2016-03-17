package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerFunc func(Request) (Response, error)

type Handler struct {
	GET  HandlerFunc
	POST HandlerFunc
}

type HTTPRequest http.Request

func (r *HTTPRequest) ReadBodyInto(val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	defer r.Body.Close()
	return err
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var fn HandlerFunc

	switch req.Method {
	case "GET":
		fn = h.GET
	case "POST":
		fn = h.POST
	}

	if fn == nil {
		http.NotFound(rw, req)
		return
	}

	r := HTTPRequest(*req)
	resp, err := fn(&r)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.WriteHeader(resp.StatusCode())
	json.NewEncoder(rw).Encode(resp.GetEncodable())
}
