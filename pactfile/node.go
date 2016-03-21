package pactfile

import (
	"encoding/json"
	"fmt"
)

type DiffFlag int

const (
	DiffFlag_AllowExtraKeysInMap = 1 << iota
	DiffFlag_AllowExtraElementsInArray
)

// Node represents part of a JSON tree structure.
type Node interface {
	//json.Marshaler
	//json.Unmarshaler
	// Optional, so we don't need to implement for builtins

	Mock() (interface{}, error)
	Diff(DiffFlag, interface{}) (*Diff, error)
}

type Diff struct {
	Match  bool
	Expect interface{}
	Got    interface{}
}

type String string

func (s String) Mock() (interface{}, error) {
	return s, nil
}

func (s String) Diff(_ DiffFlag, impl interface{}) (*Diff, error) {
	if impl != string(s) {
		return &Diff{
			Match:  false,
			Got:    impl,
			Expect: &s,
		}, nil
	}
	return &Diff{
		Match:  true,
		Got:    impl,
		Expect: impl,
	}, nil
}

type Float float64

func (f Float) Mock() (interface{}, error) {
	return f, nil
}

func (f Float) Diff(_ DiffFlag, impl interface{}) (*Diff, error) {
	implFloat, ok := coerceFloat64(impl)
	if !ok {
		return &Diff{
			Match:  false,
			Got:    impl,
			Expect: &f,
		}, nil
	}
	if implFloat != float64(f) {
		return &Diff{
			Match:  false,
			Got:    impl,
			Expect: &f,
		}, nil
	}
	return &Diff{
		Match:  true,
		Expect: impl,
		Got:    impl,
	}, nil
}

type Bool bool

func (b Bool) Mock() (interface{}, error) {
	return b, nil
}

func (b Bool) Diff(_ DiffFlag, impl interface{}) (*Diff, error) {
	if impl != bool(b) {
		return &Diff{
			Match:  false,
			Got:    impl,
			Expect: &b,
		}, nil
	}
	return &Diff{
		Match:  true,
		Expect: impl,
		Got:    impl,
	}, nil
}

type Null struct{}

func (n Null) Mock() (interface{}, error) {
	return nil, nil
}

func (n Null) Diff(_ DiffFlag, impl interface{}) (*Diff, error) {
	if impl != nil {
		return &Diff{
			Match:  false,
			Expect: nil,
			Got:    impl,
		}, nil
	}
	return &Diff{
		Match:  true,
		Expect: nil,
		Got:    nil,
	}, nil
}

type Map map[string]interface{}

func (m Map) Mock() (interface{}, error) {
	return m, nil
}

func (m Map) Diff(flags DiffFlag, implRaw interface{}) (*Diff, error) {
	impl, ok := implRaw.(map[string]interface{})
	if !ok {
		return &Diff{
			Match:  false,
			Got:    implRaw,
			Expect: &m,
		}, nil
	}

	match := true
	expect := map[string]interface{}{}
	got := map[string]interface{}{}

	for k, v := range m {
		valNode, ok := CoerceNode(v)
		if !ok {
			return nil, fmt.Errorf("Could not coerce %T into a node", v)
		}
		diff, err := valNode.Diff(flags, impl[k])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", k, err.Error())
		}

		if !diff.Match {
			match = false
		}
		expect[k] = diff.Expect
		got[k] = diff.Got
	}

	if !(flags&DiffFlag_AllowExtraKeysInMap == DiffFlag_AllowExtraKeysInMap) {
		for k, v := range impl {
			if _, ok := impl[k]; !ok {
				match = false
				expect[k] = nil
				got[k] = v
			}
		}
	}

	return &Diff{
		Match:  match,
		Expect: expect,
		Got:    got,
	}, nil
}

type Array []interface{}

func (a Array) Mock() (interface{}, error) {
	return a, nil
}

func (a Array) Diff(flags DiffFlag, implRaw interface{}) (*Diff, error) {

	impl, ok := implRaw.([]interface{})
	if !ok {
		return &Diff{
			Match:  false,
			Got:    implRaw,
			Expect: &a,
		}, nil
	}

	match := true
	expect := make([]interface{}, len(a), len(a))

	for i, v := range a {
		if len(impl)-1 < i {
			match = false
			expect[i] = v
			continue
		}
		valNode, ok := CoerceNode(v)
		if !ok {
			return nil, fmt.Errorf("Could not coerce %t into a node")
		}
		diff, err := valNode.Diff(flags, impl[i])
		if err != nil {
			return nil, fmt.Errorf("%d: %s", i, err.Error())
		}

		if !diff.Match {
			match = false
		}
		expect[i] = diff.Expect
	}

	return &Diff{
		Match:  match,
		Expect: expect,
		Got:    impl,
	}, nil
}

// ResponseNode is encoded directly to JSON,
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
