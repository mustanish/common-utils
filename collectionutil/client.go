package collectionutil

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

// CollectionClient defines the interface for collection utility operations
type CollectionClient interface {
	// Validation methods
	IsEmpty(val any) bool
	IsNotEmpty(val any) bool
	KeyExistsAndNotEmpty(value map[string]string, key string) bool

	// Type conversion methods
	ConvertToInteger(value any) (int, error)
	ConvertToInt64(value any) (int64, error)
	ConvertToFloat64(value any) (float64, error)
	ConvertToBool(value any) (bool, error)
	ConvertToString(value any) string
	ConvertToSlice(value any, separator string) ([]string, error)
	ConvertToMap(slice any, keyExtractor func(any) string) (map[string]any, error)

	// Slice operations
	SliceContains(slice []string, item string) bool
	SliceContainsAny(slice []any, item any) bool
	SliceUnique(slice []string) []string
	SliceFilter(slice []string, predicate func(string) bool) []string
	SliceMap(slice []string, mapper func(string) string) []string
	SliceReverse(slice []string) []string
	SliceChunk(slice []string, size int) [][]string

	// Map operations
	MapKeys(m map[string]any) []string
	MapValues(m map[string]any) []any
	MapMerge(maps ...map[string]any) map[string]any
	MapFilter(m map[string]any, predicate func(string, any) bool) map[string]any
	MapPick(m map[string]any, keys ...string) map[string]any
	MapOmit(m map[string]any, keys ...string) map[string]any

	// Utility methods
	FindInSlice(slice []any, predicate func(any) bool) (any, bool)

	// Additional slice operations
	SliceIntersection(slice1, slice2 []string) []string
	SliceDifference(slice1, slice2 []string) []string
	SliceUnion(slice1, slice2 []string) []string
}

type CollectionUtil struct{}

// NewCollectionUtil creates a new instance of CollectionUtil
func NewCollectionUtil() CollectionClient {
	return &CollectionUtil{}
}

// IsEmpty checks if a value is considered empty
func (c *CollectionUtil) IsEmpty(val any) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []any:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	default:
		// Use reflection for other slice/map types
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
			return rv.Len() == 0
		case reflect.Ptr:
			return rv.IsNil()
		}
		return false
	}
}

// IsNotEmpty checks if a value is not empty
func (c *CollectionUtil) IsNotEmpty(val any) bool {
	return !c.IsEmpty(val)
}

// KeyExistsAndNotEmpty checks if a key exists in the map and its value is not empty
func (c *CollectionUtil) KeyExistsAndNotEmpty(value map[string]string, key string) bool {
	if val, ok := value[key]; ok && strings.TrimSpace(val) != "" {
		return true
	}
	return false
}

// ConvertToInteger converts a value to an integer with error handling
func (c *CollectionUtil) ConvertToInteger(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		val, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		return int(val), err
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to integer", value)
	}
}

// ConvertToInt64 converts a value to int64
func (c *CollectionUtil) ConvertToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ConvertToFloat64 converts a value to float64
func (c *CollectionUtil) ConvertToFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ConvertToBool converts a value to boolean
func (c *CollectionUtil) ConvertToBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		trimmed := strings.ToLower(strings.TrimSpace(v))
		switch trimmed {
		case "true", "1", "yes", "on", "t", "y":
			return true, nil
		case "false", "0", "no", "off", "f", "n", "":
			return false, nil
		default:
			return strconv.ParseBool(v)
		}
	case int:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ConvertToString converts any value to string
func (c *CollectionUtil) ConvertToString(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// ConvertToSlice converts a string to slice using separator
func (c *CollectionUtil) ConvertToSlice(value any, separator string) ([]string, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string value for slice conversion, got %T", value)
	}

	if separator == "" {
		separator = ","
	}

	if strings.TrimSpace(strValue) == "" {
		return []string{}, nil
	}

	parts := strings.Split(strValue, separator)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts, nil
}

