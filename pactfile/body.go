package pactfile

import "encoding/json"

// Node represents part of a JSON tree structure.
type Node interface {
	json.Marshaler
	json.Unmarshaler
	//Mock(?) (interface{}, error)
	//Diff(?) (?, error)
}

// ResponseBody is encoded directly to JSON,
// TODO: Implement Maps and Arrays, whatever is required for the matcher spec
type ResponseNode struct {
	Raw interface{}
}

// MarshalJSON is used to save the pactfile to disk or present it through an API call.
// Mocks are provided by the GetMock function, not here.
func (n ResponseNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Raw)
}

// UnmarshalJSON is used to load the pactfile from disk or an API call.
func (n *ResponseNode) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &n.Raw)
}

// Mock returns a tree of simple json encodable elements which represent a
// mocked response
//TODO func(n *ResponseNode) Mock(context?) (interface{}, error)

// Diff returns a tree representing the difference between the node and the provided interface.
// Diffs are given by example.
// TODO func(n *ResponseNode) Diff(context?) Diff
