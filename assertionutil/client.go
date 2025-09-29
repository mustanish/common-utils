// Package assertionutil provides safe type assertion and extraction utilities
// for working with map[string]any data structures commonly found in JSON parsing,
// configuration handling, and dynamic data processing.
//
// Features include:
//   - Safe type assertions with existence checks
//   - Required field validation with descriptive errors
//   - Support for common Go types (string, int, int64, float64, bool, map, slice)
//   - Default value fallbacks for missing or invalid fields
//   - Numeric type conversion utilities
//   - Nested path navigation for complex data structures
//   - Bulk validation for required fields
//   - Consistent error handling patterns
//   - Zero-allocation design for performance
//
// Example usage:
//
//	util := assertionutil.NewAssertionUtil()
//	data := map[string]any{"name": "John", "age": 30.0, "config": map[string]any{"debug": true}}
//
//	// Basic extraction
//	name, ok := util.GetString(data, "name")
//	age := util.GetIntWithDefault(data, "age", 0)
//
//	// Nested access
//	debug, _ := util.GetNestedString(data, "config", "debug")
//
//	// Validation
//	err := util.ValidateRequired(data, "name", "email")
package assertionutil

import "fmt"

// AssertionClient defines the interface for safe type assertion operations
type AssertionClient interface {
	// Basic type getters
	GetString(m map[string]any, key string) (string, bool)
	GetStringRequired(m map[string]any, key string) (string, error)
	GetFloat64(m map[string]any, key string) (float64, bool)
	GetMap(m map[string]any, key string) (map[string]interface{}, bool)
	GetSlice(m map[string]any, key string) ([]interface{}, bool)
	GetBool(m map[string]any, key string) (bool, bool)

	// Integer type getters
	GetInt(m map[string]any, key string) (int, bool)
	GetInt64(m map[string]any, key string) (int64, bool)

	// Getters with default values
	GetStringWithDefault(m map[string]any, key, defaultValue string) string
	GetFloat64WithDefault(m map[string]any, key string, defaultValue float64) float64
	GetIntWithDefault(m map[string]any, key string, defaultValue int) int
	GetBoolWithDefault(m map[string]any, key string, defaultValue bool) bool

	// Numeric conversion utilities
	GetNumericAsFloat64(m map[string]any, key string) (float64, bool)
	GetNumericAsInt(m map[string]any, key string) (int, bool)

	// Nested path access
	GetNestedString(m map[string]any, path ...string) (string, bool)
	GetNestedMap(m map[string]any, path ...string) (map[string]interface{}, bool)

	// Validation utilities
	HasKey(m map[string]any, key string) bool
	HasNonEmptyString(m map[string]any, key string) bool
	ValidateRequired(m map[string]any, keys ...string) error
}

// AssertionUtil provides safe type assertion utilities for map[string]any data structures
type AssertionUtil struct{}

// NewAssertionUtil creates a new assertion utility instance
func NewAssertionUtil() AssertionClient {
	return &AssertionUtil{}
}

// GetString safely extracts a non-empty string value from a map
func (a *AssertionUtil) GetString(m map[string]any, key string) (string, bool) {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok && str != "" {
			return str, true
		}
	}
	return "", false
}

// GetStringRequired safely extracts a required non-empty string value from a map
// Returns an error if the key doesn't exist, is not a string, or is empty
func (a *AssertionUtil) GetStringRequired(m map[string]any, key string) (string, error) {
	if str, ok := a.GetString(m, key); ok {
		return str, nil
	}
	return "", fmt.Errorf("required field '%s' not found or empty", key)
}

// GetFloat64 safely extracts a float64 value from a map
func (a *AssertionUtil) GetFloat64(m map[string]any, key string) (float64, bool) {
	if val, exists := m[key]; exists {
		if f, ok := val.(float64); ok {
			return f, true
		}
	}
	return 0, false
}

// GetMap safely extracts a nested map from a map
func (a *AssertionUtil) GetMap(m map[string]any, key string) (map[string]interface{}, bool) {
	if val, exists := m[key]; exists {
		if subMap, ok := val.(map[string]interface{}); ok {
			return subMap, true
		}
	}
	return nil, false
}

// GetSlice safely extracts a non-empty slice from a map
func (a *AssertionUtil) GetSlice(m map[string]any, key string) ([]interface{}, bool) {
	if val, exists := m[key]; exists {
		if slice, ok := val.([]interface{}); ok && len(slice) > 0 {
			return slice, true
		}
	}
	return nil, false
}

