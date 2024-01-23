package odize

import (
	"errors"
	"testing"
)

func TestUnitNoEnvVarShouldRunAll(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	err := unit.
		Test("test should equal one", testShouldEqualOne).
		Test("test should equal three", testShouldEqualThree).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func TestUnitEnvVarShouldRunAll(t *testing.T) {

	t.Setenv(ODIZE_TAGS, "unit")

	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	err := unit.
		Test("testShouldEqualOne", testShouldEqualOne).
		Test("testShouldEqualThree", testShouldEqualThree).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)

}

func TestUnitEnvVarNonMatchShouldSkipAndPass(t *testing.T) {
	t.Setenv(ODIZE_TAGS, "unit")

	unit := &TestGroup{
		t:         t,
		groupTags: []string{"integration"},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	err := unit.
		Test("testShouldFail", testShouldFail).
		Run()
	AssertNoError(t, err)

}

func TestBeforeAll(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	increment := 0

	unit.BeforeAll(func() {
		increment++
	})

	AssertEqual(t, 0, increment)

	err := unit.
		Test("should increment", func(t *testing.T) {
			AssertEqual(t, 1, increment)
		}).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func TestAfterAll(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	increment := 0

	unit.AfterAll(func() {
		increment++
	})

	AssertEqual(t, 0, increment)

	err := unit.Test("should not increment", func(t *testing.T) {
		AssertEqual(t, 0, increment)
	}).Run()

	AssertNoError(t, err)
	AssertEqual(t, 1, increment)
}

func TestBeforeEach(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	increment := 0

	unit.BeforeEach(func() {
		increment++
	})

	AssertEqual(t, 0, increment)

	err := unit.
		Test("should equal 1", func(t *testing.T) {
			AssertEqual(t, 1, increment)
		}).
		Test("should equal 2", func(t *testing.T) {
			AssertEqual(t, 2, increment)
		}).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func TestAfterEach(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	increment := 0

	unit.AfterEach(func() {
		increment++
	})

	AssertEqual(t, 0, increment)

	err := unit.
		Test("should equal 0", func(t *testing.T) {
			AssertEqual(t, 0, increment)
		}).
		Test("should equal 1", func(t *testing.T) {
			AssertEqual(t, 1, increment)
		}).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func TestTestFuncWithNamedFuncs(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	unit.
		Test("testShouldEqualOne", testShouldEqualOne).
		Test("testShouldEqualThree", testShouldEqualThree)

	err := unit.Run()
	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func TestTestFuncWithAnonymousFuncs(t *testing.T) {
	unit := &TestGroup{
		t:         t,
		groupTags: []string{},
		registry:  []TestRegistryEntry{},
		cache:     map[string]struct{}{},
	}

	err := unit.
		Test("testShouldEqualOne", func(t *testing.T) {
			AssertEqual(t, 1, 1)
		}).
		Test("testShouldEqualThree", func(t *testing.T) {
			AssertEqual(t, 2, 2)
		}).
		Run()

	AssertNoError(t, err)
	AssertEqual(t, false, unit.skipped)
}

func testShouldEqualOne(t *testing.T) {
	AssertEqual(t, 1, 1)
}

func testShouldEqualThree(t *testing.T) {
	AssertEqual(t, 3, 3)
}

func testShouldFail(t *testing.T) {
	AssertEqual(t, 1, 2)
}

func TestShouldSkipTests(t *testing.T) {
	group := NewGroup(t, nil)

	err := group.
		Test("should not skip if no env vars", func(t *testing.T) {
			result := shouldSkipTests([]string{"unit"}, []string{})
			AssertTrue(t, result)
		}).
		Test("should skip if no env var does not match", func(t *testing.T) {
			result := shouldSkipTests([]string{"unit"}, []string{"system"})
			AssertTrue(t, result)
		}).
		Test("should not skip if env var does match", func(t *testing.T) {
			result := shouldSkipTests([]string{"unit"}, []string{"unit"})
			AssertFalse(t, result)
		}).
		Test("should not skip if env var present and group has no tags", func(t *testing.T) {
			result := shouldSkipTests([]string{}, []string{"unit"})
			AssertTrue(t, result)
		}).
		Run()

	AssertNoError(t, err)
}

func TestRegisterCleanupTaskShouldNotFailIfSkipped(t *testing.T) {
	tg := TestGroup{
		t:     t,
		cache: map[string]struct{}{},
	}

	err := tg.registerTest("test", func(t *testing.T) {
		t.Skip()
	}, TestOpts{})
	AssertNoError(t, err)

	err = tg.Run()
	AssertNoError(t, err)

	tg.registerCleanupTasks()
}

func TestRegisterCleanupTaskShouldNotFailIfComplete(t *testing.T) {
	tg := TestGroup{
		t:        t,
		cache:    map[string]struct{}{},
		complete: true,
	}

	err := tg.registerTest("test", func(t *testing.T) {}, TestOpts{})
	AssertNoError(t, err)

	tg.registerCleanupTasks()
}

func TestOptionSkip(t *testing.T) {
	tg := NewGroup(t, nil)

	testExecute := false

	err := tg.
		Test("should not execute test", func(t *testing.T) {
			testExecute = true
		}, Skip()).
		Run()
	AssertNoError(t, err)
	AssertFalse(t, testExecute)
}

func TestOptionOnly(t *testing.T) {
	tg := NewGroup(t, nil)

	testCall := 0

	err := tg.
		Test("should execute test ", func(t *testing.T) {
			testCall++
		}, Only()).
		Test("should not execute test", func(t *testing.T) {
			testCall++
		}).
		Run()
	AssertNoError(t, err)
	AssertEqual(t, 1, testCall)
}

func TestCITest(t *testing.T) {
	t.Setenv(ENV_CI, "true")

	defer func() {
		t.Setenv(ENV_CI, "false")
	}()

	tg := NewGroup(t, nil)

	err := tg.
		Test("should fail group if ENV_CI is true", func(t *testing.T) {
			AssertTrue(t, true)
		}, Only()).
		Test("should not execute test", func(t *testing.T) {
			AssertTrue(t, false)
		}).
		Run()
	AssertTrue(t, errors.Is(err, ErrTestOptionNotAllowedInCI))
}

func TestCITestNoOnly(t *testing.T) {
	t.Setenv(ENV_CI, "true")

	defer func() {
		t.Setenv(ENV_CI, "false")
	}()

	tg := NewGroup(t, nil)

	err := tg.
		Test("should pass", func(t *testing.T) {
			AssertTrue(t, true)
		}).
		Run()
	AssertNoError(t, err)
}
