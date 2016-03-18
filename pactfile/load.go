package pactfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

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

func LoadFile(filename string) (*Root, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadFromReader(file)
}

func LoadFromReader(reader io.Reader) (*Root, error) {
	r := &Root{}
	err := json.NewDecoder(reader).Decode(r)
	return r, err
}

func LoadHTTP(url string) (*Root, error) {
	return nil, fmt.Errorf("HTTP Loader Not Implemented")
}
