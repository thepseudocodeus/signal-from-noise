package assert

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// Assert tests an assumption that must be true for the function to work correctly.
// If the assumption is false, it logs the failure with context and panics.
// This makes the "physics" of the system explicit - what must be true for things to work.
//
// Usage:
//   assert.That(len(items) > 0, "items must not be empty for processing")
//   assert.That(user != nil, "user must be initialized before query")
//
// The comment should clearly state the assumption being tested.
func That(condition bool, assumption string) {
	if !condition {
		// Get caller information for context
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			// Extract just the function name
			parts := strings.Split(funcName, ".")
			funcName = parts[len(parts)-1]
			
			log.Printf("❌ ASSERTION FAILED: %s", assumption)
			log.Printf("   Location: %s:%d in %s()", file, line, funcName)
			panic(fmt.Sprintf("ASSERTION FAILED: %s (at %s:%d)", assumption, file, line))
		}
		log.Printf("❌ ASSERTION FAILED: %s", assumption)
		panic(fmt.Sprintf("ASSERTION FAILED: %s", assumption))
	}
}

// ThatNotNil tests that a pointer/interface is not nil.
// Assumption should state what the value represents and why it must exist.
func ThatNotNil(value interface{}, assumption string) {
	if value == nil {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			parts := strings.Split(funcName, ".")
			funcName = parts[len(parts)-1]
			
			log.Printf("❌ ASSERTION FAILED: %s", assumption)
			log.Printf("   Location: %s:%d in %s()", file, line, funcName)
			panic(fmt.Sprintf("ASSERTION FAILED: %s (at %s:%d)", assumption, file, line))
		}
		log.Printf("❌ ASSERTION FAILED: %s", assumption)
		panic(fmt.Sprintf("ASSERTION FAILED: %s", assumption))
	}
}

// ThatNotEmpty tests that a string/slice is not empty.
func ThatNotEmpty(value interface{}, assumption string) {
	var isEmpty bool
	switch v := value.(type) {
	case string:
		isEmpty = v == ""
	case []interface{}:
		isEmpty = len(v) == 0
	default:
		// Try to get length if it's a slice
		if v == nil {
			isEmpty = true
		}
	}
	
	if isEmpty {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			parts := strings.Split(funcName, ".")
			funcName = parts[len(parts)-1]
			
			log.Printf("❌ ASSERTION FAILED: %s", assumption)
			log.Printf("   Location: %s:%d in %s()", file, line, funcName)
			panic(fmt.Sprintf("ASSERTION FAILED: %s (at %s:%d)", assumption, file, line))
		}
		log.Printf("❌ ASSERTION FAILED: %s", assumption)
		panic(fmt.Sprintf("ASSERTION FAILED: %s", assumption))
	}
}
