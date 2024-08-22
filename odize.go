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
	// ENV variable declared in pipelines such as Github Actions
	ENV_CI = "CI"
)

var (
	ErrTestOptionNotAllowedInCI = fmt.Errorf("test option 'Only' not allowed in CI environment")
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
		t:         t,
		groupTags: *groupTags,
		envTags:   env.GetAsSlice(ODIZE_TAGS, ","),
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
		isCIEnv:   env.GetAsBool(ENV_CI),
	}

	tg.registerCleanupTasks()

	return tg
}

// Test - Add a test to the group
func (tg *TestGroup) Test(name string, testFn TestFn, options ...TestFuncOpts) *TestGroup {
	testOpts := TestOpts{}
	for _, opt := range options {
		opt(&testOpts)
	}

	if err := tg.registerTest(name, testFn, testOpts); err != nil {
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
	tg.t.Helper()

	if tg.errors.Len() > 0 {
		tg.complete = true
		return &tg.errors
	}

	if shouldSkipTests(tg.groupTags, tg.envTags) {
		tg.skipped = true
		tg.t.Skip("Skipping test group ", tg.t.Name())
		return nil
	}

	tg.sanitiseLifecycle()

	tg.beforeAll()

	entries, err := filterExecutableTests(tg.t, tg.isCIEnv, tg.registry)
	if err != nil {
		// Stop Run, suite is in an invalid state
		tg.complete = true
		return fmt.Errorf("Test group \"%s\" error: %w", tg.t.Name(), err)
	}

	for _, entry := range entries {
		tg.beforeEach()
		tg.t.Run(entry.name, entry.fn)
		tg.afterEach()
	}

	tg.afterAll()

	tg.complete = true

	return nil
}

// registerTest registers a test to the group. Do not overwrite existing tests.
func (tg *TestGroup) registerTest(name string, testFn TestFn, options TestOpts) error {
	if _, ok := tg.cache[name]; ok {
		return fmt.Errorf(fmt.Sprintf("test already exists: %s", name))
	}

	tg.cache[name] = struct{}{}
	tg.registry = append(tg.registry, TestRegistryEntry{
		name:    name,
		fn:      testFn,
		options: options,
	})
	return nil
}

// registerCleanupTasks registers cleanup tasks to ensure that the test group is run
func (tg *TestGroup) registerCleanupTasks() {
	tg.t.Helper()

	tg.t.Cleanup(func() {
		tg.t.Helper()
		if tg.t.Skipped() {
			return
		}

		if !tg.complete && len(tg.registry) > 0 {
			tg.t.Fatalf(fmt.Sprintf("test group \"%s\" did not run. Make sure you use the .Run() method to execute test group", tg.t.Name()))
		}
	})
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
func shouldSkipTests(groupTags []string, envTags []string) bool {

	fmt.Println("groupTags", groupTags, "len", len(groupTags))
	fmt.Println("envTags", envTags, "len", len(envTags))

	if len(groupTags) == 0 && len(envTags) == 0 {
		// run all tests
		fmt.Println("running test")
		return false
	}

	for _, groupTag := range groupTags {
		if slices.Contains(envTags, groupTag) {
			return false
		}
	}

	fmt.Println("hitting here")

	return true
}

// filterExecutableTests filters tests that are executable within the test group
// Note that test option 'Only' is only used for debugging tests, and should not be used in a CI env.
func filterExecutableTests(t *testing.T, isCIEnv bool, tests []TestRegistryEntry) ([]TestRegistryEntry, error) {
	filtered, err := filterOnlyAllowedTests(t, isCIEnv, tests)
	if err != nil {
		return filtered, err
	}

	if len(filtered) > 0 {
		// if there are tests that are marked as only, then return 'only' those tests
		return filtered, nil
	}

	for _, test := range tests {
		if test.options.Skip {
			filtered = append(filtered, TestRegistryEntry{
				name: test.name,
				fn: func(t *testing.T) {
					t.Skip("skipping test ", test.name)
				},
			})

			continue
		}

		filtered = append(filtered, test)
	}

	return filtered, nil
}

// filterOnlyAllowedTests filters tests that are marked as only within a test group
// If the framework detects that the test is running under a CI environment and the group has tests with 'Only', then it will return an error
func filterOnlyAllowedTests(t *testing.T, isCIEnv bool, tests []TestRegistryEntry) ([]TestRegistryEntry, error) {
	filtered := []TestRegistryEntry{}

	for _, test := range tests {
		if test.options.Only && isCIEnv {
			return filtered, ErrTestOptionNotAllowedInCI
		}

		if test.options.Only {
			filtered = append(filtered, test)
		}
	}

	return filtered, nil
}
