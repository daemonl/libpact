package pactfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"
)

func TestV2RequestBody(t *testing.T) {
	files, err := getGithubFiles("version-2", "request/body")
	if err != nil {
		t.Error(err)
		return
	}

	for _, file := range files {
		testCase, err := file.DownloadTest()
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Compare file %s\n", file.Path)
		got := http.Request(testCase.Actual)
		matches := testCase.Expected.MatchesRequest(&got)
		if matches == testCase.Match {
			t.Logf("Test case %s OK", file.Path)
		} else {
			t.Errorf("Test case %s, Should: %t, Got: %t", file.Path, testCase.Match, matches)
		}

	}

}

type testRequest http.Request

func (tr *testRequest) UnmarshalJSON(data []byte) error {
	r := &struct {
		Method  string          `json:"method"`
		Path    string          `json:"path"`
		Query   string          `json:"query"`
		Headers Headers         `json:"headers"`
		Body    json.RawMessage `json:"body"`
	}{}
	err := json.Unmarshal(data, r)
	if err != nil {
		return err
	}

	tr.Method = r.Method
	tr.URL = &url.URL{
		Path: r.Path,
	}
	tr.Header = http.Header{}
	r.Headers.Set(tr.Header)

	tr.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(r.Body)))
	return nil
}

type TestCase struct {
	Match    bool        `json:"match"`
	Comment  string      `json:"comment"`
	Expected Request     `json:"expected"`
	Actual   testRequest `json:"actual"`
}

type githubFile struct {
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}

func (f *githubFile) DownloadTest() (*TestCase, error) {
	resp, err := http.Get(f.DownloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tc := &TestCase{}
	return tc, json.NewDecoder(resp.Body).Decode(tc)
}

func getGithubFiles(version string, testPath string) ([]githubFile, error) {

	resp, err := http.Get("https://" + path.Join("api.github.com/repos/pact-foundation/pact-specification/contents/testcases/", testPath) + "?ref=" + version)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	files := []githubFile{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	return files, err
}
