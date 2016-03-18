package pactfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Save writes a Root (pactfile) to the 'dst', which could be file://, http://, https://
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

// SaveFile writes a pactfile to the filesystem
func SaveFile(filename string, root *Root) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return SaveToWriter(file, root)
}

// SaveToWriter is a utility method for other savers. Exported to allow custom savers
func SaveToWriter(writer io.Writer, root *Root) error {
	return json.NewEncoder(writer).Encode(root)
}

// SaveHTTP will write to a URL
func SaveHTTP(url string, root *Root) error {
	return fmt.Errorf("HTTP Loader Not Implemented")
}
