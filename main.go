package main

import (
	"github.com/dius/libpact/api"
	"github.com/dius/libpact/mock"
	"github.com/dius/libpact/pactfile"
)

func main() {
	pf := pactfile.New()
	go mock.Serve(":8080", pf)
	api.ConsumerServe(":5550", pf)
}