// ConvertToMap converts a slice to a map using a key extractor function
func (c *CollectionUtil) ConvertToMap(slice any, keyExtractor func(any) string) (map[string]any, error) {
	if slice == nil {
		return map[string]any{}, nil
	}

	// Use reflection to handle different slice types
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("expected slice or array, got %T", slice)
	}

	result := make(map[string]any)
	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		key := keyExtractor(item)
		result[key] = item
	}

	return result, nil
}

// SliceContains checks if a string slice contains a specific item
func (c *CollectionUtil) SliceContains(slice []string, item string) bool {
	return funk.Contains(slice, item)
}

// SliceContainsAny checks if a slice contains a specific item (any type)
func (c *CollectionUtil) SliceContainsAny(slice []any, item any) bool {
	return funk.Contains(slice, item)
}

// SliceUnique returns a slice with unique elements
func (c *CollectionUtil) SliceUnique(slice []string) []string {
	return funk.UniqString(slice)
}

// SliceFilter filters a slice based on a predicate function
func (c *CollectionUtil) SliceFilter(slice []string, predicate func(string) bool) []string {
	return funk.FilterString(slice, predicate)
}

// SliceMap transforms each element in a slice using a mapper function
func (c *CollectionUtil) SliceMap(slice []string, mapper func(string) string) []string {
	return funk.Map(slice, mapper).([]string)
}

// SliceReverse returns a reversed copy of the slice
func (c *CollectionUtil) SliceReverse(slice []string) []string {
	reversed := funk.ReverseStrings(slice)
	return reversed
}

// SliceChunk splits a slice into chunks of specified size
func (c *CollectionUtil) SliceChunk(slice []string, size int) [][]string {
	if size <= 0 {
		return [][]string{}
	}

	var chunks [][]string
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// MapKeys returns all keys from a map
func (c *CollectionUtil) MapKeys(m map[string]any) []string {
	keys := funk.Keys(m).([]string)
	sort.Strings(keys) // Sort for consistent output
	return keys
}

// MapValues returns all values from a map
func (c *CollectionUtil) MapValues(m map[string]any) []any {
	return funk.Values(m).([]any)
}

// MapMerge merges multiple maps into one (later maps override earlier ones)
func (c *CollectionUtil) MapMerge(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MapFilter filters a map based on a predicate function
func (c *CollectionUtil) MapFilter(m map[string]any, predicate func(string, any) bool) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapPick creates a new map with only the specified keys
func (c *CollectionUtil) MapPick(m map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)
	for _, key := range keys {
		if value, exists := m[key]; exists {
			result[key] = value
		}
	}
	return result
}

// MapOmit creates a new map without the specified keys
func (c *CollectionUtil) MapOmit(m map[string]any, keys ...string) map[string]any {
	omitSet := make(map[string]bool)
	for _, key := range keys {
		omitSet[key] = true
	}

	result := make(map[string]any)
	for k, v := range m {
		if !omitSet[k] {
			result[k] = v
		}
	}
	return result
}

// FindInSlice finds the first element in a slice that matches the predicate
func (c *CollectionUtil) FindInSlice(slice []any, predicate func(any) bool) (any, bool) {
	result := funk.Find(slice, predicate)
	if result == nil {
		return nil, false
	}
	return result, true
}

// SliceIntersection returns elements that exist in both slices
func (c *CollectionUtil) SliceIntersection(slice1, slice2 []string) []string {
	return funk.IntersectString(slice1, slice2)
}

// SliceDifference returns elements that exist in slice1 but not in slice2
func (c *CollectionUtil) SliceDifference(slice1, slice2 []string) []string {
	left, _ := funk.DifferenceString(slice1, slice2)
	return left
}

// SliceUnion returns unique elements from both slices combined
func (c *CollectionUtil) SliceUnion(slice1, slice2 []string) []string {
	combined := append(slice1, slice2...)
	return funk.UniqString(combined)
}
