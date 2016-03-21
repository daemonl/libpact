package pactfile

func CoerceNode(v interface{}) (Node, bool) {

	// If it is already a node:
	if v, ok := v.(Node); ok {
		return v, true
	}

	// Basic json types
	switch v := v.(type) {
	case string:
		s := String(v)
		return &s, true

	case bool:
		b := Bool(v)
		return &b, true

	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		fv, _ := coerceFloat64(v)
		float := Float(fv)
		return &float, true

	case nil:
		return &Null{}, true

	case map[string]interface{}:
		m := Map(v)
		return &m, true

	case []interface{}:
		a := Array(v)
		return &a, true

	default:
		return nil, false
	}
}

// Does every language have one of these somewhere?
func coerceFloat64(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true

	default:
		return 0, false
	}
}