// GetBool safely extracts a boolean value from a map
func (a *AssertionUtil) GetBool(m map[string]any, key string) (bool, bool) {
	if val, exists := m[key]; exists {
		if b, ok := val.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// GetInt safely extracts an int value from a map
// Handles both int and float64 types from JSON unmarshaling
func (a *AssertionUtil) GetInt(m map[string]any, key string) (int, bool) {
	if val, exists := m[key]; exists {
		switch v := val.(type) {
		case int:
			return v, true
		case float64:
			// JSON numbers are unmarshaled as float64
			if v == float64(int(v)) {
				return int(v), true
			}
		}
	}
	return 0, false
}

// GetInt64 safely extracts an int64 value from a map
// Handles both int64 and float64 types from JSON unmarshaling
func (a *AssertionUtil) GetInt64(m map[string]any, key string) (int64, bool) {
	if val, exists := m[key]; exists {
		switch v := val.(type) {
		case int64:
			return v, true
		case int:
			return int64(v), true
		case float64:
			// JSON numbers are unmarshaled as float64
			if v == float64(int64(v)) {
				return int64(v), true
			}
		}
	}
	return 0, false
}

// GetStringWithDefault safely extracts a string value with a fallback default
func (a *AssertionUtil) GetStringWithDefault(m map[string]any, key, defaultValue string) string {
	if val, ok := a.GetString(m, key); ok {
		return val
	}
	return defaultValue
}

// GetFloat64WithDefault safely extracts a float64 value with a fallback default
func (a *AssertionUtil) GetFloat64WithDefault(m map[string]any, key string, defaultValue float64) float64 {
	if val, ok := a.GetFloat64(m, key); ok {
		return val
	}
	return defaultValue
}

// GetIntWithDefault safely extracts an int value with a fallback default
func (a *AssertionUtil) GetIntWithDefault(m map[string]any, key string, defaultValue int) int {
	if val, ok := a.GetInt(m, key); ok {
		return val
	}
	return defaultValue
}

// GetBoolWithDefault safely extracts a boolean value with a fallback default
func (a *AssertionUtil) GetBoolWithDefault(m map[string]any, key string, defaultValue bool) bool {
	if val, ok := a.GetBool(m, key); ok {
		return val
	}
	return defaultValue
}

// GetNumericAsFloat64 attempts to extract any numeric value as float64
// Handles int, int64, float32, float64 types
func (a *AssertionUtil) GetNumericAsFloat64(m map[string]any, key string) (float64, bool) {
	if val, exists := m[key]; exists {
		switch v := val.(type) {
		case float64:
			return v, true
		case float32:
			return float64(v), true
		case int:
			return float64(v), true
		case int64:
			return float64(v), true
		case int32:
			return float64(v), true
		}
	}
	return 0, false
}

// GetNumericAsInt attempts to extract any numeric value as int
// Only succeeds if the value can be represented as an integer without loss
func (a *AssertionUtil) GetNumericAsInt(m map[string]any, key string) (int, bool) {
	if val, exists := m[key]; exists {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			if v >= int64(int(^uint(0)>>1)*-1) && v <= int64(int(^uint(0)>>1)) {
				return int(v), true
			}
		case float64:
			if v == float64(int(v)) {
				return int(v), true
			}
		}
	}
	return 0, false
}

// GetNestedString safely extracts a string value from nested maps using a path
func (a *AssertionUtil) GetNestedString(m map[string]any, path ...string) (string, bool) {
	current := m
	for i, key := range path {
		if i == len(path)-1 {
			// Last key, get the string value
			return a.GetString(current, key)
		}
		// Navigate deeper
		next, ok := a.GetMap(current, key)
		if !ok {
			return "", false
		}
		current = next
	}
	return "", false
}

// GetNestedMap safely extracts a nested map using a path
func (a *AssertionUtil) GetNestedMap(m map[string]any, path ...string) (map[string]interface{}, bool) {
	current := m
	for _, key := range path {
		next, ok := a.GetMap(current, key)
		if !ok {
			return nil, false
		}
		current = next
	}
	return current, true
}

// HasKey checks if a key exists in the map (regardless of value type or nil)
func (a *AssertionUtil) HasKey(m map[string]any, key string) bool {
	_, exists := m[key]
	return exists
}

// HasNonEmptyString checks if a key exists and contains a non-empty string
func (a *AssertionUtil) HasNonEmptyString(m map[string]any, key string) bool {
	_, ok := a.GetString(m, key)
	return ok
}

// ValidateRequired checks that all specified keys exist and contain non-empty strings
func (a *AssertionUtil) ValidateRequired(m map[string]any, keys ...string) error {
	var missing []string
	for _, key := range keys {
		if !a.HasNonEmptyString(m, key) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required fields missing or empty: %v", missing)
	}
	return nil
}
