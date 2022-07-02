package utils

import "time"

// StringSlice creates a slice of interface{} from list of strings
func StringSlice(v ...string) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}

// IntSlice creates a slice of interface{} from list of ints
func IntSlice(v ...int) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}

// Int64Slice creates a slice of interface{} from list of int64
func Int64Slice(v ...int64) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}

// Float32Slice creates a slice of interface{} from list of float32
func Float32Slice(v ...float32) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}

// Float64Slice creates a slice of interface{} from list of float64
func Float64Slice(v ...float64) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}

// TimeSlice creates a slice of interface{} from list of Time
func TimeSlice(v ...time.Time) []interface{} {
	r := make([]interface{}, len(v))
	for i, s := range v {
		r[i] = s
	}
	return r
}
