package mock

import (
	"encoding/json"
	"net/http"

	"github.com/dius/libpact/pactfile"
)

type Server struct {
	Pact *pactfile.Root
}

func descriptionsOfInteractions(interactions []pactfile.Interaction) []string {
	desc := make([]string, len(interactions), len(interactions))
	for idx, interaction := range interactions {
		desc[idx] = interaction.Description
	}
	return desc
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	matchedInteractions := []pactfile.Interaction{}
	for _, interaction := range s.Pact.Interactions {
		if interaction.Request.MatchesRequest(req) {
			matchedInteractions = append(matchedInteractions, interaction)
		}
	}

	if len(matchedInteractions) > 1 {
		rw.WriteHeader(500)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":      "More than one interaction matched",
			"interactions": descriptionsOfInteractions(matchedInteractions),
		})

		return
	}
	if len(matchedInteractions) < 1 {
		rw.WriteHeader(500)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":      "No interaction matched",
			"interactions": descriptionsOfInteractions(s.Pact.Interactions),
		})
		return
	}
	interaction := matchedInteractions[0]
	interaction.Response.ServeHTTP(rw, req)
	interaction.RunInfo.Count += 1

	/*
		Handling requests

		When a request comes in to the mock service, the request is compared with
		each interaction that has been registered with the mock service to find the
		right response to return.

		The rules for determining whether a request "matches" or not are defined by
		the pact-specification. It must "match" the path, query, headers and body
		according to the pact specification matching rules.

		If no interactions match the given request, then a 500 error should be
		returned by the mock service with an error indicating that no matches have
		been found. Include a list of the registered interactions to assist with
		debugging.

		If more than one interaction matches the given request, then a 500 error
		should be returned by the mock service, with a helpful error message. Each
		matching interaction should be logged and returned in the response body to
		assist with debugging.

		If exactly one interaction matches the given request, than the
		corresponding response should be returned, and that interaction should be
		marked as received.

		https://github.com/pact-foundation/pact-specification/blob/master/implem
		entation-guidelines/README.md#handling-requests
	*/
}

func Serve(bind string, pact *pactfile.Root) error {
	return http.ListenAndServe(bind, &Server{
		Pact: pact,
	})
}
