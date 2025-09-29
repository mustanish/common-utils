package assertionutil

import (
	"reflect"
	"testing"
)

func TestNewAssertionUtil(t *testing.T) {
	util := NewAssertionUtil()
	if util == nil {
		t.Error("NewAssertionUtil() returned nil")
	}
	if _, ok := util.(*AssertionUtil); !ok {
		t.Error("NewAssertionUtil() did not return *AssertionUtil")
	}
}

func TestGetString(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected string
		ok       bool
	}{
		{
			name:     "valid string",
			data:     map[string]any{"key": "value"},
			key:      "key",
			expected: "value",
			ok:       true,
		},
		{
			name:     "empty string",
			data:     map[string]any{"key": ""},
			key:      "key",
			expected: "",
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: "",
			ok:       false,
		},
		{
			name:     "wrong type",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: "",
			ok:       false,
		},
		{
			name:     "nil value",
			data:     map[string]any{"key": nil},
			key:      "key",
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetString(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetString() = (%q, %v), want (%q, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetStringRequired(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name      string
		data      map[string]any
		key       string
		expected  string
		shouldErr bool
	}{
		{
			name:      "valid string",
			data:      map[string]any{"key": "value"},
			key:       "key",
			expected:  "value",
			shouldErr: false,
		},
		{
			name:      "empty string",
			data:      map[string]any{"key": ""},
			key:       "key",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "missing key",
			data:      map[string]any{},
			key:       "key",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "wrong type",
			data:      map[string]any{"key": 123},
			key:       "key",
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.GetStringRequired(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("GetStringRequired() result = %q, want %q", result, tt.expected)
			}
			if (err != nil) != tt.shouldErr {
				t.Errorf("GetStringRequired() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected float64
		ok       bool
	}{
		{
			name:     "valid float64",
			data:     map[string]any{"key": 3.14},
			key:      "key",
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "zero value",
			data:     map[string]any{"key": 0.0},
			key:      "key",
			expected: 0.0,
			ok:       true,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "wrong type - int",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "3.14"},
			key:      "key",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetFloat64(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetFloat64() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected int
		ok       bool
	}{
		{
			name:     "valid int",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "float64 from JSON",
			data:     map[string]any{"key": 123.0},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "float64 with decimals",
			data:     map[string]any{"key": 123.5},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "123"},
			key:      "key",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetInt(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetInt() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetInt64(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected int64
		ok       bool
	}{
		{
			name:     "valid int64",
			data:     map[string]any{"key": int64(123456789)},
			key:      "key",
			expected: 123456789,
			ok:       true,
		},
		{
			name:     "valid int converted to int64",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "float64 from JSON",
			data:     map[string]any{"key": 123456789.0},
			key:      "key",
			expected: 123456789,
			ok:       true,
		},
		{
			name:     "float64 with decimals",
			data:     map[string]any{"key": 123.5},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetInt64(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetInt64() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected bool
		ok       bool
	}{
		{
			name:     "true value",
			data:     map[string]any{"key": true},
			key:      "key",
			expected: true,
			ok:       true,
		},
		{
			name:     "false value",
			data:     map[string]any{"key": false},
			key:      "key",
			expected: false,
			ok:       true,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: false,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "true"},
			key:      "key",
			expected: false,
			ok:       false,
		},
		{
			name:     "wrong type - int",
			data:     map[string]any{"key": 1},
			key:      "key",
			expected: false,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetBool(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetBool() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetMap(t *testing.T) {
	util := NewAssertionUtil()

	nestedMap := map[string]interface{}{"nested": "value"}

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected map[string]interface{}
		ok       bool
	}{
		{
			name:     "valid map",
			data:     map[string]any{"key": nestedMap},
			key:      "key",
			expected: nestedMap,
			ok:       true,
		},
		{
			name:     "empty map",
			data:     map[string]any{"key": map[string]interface{}{}},
			key:      "key",
			expected: map[string]interface{}{},
			ok:       true,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: nil,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "not a map"},
			key:      "key",
			expected: nil,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetMap(tt.data, tt.key)
			if !reflect.DeepEqual(result, tt.expected) || ok != tt.ok {
				t.Errorf("GetMap() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetSlice(t *testing.T) {
	util := NewAssertionUtil()

	slice := []interface{}{"item1", "item2"}

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected []interface{}
		ok       bool
	}{
		{
			name:     "valid slice",
			data:     map[string]any{"key": slice},
			key:      "key",
			expected: slice,
			ok:       true,
		},
		{
			name:     "empty slice",
			data:     map[string]any{"key": []interface{}{}},
			key:      "key",
			expected: nil,
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: nil,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "not a slice"},
			key:      "key",
			expected: nil,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetSlice(tt.data, tt.key)
			if !reflect.DeepEqual(result, tt.expected) || ok != tt.ok {
				t.Errorf("GetSlice() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetStringWithDefault(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name         string
		data         map[string]any
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "valid string",
			data:         map[string]any{"key": "value"},
			key:          "key",
			defaultValue: "default",
			expected:     "value",
		},
		{
			name:         "missing key",
			data:         map[string]any{},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "empty string",
			data:         map[string]any{"key": ""},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "wrong type",
			data:         map[string]any{"key": 123},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetStringWithDefault(tt.data, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetStringWithDefault() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetIntWithDefault(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name         string
		data         map[string]any
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid int",
			data:         map[string]any{"key": 123},
			key:          "key",
			defaultValue: 999,
			expected:     123,
		},
		{
			name:         "missing key",
			data:         map[string]any{},
			key:          "key",
			defaultValue: 999,
			expected:     999,
		},
		{
			name:         "wrong type",
			data:         map[string]any{"key": "123"},
			key:          "key",
			defaultValue: 999,
			expected:     999,
		},
		{
			name:         "float64 from JSON",
			data:         map[string]any{"key": 123.0},
			key:          "key",
			defaultValue: 999,
			expected:     123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetIntWithDefault(tt.data, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetIntWithDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetNumericAsFloat64(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected float64
		ok       bool
	}{
		{
			name:     "float64",
			data:     map[string]any{"key": 3.14},
			key:      "key",
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "float32",
			data:     map[string]any{"key": float32(3.14)},
			key:      "key",
			expected: float64(float32(3.14)),
			ok:       true,
		},
		{
			name:     "int",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: 123.0,
			ok:       true,
		},
		{
			name:     "int64",
			data:     map[string]any{"key": int64(123)},
			key:      "key",
			expected: 123.0,
			ok:       true,
		},
		{
			name:     "int32",
			data:     map[string]any{"key": int32(123)},
			key:      "key",
			expected: 123.0,
			ok:       true,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "123"},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetNumericAsFloat64(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetNumericAsFloat64() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetNumericAsInt(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected int
		ok       bool
	}{
		{
			name:     "int",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "int64 within int range",
			data:     map[string]any{"key": int64(123)},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "float64 without decimals",
			data:     map[string]any{"key": 123.0},
			key:      "key",
			expected: 123,
			ok:       true,
		},
		{
			name:     "float64 with decimals",
			data:     map[string]any{"key": 123.5},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "wrong type - string",
			data:     map[string]any{"key": "123"},
			key:      "key",
			expected: 0,
			ok:       false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetNumericAsInt(tt.data, tt.key)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetNumericAsInt() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetNestedString(t *testing.T) {
	util := NewAssertionUtil()

	data := map[string]any{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"value": "nested_value",
				"empty": "",
			},
		},
	}

	tests := []struct {
		name     string
		data     map[string]any
		path     []string
		expected string
		ok       bool
	}{
		{
			name:     "valid nested path",
			data:     data,
			path:     []string{"level1", "level2", "value"},
			expected: "nested_value",
			ok:       true,
		},
		{
			name:     "empty string at end",
			data:     data,
			path:     []string{"level1", "level2", "empty"},
			expected: "",
			ok:       false,
		},
		{
			name:     "invalid middle path",
			data:     data,
			path:     []string{"level1", "nonexistent", "value"},
			expected: "",
			ok:       false,
		},
		{
			name:     "invalid final key",
			data:     data,
			path:     []string{"level1", "level2", "nonexistent"},
			expected: "",
			ok:       false,
		},
		{
			name:     "single level path",
			data:     map[string]any{"key": "value"},
			path:     []string{"key"},
			expected: "value",
			ok:       true,
		},
		{
			name:     "empty path",
			data:     data,
			path:     []string{},
			expected: "",
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetNestedString(tt.data, tt.path...)
			if result != tt.expected || ok != tt.ok {
				t.Errorf("GetNestedString() = (%q, %v), want (%q, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestGetNestedMap(t *testing.T) {
	util := NewAssertionUtil()

	innerMap := map[string]interface{}{"value": "test"}
	data := map[string]any{
		"level1": map[string]interface{}{
			"level2": innerMap,
		},
	}

	tests := []struct {
		name     string
		data     map[string]any
		path     []string
		expected map[string]interface{}
		ok       bool
	}{
		{
			name:     "valid nested map",
			data:     data,
			path:     []string{"level1", "level2"},
			expected: innerMap,
			ok:       true,
		},
		{
			name:     "single level",
			data:     data,
			path:     []string{"level1"},
			expected: map[string]interface{}{"level2": innerMap},
			ok:       true,
		},
		{
			name:     "invalid path",
			data:     data,
			path:     []string{"level1", "nonexistent"},
			expected: nil,
			ok:       false,
		},
		{
			name:     "empty path returns original map",
			data:     data,
			path:     []string{},
			expected: data,
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := util.GetNestedMap(tt.data, tt.path...)
			if !reflect.DeepEqual(result, tt.expected) || ok != tt.ok {
				t.Errorf("GetNestedMap() = (%v, %v), want (%v, %v)", result, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestHasKey(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected bool
	}{
		{
			name:     "key exists with value",
			data:     map[string]any{"key": "value"},
			key:      "key",
			expected: true,
		},
		{
			name:     "key exists with nil value",
			data:     map[string]any{"key": nil},
			key:      "key",
			expected: true,
		},
		{
			name:     "key exists with empty string",
			data:     map[string]any{"key": ""},
			key:      "key",
			expected: true,
		},
		{
			name:     "key does not exist",
			data:     map[string]any{},
			key:      "key",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.HasKey(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("HasKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasNonEmptyString(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name     string
		data     map[string]any
		key      string
		expected bool
	}{
		{
			name:     "non-empty string",
			data:     map[string]any{"key": "value"},
			key:      "key",
			expected: true,
		},
		{
			name:     "empty string",
			data:     map[string]any{"key": ""},
			key:      "key",
			expected: false,
		},
		{
			name:     "nil value",
			data:     map[string]any{"key": nil},
			key:      "key",
			expected: false,
		},
		{
			name:     "wrong type",
			data:     map[string]any{"key": 123},
			key:      "key",
			expected: false,
		},
		{
			name:     "missing key",
			data:     map[string]any{},
			key:      "key",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.HasNonEmptyString(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("HasNonEmptyString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	util := NewAssertionUtil()

	tests := []struct {
		name      string
		data      map[string]any
		keys      []string
		shouldErr bool
		errMsg    string
	}{
		{
			name: "all required fields present",
			data: map[string]any{
				"name":  "John",
				"email": "john@example.com",
				"phone": "123-456-7890",
			},
			keys:      []string{"name", "email", "phone"},
			shouldErr: false,
		},
		{
			name: "one missing field",
			data: map[string]any{
				"name":  "John",
				"email": "john@example.com",
			},
			keys:      []string{"name", "email", "phone"},
			shouldErr: true,
			errMsg:    "required fields missing or empty: [phone]",
		},
		{
			name: "multiple missing fields",
			data: map[string]any{
				"name": "John",
			},
			keys:      []string{"name", "email", "phone"},
			shouldErr: true,
			errMsg:    "required fields missing or empty: [email phone]",
		},
		{
			name: "empty string field",
			data: map[string]any{
				"name":  "",
				"email": "john@example.com",
			},
			keys:      []string{"name", "email"},
			shouldErr: true,
			errMsg:    "required fields missing or empty: [name]",
		},
		{
			name:      "no required fields",
			data:      map[string]any{"key": "value"},
			keys:      []string{},
			shouldErr: false,
		},
		{
			name: "wrong type field",
			data: map[string]any{
				"name":  123,
				"email": "john@example.com",
			},
			keys:      []string{"name", "email"},
			shouldErr: true,
			errMsg:    "required fields missing or empty: [name]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidateRequired(tt.data, tt.keys...)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateRequired() error = %v, shouldErr %v", err, tt.shouldErr)
				return
			}
			if tt.shouldErr && err.Error() != tt.errMsg {
				t.Errorf("ValidateRequired() error message = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetString(b *testing.B) {
	util := NewAssertionUtil()
	data := map[string]any{"key": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.GetString(data, "key")
	}
}

func BenchmarkGetNestedString(b *testing.B) {
	util := NewAssertionUtil()
	data := map[string]any{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"value": "nested_value",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.GetNestedString(data, "level1", "level2", "value")
	}
}

func BenchmarkValidateRequired(b *testing.B) {
	util := NewAssertionUtil()
	data := map[string]any{
		"name":  "John",
		"email": "john@example.com",
		"phone": "123-456-7890",
	}
	keys := []string{"name", "email", "phone"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.ValidateRequired(data, keys...)
	}
}
