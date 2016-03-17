package pactfile

import "encoding/json"

type Node interface {
	json.Marshaler
	json.Unmarshaler
	IsMatchedBy(interface{}) (bool, error)
}

// ResponseBody is encoded directly to JSON,
// TODO: Implement Maps and Arrays, whatever is required for the matcher spec
type ResponseNode struct {
	Raw interface{}
}

// MarshalJSON is used to provide an example implementation when mocking
// any funky attributes should implement json.Encodable, giving an example.
func (n ResponseNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Raw)
}

// UnmarshalJSON is used to read from the pactfile, so any funky attributes
// shoule decode into an implementation which can be matched, and marshalled to
// an example
func (n *ResponseNode) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &n.Raw)
}

// TODO: A third use for Matcher, which implements the same style as the
// standard JSON/XML libraries
func (n *ResponseNode) IsMatchedBy(i interface{}) (bool, error) {
	return true, nil
}
