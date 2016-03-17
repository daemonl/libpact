package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// the server is a thin wrapper around pactfiles, hard
	// linked, not through interfaces
	"github.com/dius/libpact/pactfile"
)

// This is a very basic test of the HTTP server functionality only.
// Matching of Request and Response is tested in the pactfile package.

func Test404(t *testing.T) {
	s := getTestServer()
	w, err := doReq(s, "GET", "/nopath", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%d - (%s)", w.Code, strings.TrimSpace(w.Body.String()))

	if w.Code != 500 {
		t.Errorf("/nopath should give a 500, got %d", w.Code)
	}
}

func TestGET(t *testing.T) {
	s := getTestServer()
	w, err := doReq(s, "GET", "/test1", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	sb := strings.TrimSpace(w.Body.String())
	t.Logf("%d - (%s)", w.Code, sb)

	if w.Code != 123 {
		t.Errorf("/test1 should give status 123, got %d", w.Code)
	}
	if w.Header().Get("Key1") != "Value1" {
		t.Errorf("/test1 header Key1 mismatch")
	}
	if sb != `"BODY1"` {
		t.Errorf("/test1 body, got (%s)", sb)
	}
}

func doReq(s *Server, method string, path string, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, "http://example.com"+path, body)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	return w, nil
}

func getTestServer() *Server {
	test1 := pactfile.Interaction{
		Request: pactfile.Request{
			Method: "GET",
			Path:   "/test1",
		},
		Response: pactfile.Response{
			Status: 123,
			Headers: pactfile.Headers{
				"Key1": "Value1",
			},
			Body: pactfile.ResponseNode{
				Raw: "BODY1",
			},
		},
	}
	s := &Server{
		Pact: &pactfile.Root{
			Interactions: []pactfile.Interaction{test1},
		},
	}
	return s
}
