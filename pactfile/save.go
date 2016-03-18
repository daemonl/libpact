package pactfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func Save(dst string, root *Root) error {
	parts := strings.SplitN(dst, "://", 2)
	if len(parts) != 2 {
		return fmt.Errorf("No protocol for %s", dst)
	}
	protocol := parts[0]
	name := parts[1]
	switch protocol {
	case "file":
		return SaveFile(name, root)
	case "http", "https":
		return SaveHTTP(dst, root)
	default:
		return fmt.Errorf("Unknwon protocol '%s'", protocol)
	}
}

func SaveFile(filename string, root *Root) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return SaveToWriter(file, root)
}

func SaveToWriter(writer io.Writer, root *Root) error {
	return json.NewEncoder(writer).Encode(root)
}

func SaveHTTP(url string, root *Root) error {
	return fmt.Errorf("HTTP Loader Not Implemented")
}
