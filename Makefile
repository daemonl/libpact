VERSION=$(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION)

libpact.so:
	go build -ldflags="$(FLAGS)" -o libpact.so -buildmode=c-shared libpact.go

pactserver:
	go build -ldflags="$(FLAGS)" -o pactserver main.go

.PHONY: libpact.so pactserver
