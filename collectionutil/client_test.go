package collectionutil

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewCollectionUtil(t *testing.T) {
	util := NewCollectionUtil()
	if util == nil {
		t.Error("NewCollectionUtil() returned nil")
	}
}

// =================== Test Validation Methods ===================

func TestIsEmpty(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"whitespace string", "   ", true},
		{"non-empty string", "hello", false},
		{"empty slice", []any{}, true},
		{"non-empty slice", []any{1, 2, 3}, false},
		{"empty string slice", []string{}, true},
		{"non-empty string slice", []string{"a", "b"}, false},
		{"empty map", map[string]any{}, true},
		{"non-empty map", map[string]any{"key": "value"}, false},
		{"zero int", 0, false},
		{"non-zero int", 42, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.IsEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmpty(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	util := NewCollectionUtil()

	if util.IsNotEmpty("") {
		t.Error("IsNotEmpty(\"\") should return false")
	}

	if !util.IsNotEmpty("hello") {
		t.Error("IsNotEmpty(\"hello\") should return true")
	}
}

func TestKeyExistsAndNotEmpty(t *testing.T) {
	util := NewCollectionUtil()
	m := map[string]string{
		"existing": "value",
		"empty":    "",
		"spaces":   "   ",
	}

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"existing with value", "existing", true},
		{"existing but empty", "empty", false},
		{"existing but spaces", "spaces", false},
		{"non-existent", "missing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.KeyExistsAndNotEmpty(m, tt.key)
			if result != tt.expected {
				t.Errorf("KeyExistsAndNotEmpty(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

// =================== Test Type Conversion Methods ===================

func TestConvertToInteger(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name      string
		input     any
		expected  int
		expectErr bool
	}{
		{"int", 42, 42, false},
		{"int32", int32(42), 42, false},
		{"int64", int64(42), 42, false},
		{"float32", float32(42.0), 42, false},
		{"float64", 42.0, 42, false},
		{"string number", "42", 42, false},
		{"string with spaces", "  42  ", 42, false},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"invalid string", "abc", 0, true},
		{"slice", []string{"a"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToInteger(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToInteger(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ConvertToInteger(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToBool(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name      string
		input     any
		expected  bool
		expectErr bool
	}{
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"string true", "true", true, false},
		{"string false", "false", false, false},
		{"string 1", "1", true, false},
		{"string 0", "0", false, false},
		{"string yes", "yes", true, false},
		{"string no", "no", false, false},
		{"string empty", "", false, false},
		{"int 1", 1, true, false},
		{"int 0", 0, false, false},
		{"float 1.0", 1.0, true, false},
		{"float 0.0", 0.0, false, false},
		{"invalid string", "maybe", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToBool(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToBool(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ConvertToBool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool", true, "true"},
		{"nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ConvertToString(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToSlice(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name      string
		input     any
		separator string
		expected  []string
		expectErr bool
	}{
		{"comma separated", "a,b,c", ",", []string{"a", "b", "c"}, false},
		{"pipe separated", "a|b|c", "|", []string{"a", "b", "c"}, false},
		{"with spaces", "a, b, c", ",", []string{"a", "b", "c"}, false},
		{"empty string", "", ",", []string{}, false},
		{"single item", "hello", ",", []string{"hello"}, false},
		{"empty separator", "a,b,c", "", []string{"a", "b", "c"}, false},
		{"non-string input", 123, ",", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToSlice(tt.input, tt.separator)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToSlice(%v, %q) error = %v, expectErr %v", tt.input, tt.separator, err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToSlice(%v, %q) = %v, want %v", tt.input, tt.separator, result, tt.expected)
			}
		})
	}
}

func TestConvertToMap(t *testing.T) {
	util := NewCollectionUtil()

	// Test data structures
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	stringSlice := []string{"apple", "banana", "cherry"}
	intSlice := []int{1, 2, 3}

	tests := []struct {
		name         string
		input        any
		keyExtractor func(any) string
		expected     map[string]any
		expectErr    bool
	}{
		{
			name:  "convert person slice to map by name",
			input: people,
			keyExtractor: func(item any) string {
				person := item.(Person)
				return person.Name
			},
			expected: map[string]any{
				"Alice":   Person{Name: "Alice", Age: 30},
				"Bob":     Person{Name: "Bob", Age: 25},
				"Charlie": Person{Name: "Charlie", Age: 35},
			},
			expectErr: false,
		},
		{
			name:  "convert string slice to map with index as key",
			input: stringSlice,
			keyExtractor: func(item any) string {
				return item.(string)
			},
			expected: map[string]any{
				"apple":  "apple",
				"banana": "banana",
				"cherry": "cherry",
			},
			expectErr: false,
		},
		{
			name:  "convert int slice to map with string representation as key",
			input: intSlice,
			keyExtractor: func(item any) string {
				return fmt.Sprintf("key_%d", item.(int))
			},
			expected: map[string]any{
				"key_1": 1,
				"key_2": 2,
				"key_3": 3,
			},
			expectErr: false,
		},
		{
			name:         "nil slice",
			input:        nil,
			keyExtractor: func(item any) string { return "key" },
			expected:     map[string]any{},
			expectErr:    false,
		},
		{
			name:         "empty slice",
			input:        []string{},
			keyExtractor: func(item any) string { return item.(string) },
			expected:     map[string]any{},
			expectErr:    false,
		},
		{
			name:         "non-slice input",
			input:        "not a slice",
			keyExtractor: func(item any) string { return "key" },
			expected:     nil,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToMap(tt.input, tt.keyExtractor)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToMap() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertToInt64(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name      string
		input     any
		expected  int64
		expectErr bool
	}{
		{"int", 42, 42, false},
		{"int64", int64(42), 42, false},
		{"string", "42", 42, false},
		{"invalid string", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToInt64(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToInt64(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ConvertToInt64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToFloat64(t *testing.T) {
	util := NewCollectionUtil()

	tests := []struct {
		name      string
		input     any
		expected  float64
		expectErr bool
	}{
		{"int", 42, 42.0, false},
		{"float64", 42.5, 42.5, false},
		{"string", "42.5", 42.5, false},
		{"invalid string", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ConvertToFloat64(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertToFloat64(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ConvertToFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSliceContainsAny(t *testing.T) {
	util := NewCollectionUtil()
	slice := []any{1, "hello", 3.14, true}

	if !util.SliceContainsAny(slice, "hello") {
		t.Error("SliceContainsAny should return true for existing string")
	}

	if !util.SliceContainsAny(slice, 1) {
		t.Error("SliceContainsAny should return true for existing int")
	}

	if util.SliceContainsAny(slice, "world") {
		t.Error("SliceContainsAny should return false for non-existing item")
	}
}

func TestMapFilter(t *testing.T) {
	util := NewCollectionUtil()
	input := map[string]any{
		"a": 1,
		"b": "hello",
		"c": 3,
	}

	// Filter for numeric values
	result := util.MapFilter(input, func(k string, v any) bool {
		_, ok := v.(int)
		return ok
	})

	expected := map[string]any{
		"a": 1,
		"c": 3,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapFilter() = %v, want %v", result, expected)
	}
}

func TestMapValues(t *testing.T) {
	util := NewCollectionUtil()
	input := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	values := util.MapValues(input)
	if len(values) != 3 {
		t.Errorf("MapValues() should return 3 values, got %d", len(values))
	}

	// Check that all values are present (order doesn't matter)
	valueSet := make(map[any]bool)
	for _, v := range values {
		valueSet[v] = true
	}

	for _, expectedValue := range []any{1, 2, 3} {
		if !valueSet[expectedValue] {
			t.Errorf("MapValues() missing expected value %v", expectedValue)
		}
	}
}

func TestSliceChunkEdgeCases(t *testing.T) {
	util := NewCollectionUtil()

	// Test with size <= 0
	result := util.SliceChunk([]string{"a", "b", "c"}, 0)
	if len(result) != 0 {
		t.Error("SliceChunk with size 0 should return empty slice")
	}

	// Test with size larger than slice
	result = util.SliceChunk([]string{"a", "b"}, 5)
	expected := [][]string{{"a", "b"}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceChunk() = %v, want %v", result, expected)
	}
}

// =================== Test Go-Funk Powered Methods ===================

func TestSliceIntersection(t *testing.T) {
	util := NewCollectionUtil()
	slice1 := []string{"a", "b", "c", "d"}
	slice2 := []string{"c", "d", "e", "f"}
	expected := []string{"c", "d"}

	result := util.SliceIntersection(slice1, slice2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceIntersection(%v, %v) = %v, want %v", slice1, slice2, result, expected)
	}
}

func TestSliceDifference(t *testing.T) {
	util := NewCollectionUtil()
	slice1 := []string{"a", "b", "c", "d"}
	slice2 := []string{"c", "d", "e", "f"}
	expected := []string{"a", "b"}

	result := util.SliceDifference(slice1, slice2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceDifference(%v, %v) = %v, want %v", slice1, slice2, result, expected)
	}
}

func TestSliceUnion(t *testing.T) {
	util := NewCollectionUtil()
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"c", "d", "e"}

	result := util.SliceUnion(slice1, slice2)

	// Check that all unique elements are present
	expected := []string{"a", "b", "c", "d", "e"}
	if len(result) != len(expected) {
		t.Errorf("SliceUnion(%v, %v) length = %d, want %d", slice1, slice2, len(result), len(expected))
	}

	// Check each expected element is present
	for _, item := range expected {
		if !util.SliceContains(result, item) {
			t.Errorf("SliceUnion result missing expected item: %s", item)
		}
	}
}

func TestSliceContains(t *testing.T) {
	util := NewCollectionUtil()
	slice := []string{"apple", "banana", "cherry"}

	if !util.SliceContains(slice, "banana") {
		t.Error("SliceContains should return true for existing item")
	}

	if util.SliceContains(slice, "grape") {
		t.Error("SliceContains should return false for non-existing item")
	}
}

func TestSliceUnique(t *testing.T) {
	util := NewCollectionUtil()
	input := []string{"a", "b", "a", "c", "b", "d"}
	expected := []string{"a", "b", "c", "d"}

	result := util.SliceUnique(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceUnique(%v) = %v, want %v", input, result, expected)
	}
}

func TestSliceFilter(t *testing.T) {
	util := NewCollectionUtil()
	input := []string{"apple", "banana", "apricot", "cherry"}
	predicate := func(s string) bool { return len(s) > 5 }
	expected := []string{"banana", "apricot", "cherry"}

	result := util.SliceFilter(input, predicate)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceFilter() = %v, want %v", result, expected)
	}
}

func TestSliceMap(t *testing.T) {
	util := NewCollectionUtil()
	input := []string{"hello", "world"}
	mapper := func(s string) string { return s + "!" }
	expected := []string{"hello!", "world!"}

	result := util.SliceMap(input, mapper)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceMap() = %v, want %v", result, expected)
	}
}

func TestSliceReverse(t *testing.T) {
	util := NewCollectionUtil()
	input := []string{"a", "b", "c"}
	expected := []string{"c", "b", "a"}

	result := util.SliceReverse(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceReverse(%v) = %v, want %v", input, result, expected)
	}
}

func TestSliceChunk(t *testing.T) {
	util := NewCollectionUtil()
	input := []string{"a", "b", "c", "d", "e"}
	expected := [][]string{{"a", "b"}, {"c", "d"}, {"e"}}

	result := util.SliceChunk(input, 2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SliceChunk(%v, 2) = %v, want %v", input, result, expected)
	}
}

// =================== Test Map Operations ===================

func TestMapKeys(t *testing.T) {
	util := NewCollectionUtil()
	input := map[string]any{"b": 2, "a": 1, "c": 3}
	expected := []string{"a", "b", "c"} // sorted

	result := util.MapKeys(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapKeys(%v) = %v, want %v", input, result, expected)
	}
}

func TestMapMerge(t *testing.T) {
	util := NewCollectionUtil()
	map1 := map[string]any{"a": 1, "b": 2}
	map2 := map[string]any{"b": 3, "c": 4}
	expected := map[string]any{"a": 1, "b": 3, "c": 4}

	result := util.MapMerge(map1, map2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapMerge() = %v, want %v", result, expected)
	}
}

func TestMapPick(t *testing.T) {
	util := NewCollectionUtil()
	input := map[string]any{"a": 1, "b": 2, "c": 3}
	expected := map[string]any{"a": 1, "c": 3}

	result := util.MapPick(input, "a", "c", "d") // "d" doesn't exist
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapPick() = %v, want %v", result, expected)
	}
}

func TestMapOmit(t *testing.T) {
	util := NewCollectionUtil()
	input := map[string]any{"a": 1, "b": 2, "c": 3}
	expected := map[string]any{"a": 1, "c": 3}

	result := util.MapOmit(input, "b")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MapOmit() = %v, want %v", result, expected)
	}
}

// =================== Test Utility Methods ===================

func TestFindInSlice(t *testing.T) {
	util := NewCollectionUtil()
	slice := []any{1, "hello", 3.14, true}

	// Find string
	result, found := util.FindInSlice(slice, func(item any) bool {
		s, ok := item.(string)
		return ok && s == "hello"
	})

	if !found {
		t.Error("FindInSlice should find the string item")
	}

	if result != "hello" {
		t.Errorf("FindInSlice() = %v, want %v", result, "hello")
	}

	// Find non-existent
	_, found = util.FindInSlice(slice, func(item any) bool {
		s, ok := item.(string)
		return ok && s == "world"
	})

	if found {
		t.Error("FindInSlice should not find non-existent item")
	}
}

// =================== Benchmarks ===================

func BenchmarkIsEmpty(b *testing.B) {
	util := NewCollectionUtil()
	testData := []any{
		"",
		"hello",
		[]any{},
		[]any{1, 2, 3},
		map[string]any{},
		map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, data := range testData {
			util.IsEmpty(data)
		}
	}
}

func BenchmarkSliceContains(b *testing.B) {
	util := NewCollectionUtil()
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.SliceContains(slice, "cherry")
	}
}

func BenchmarkConvertToInteger(b *testing.B) {
	util := NewCollectionUtil()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		util.ConvertToInteger("12345")
	}
}
