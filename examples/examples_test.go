package examples

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestBasicUnitTestExample(t *testing.T) {
	// will only run on unit test tag
	group := odize.NewGroup(t, &[]string{"unit"})

	err := group.
		Test("should pass", func(t *testing.T) {
			odize.AssertTrue(t, true)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestBasicUnitTestExampleWithBeforeEach(t *testing.T) {
	// will always run
	group := odize.NewGroup(t, nil)

	type UserEntity struct {
		Name string
		Age  int
	}

	seedAge := 1
	var user UserEntity

	group.BeforeEach(func() {
		seedAge++
		user = UserEntity{
			Name: "John",
			Age:  seedAge,
		}
	})

	err := group.
		Test("user age should equal 2", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 3", func(t *testing.T) {
			odize.AssertEqual(t, 3, user.Age)
		}).
		Test("user age should equal 4", func(t *testing.T) {
			odize.AssertEqual(t, 4, user.Age)
		}).
		Test("user age should equal 5", func(t *testing.T) {
			odize.AssertEqual(t, 5, user.Age)
		}).
		Test("user age should equal 6", func(t *testing.T) {
			odize.AssertEqual(t, 6, user.Age)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestBasicUnitTestExampleWithBeforeAndAfterEach(t *testing.T) {
	// will always run
	group := odize.NewGroup(t, nil)

	type UserEntity struct {
		Name string
		Age  int
	}

	seedAge := 1
	var user UserEntity

	group.BeforeEach(func() {
		seedAge++
		user = UserEntity{
			Name: "John",
			Age:  seedAge,
		}
	})

	group.AfterEach(func() {
		user = UserEntity{}
	})

	err := group.
		Test("user age should equal 2", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 3", func(t *testing.T) {
			odize.AssertEqual(t, 3, user.Age)
		}).
		Test("user age should equal 4", func(t *testing.T) {
			odize.AssertEqual(t, 4, user.Age)
		}).
		Test("user age should equal 5", func(t *testing.T) {
			odize.AssertEqual(t, 5, user.Age)
		}).
		Test("user age should equal 6", func(t *testing.T) {
			odize.AssertEqual(t, 6, user.Age)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestBasicUnitTestExampleWithResetState(t *testing.T) {
	// will always run
	group := odize.NewGroup(t, nil)

	type UserEntity struct {
		Name string
		Age  int
	}

	seedAge := 1
	var user UserEntity

	group.BeforeEach(func() {
		seedAge++
		user = UserEntity{
			Name: "John",
			Age:  seedAge,
		}
	})

	group.AfterEach(func() {
		seedAge = 1
		user = UserEntity{}
	})

	err := group.
		Test("user age should equal 2 on first run", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 2 on second run", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 2 on third run", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 2 on forth run", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 2 on firth run", func(t *testing.T) {
			odize.AssertEqual(t, 2, user.Age)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestSkipGroup(t *testing.T) {
	e2e := odize.NewGroup(t, &[]string{"e2e"})
	err := e2e.Test("should not run", func(t *testing.T) {
		t.Error("should not run")
	}).
		Run()

	odize.AssertNoError(t, err)
}

func TestSkipTestWithinGroup(t *testing.T) {
	group := odize.NewGroup(t, nil)
	err := group.
		Test("should equal 1", func(t *testing.T) {
			odize.AssertEqual(t, 1, 1)
		}).
		Test("should equal 2", func(t *testing.T) {
			odize.AssertEqual(t, 2, 2)
		}).
		Test("should skip this test", func(t *testing.T) {
			odize.AssertEqual(t, 3, 3)
		}, odize.Skip()).
		Run()

	odize.AssertNoError(t, err)
}

func TestRunOnlyOneWithinGroup(t *testing.T) {
	group := odize.NewGroup(t, nil)
	err := group.
		Test("should equal 1", func(t *testing.T) {
			odize.AssertEqual(t, 1, 1)
		}).
		Test("should equal 2", func(t *testing.T) {
			odize.AssertEqual(t, 2, 2)
		}).
		Test("should only run this test", func(t *testing.T) {
			odize.AssertEqual(t, 3, 3)
		}, odize.Only()).
		Run()

	odize.AssertNoError(t, err)
}
