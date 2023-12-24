// Package odize lightweight wrapper over the standard testing lib that enables some additional features such as tagging and test lifecycle hooks.
package odize

import (
	"fmt"
	"slices"
	"testing"

	"github.com/code-gorilla-au/env"
)

const (
	// ODIZE_TAGS is the environment variable that is used to filter tests
	ODIZE_TAGS = "ODIZE_TAGS"
)

// NewGroup -  Create a new test group.
//
// If t he ODIZE_TAGS environment variable is set, then only tests with matching tags will be run.
func NewGroup(t *testing.T, tags *[]string) *TestGroup {
	groupTags := tags
	if groupTags == nil {
		groupTags = &[]string{}
	}

	tg := &TestGroup{
		t:        t,
		tags:     *groupTags,
		registry: []TestRegistryEntry{},
		cache:    map[string]struct{}{},
	}

	t.Cleanup(func() {
		if t.Skipped() {
			return
		}

		if !tg.complete && len(tg.registry) > 0 {
			tg.t.Fatalf(fmt.Sprintf("test group \"%s\" did not run. Make sure you use the .Run() method to execute test group", t.Name()))
		}
	})

	return tg
}

// Test - Add a test to the group
func (tg *TestGroup) Test(name string, testFn TestFn) *TestGroup {
	if err := tg.registerTest(name, testFn); err != nil {
		tg.errors.Append(err)
	}

	return tg
}

// BeforeEach - Run before each test
func (tg *TestGroup) BeforeEach(fn func()) {
	tg.beforeEach = fn
}

// BeforeAll - Run before all tests
func (tg *TestGroup) BeforeAll(fn func()) {
	tg.beforeAll = fn
}

// AfterEach - Run after each test
func (tg *TestGroup) AfterEach(fn func()) {
	tg.afterEach = fn
}

// AfterAll - Run after all tests
func (tg *TestGroup) AfterAll(fn func()) {
	tg.afterAll = fn
}

// Run - Run all tests within a group.If the ODIZE_TAGS environment variable is set, then only tests with matching tags will be run.
//
// If errors are encountered, tests will not run.
func (tg *TestGroup) Run() error {
	if tg.errors.Len() > 0 {
		tg.complete = true
		return &tg.errors
	}

	if shouldSkipTests(tg.tags) {
		tg.skipped = true
		tg.t.Skip("Skipping test group ", tg.t.Name())
		return nil
	}

	tg.sanitiseLifecycle()

	tg.beforeAll()

	for _, entry := range tg.registry {
		tg.beforeEach()
		tg.t.Run(entry.name, entry.fn)
		tg.afterEach()
	}

	tg.afterAll()

	tg.complete = true

	return nil
}

// registerTest registers a test to the group. Do not overwrite existing tests.
func (tg *TestGroup) registerTest(name string, testFn TestFn) error {
	if _, ok := tg.cache[name]; ok {
		return fmt.Errorf(fmt.Sprintf("test already exists: %s", name))
	}

	tg.cache[name] = struct{}{}
	tg.registry = append(tg.registry, TestRegistryEntry{
		name: name,
		fn:   testFn,
	})
	return nil
}

// sanitiseLifecycle ensures that the lifecycle functions are not nil by adding no op funcs
func (tg *TestGroup) sanitiseLifecycle() {
	if tg.beforeAll == nil {
		tg.beforeAll = func() {}
	}

	if tg.beforeEach == nil {
		tg.beforeEach = func() {}
	}

	if tg.afterAll == nil {
		tg.afterAll = func() {}
	}

	if tg.afterEach == nil {
		tg.afterEach = func() {}
	}

}

// shouldSkipTests checks if the test group should be skipped based on environment tags
func shouldSkipTests(groupTags []string) bool {
	tags := env.GetAsSlice(ODIZE_TAGS, ",")

	if len(groupTags) == 0 || len(tags) == 0 {
		// run all tests
		return false
	}

	for _, groupTag := range groupTags {
		if slices.Contains(tags, groupTag) {
			return false
		}
	}

	return true
}
