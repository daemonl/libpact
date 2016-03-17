package mock

import (
	"net/http"

	"github.com/dius/libpact/pactfile"
)

type Server struct {
	Pact *pactfile.Root
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, interaction := range s.Pact.Interactions {
		if interaction.Request.MatchesRequest(req) {
			interaction.Response.ServeHTTP(rw, req)
			return
		}
	}
	http.NotFound(rw, req)
}

func Serve(bind string, pact *pactfile.Root) error {
	return http.ListenAndServe(bind, &Server{
		Pact: pact,
	})
}
