package logging

import (
	"fmt"
	"log"
	"time"
)

// Operation represents a unit of work being logged
type Operation struct {
	name        string
	startTime   time.Time
	context     map[string]interface{}
}

// StartOperation begins logging a new operation with context
// This makes the "experiment" explicit - we're testing assumptions about what should happen
func StartOperation(name string, context map[string]interface{}) *Operation {
	op := &Operation{
		name:      name,
		startTime: time.Now(),
		context:   context,
	}
	
	log.Printf("▶️  START: %s", name)
	if len(context) > 0 {
		for key, value := range context {
			log.Printf("   %s: %v", key, value)
		}
	}
	
	return op
}

// EndOperation logs the completion of an operation with duration
func (op *Operation) EndOperation() {
	duration := time.Since(op.startTime)
	log.Printf("✅ END: %s (duration: %v)", op.name, duration)
}

// EndOperationWithResult logs completion with a result
func (op *Operation) EndOperationWithResult(result interface{}) {
	duration := time.Since(op.startTime)
	log.Printf("✅ END: %s (duration: %v)", op.name, duration)
	log.Printf("   Result: %v", result)
}

// LogAssumption logs an assumption being tested (before assertion)
// This documents the "hypothesis" - what we assume must be true
func LogAssumption(assumption string, context map[string]interface{}) {
	log.Printf("🔬 TESTING ASSUMPTION: %s", assumption)
	if len(context) > 0 {
		for key, value := range context {
			log.Printf("   %s: %v", key, value)
		}
	}
}

// LogInvariant logs that an invariant is being checked
// Invariants are properties that must always be true
func LogInvariant(invariant string, value interface{}) {
	log.Printf("🔒 CHECKING INVARIANT: %s = %v", invariant, value)
}

// LogState logs the current state at a checkpoint
// Useful for understanding what the system "knows" at this point
func LogState(component string, state map[string]interface{}) {
	log.Printf("📊 STATE [%s]:", component)
	for key, value := range state {
		log.Printf("   %s: %v", key, value)
	}
}

// LogError logs an error with full context
func LogError(operation string, err error, context map[string]interface{}) {
	log.Printf("❌ ERROR in %s: %v", operation, err)
	if len(context) > 0 {
		for key, value := range context {
			log.Printf("   %s: %v", key, value)
		}
	}
}

// LogTransition logs a state transition
// Documents how the system moves from one state to another
func LogTransition(from string, to string, reason string) {
	log.Printf("🔄 TRANSITION: %s → %s (reason: %s)", from, to, reason)
}

// LogQuery logs a database query with parameters
// Makes data access transparent
func LogQuery(query string, params map[string]interface{}) {
	log.Printf("📝 QUERY: %s", query)
	if len(params) > 0 {
		for key, value := range params {
			log.Printf("   %s: %v", key, value)
		}
	}
}

// LogResult logs a query result with count
func LogResult(operation string, count int, additionalInfo map[string]interface{}) {
	log.Printf("📦 RESULT [%s]: %d items", operation, count)
	if len(additionalInfo) > 0 {
		for key, value := range additionalInfo {
			log.Printf("   %s: %v", key, value)
		}
	}
}

// LogCheckpoint logs a checkpoint in execution
// Useful for understanding execution flow
func LogCheckpoint(location string, data map[string]interface{}) {
	log.Printf("📍 CHECKPOINT: %s", location)
	if len(data) > 0 {
		for key, value := range data {
			log.Printf("   %s: %v", key, value)
		}
	}
}

// FormatDuration formats a duration for logging
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
