package pactfile

import (
	"encoding/json"
	"net/http"
)

// New creates an empty Root
func New() *Root {
	return &Root{
		Interactions: []Interaction{},
	}
}

// Root is the root pactfile.
// TODO: All of it, and versioning.
type Root struct {
	Interactions []Interaction `json:"interactions"`
}

// RunInfo should be removed - it is used to count the number of times an
// interaction was run, but should instead be included in the mock server
type RunInfo struct {
	Count int `json:"count"`
}

// Interaction describes a request and expected response, given a provider
// state
type Interaction struct {
	RunInfo       `json:"-"`
	Description   string   `json:"description"`
	ProviderState string   `json:"provider_state"`
	Request       Request  `json:"request"`
	Response      Response `json:"response"`
}

// Request describes a request from the Client to the Provider
type Request struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Query   string      `json:"query"`
	Headers Headers     `json:"headers"`
	Body    interface{} `json:"body"`
}

// MatchesRequest should implement the Pact Matcher on the Method, Path,
// Querystring and Body of the request
func (pact *Request) MatchesRequest(h *http.Request) bool {
	if pact.Method != h.Method {
		return false
	}
	if pact.Path != h.URL.Path {
		return false
	}
	// TODO: How are Queries matched?
	return true
}

// Response describes how the mock server should respond to a request, as well
// as how the provider tests should verify responses.
type Response struct {
	Status  int          `json:"status"`
	Headers Headers      `json:"headers"`
	Body    ResponseNode `json:"body"`
}

// ServeHTTP serves a mock of the interaction to a consumer
func (pact *Response) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	pact.Headers.Set(rw.Header())
	rw.WriteHeader(pact.Status)

	json.NewEncoder(rw).Encode(pact.Body)
	// TODO: Something useful with the error from encode. The only candidates
	// here are Panic, Log or Write into the response (Which already has a
	// status code)
}

// Headers represents an expected or given header set. Only supports one value
// per key.
type Headers map[string]string

// Set sets each header into a http.Header. Linters are fussy. This was obvious.
func (pact Headers) Set(h http.Header) {
	for k, v := range pact {
		h.Add(k, v)
	}
}
