package pactfile

import (
	"encoding/json"
	"testing"
)

func TestDifference(t *testing.T) {
	s := String("a")
	diff, err := (&s).Diff(0, "b")
	if err != nil {
		t.Error(err)
		return
	}
	if diff.Match {
		t.Errorf("Expected a != b, but got a match")
		return
	}
	e, _ := json.Marshal(diff.Expect)
	se := string(e)
	if se != `"a"` {
		t.Errorf(`Expected diff.Expect == "a", got %s`, se)
	}
	if diff.Got != "b" {
		t.Errorf(`Expected diff.Got == "b", got %#v`, diff.Got)
	}
}

func TestFundamentals(t *testing.T) {
	/* Fundamental Types
	   https://golang.org/pkg/encoding/json/
	   bool, for JSON booleans
	   float64, for JSON numbers
	   string, for JSON strings
	   []interface{}, for JSON arrays
	   map[string]interface{}, for JSON objects
	   nil for JSON null
	*/

	match := func(name string, n Node, impl interface{}) {
		diff, err := n.Diff(0, impl)
		if err != nil {
			t.Errorf("Error diff for %s: %s", name, err.Error())
		} else if !diff.Match {
			t.Errorf("Should be no diff for %s, got %#v", name, diff)
		}
	}

	s := String("str")
	match("String", &s, "str")

	f := Float(3.4)
	match("Float", &f, 3.4)
	f = Float(4)
	match("Float", &f, 4)

	b := Bool(true)
	match("Bool", &b, true)

	m := Map(map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	})
	match("Map<string>", &m, map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	})

	a := Array([]interface{}{"a", "b"})
	match("Array<string>", &a, []interface{}{"a", "b"})
}
