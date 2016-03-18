package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPRequest http.Request

func (r *HTTPRequest) ReadBodyInto(val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	defer r.Body.Close()
	return err
}

func (fn HandlerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
