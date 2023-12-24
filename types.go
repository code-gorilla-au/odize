package odize

import (
	"testing"
)

// TestGroup - Group tests together, contains lifecycle context.
type TestGroup struct {
	t          *testing.T
	beforeAll  func()
	beforeEach func()
	afterEach  func()
	afterAll   func()
	groupTags  []string
	envTags    []string
	skipped    bool
	complete   bool
	registry   []TestRegistryEntry
	cache      map[string]struct{}
	errors     ErrorList
}

// TestFn - Test function
type TestFn = func(t *testing.T)

// TestRegistryEntry - Test name and function to execute on run
type TestRegistryEntry struct {
	// Name of the test
	name string
	// Test function to execute with context
	fn TestFn
}

// ErrorList - keep track of a number of errors
type ErrorList struct {
	errors []error
}
