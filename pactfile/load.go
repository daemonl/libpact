package pactfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Load reads a Root (pactfile) from the 'src', which could be file://, http://, https://
func Load(src string) (*Root, error) {
	parts := strings.SplitN(src, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("No protocol for %s", src)
	}
	protocol := parts[0]
	name := parts[1]
	switch protocol {
	case "file":
		return LoadFile(name)
	case "http", "https":
		return LoadHTTP(src)
	default:
		return nil, fmt.Errorf("Unknwon protocol '%s'", protocol)
	}
}

// LoadFile reads a pactfile from the filesystem
func LoadFile(filename string) (*Root, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadFromReader(file)
}

// LoadFromReader is a utility method for other loaders. Exported to allow
// custom loaders
func LoadFromReader(reader io.Reader) (*Root, error) {
	r := &Root{}
	err := json.NewDecoder(reader).Decode(r)
	return r, err
}

// LoadHTTP will read from a url
func LoadHTTP(url string) (*Root, error) {
	return nil, fmt.Errorf("HTTP Loader Not Implemented")
}
