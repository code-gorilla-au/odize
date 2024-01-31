package odize

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecorateBlock(t *testing.T) {
	result := decorateBlock("test", "content", "++")
	group := NewGroup(t, nil)
	err := group.
		Test("should contain label", func(t *testing.T) {
			AssertTrue(t, strings.Contains(result, "test"))
		}).
		Test("should contain content", func(t *testing.T) {
			AssertTrue(t, strings.Contains(result, "content"))
		}).
		Test("should contain line decorator", func(t *testing.T) {
			AssertTrue(t, strings.Contains(result, "++"))
		}).
		Run()

	AssertNoError(t, err)

}

func TestDecorateDiff(t *testing.T) {
	result := decorateDiff("expected", "actual")
	group := NewGroup(t, nil)
	err := group.
		Test("should contain expected", func(t *testing.T) {
			AssertTrue(t, strings.Contains(result, "expected"))
		}).
		Test("should contain actual", func(t *testing.T) {
			AssertTrue(t, strings.Contains(result, "actual"))
		}).
		Run()

	AssertNoError(t, err)
}

func TestAssertNil(t *testing.T) {
	group := NewGroup(t, nil)

	var nilString *string
	var nilInt *int
	var nilFloat *float32
	var nilBool *bool
	var nilSlice []string

	err := group.
		Test("should pass nil string", func(t *testing.T) {
			AssertNil(t, nilString)
		}).
		Test("should pass nil int", func(t *testing.T) {
			AssertNil(t, nilInt)
		}).
		Test("should pass nil float", func(t *testing.T) {
			AssertNil(t, nilFloat)
		}).
		Test("should pass nil bool", func(t *testing.T) {
			AssertNil(t, nilBool)
		}).
		Test("should pass nil empty slice", func(t *testing.T) {
			AssertNil(t, nilSlice)
		}).
		Run()

	AssertNoError(t, err)
}

func TestAssertEqual(t *testing.T) {
	testStruct := struct{ name string }{name: "test"}

	group := NewGroup(t, nil)

	err := group.
		Test("should pass equal strings", func(t *testing.T) {
			AssertEqual(t, "a", "a")
		}).
		Test("should pass equal ints", func(t *testing.T) {
			AssertEqual(t, 1, 1)
		}).
		Test("should pass equal floats", func(t *testing.T) {
			AssertEqual(t, 1.1, 1.1)
		}).
		Test("should pass equal bool", func(t *testing.T) {
			AssertEqual(t, true, true)
		}).
		Test("should pass identical", func(t *testing.T) {
			AssertEqual(t, testStruct, testStruct)
		}).
		Test("should pass same struct signature", func(t *testing.T) {
			AssertEqual(t, testStruct, struct{ name string }{name: "test"})
		}).
		Test("should pass same slice signature", func(t *testing.T) {
			AssertEqual(t, []string{
				"a",
			}, []string{
				"a",
			})
		}).
		Test("should pass two empty slices", func(t *testing.T) {
			AssertEqual(t, []string{}, []string{})
		}).
		Run()

	AssertNoError(t, err)
}

func TestAssertTrue(t *testing.T) {
	group := NewGroup(t, nil)
	err := group.
		Test("should pass on literal", func(t *testing.T) {
			AssertTrue(t, true)
		}).
		Test("should pass on string assertion", func(t *testing.T) {
			AssertTrue(t, strings.Contains("hello", "hello"))
		}).
		Run()

	AssertNoError(t, err)
}

func TestAssertFalse(t *testing.T) {
	AssertFalse(t, false)
}

func TestAssertError(t *testing.T) {
	AssertError(t, fmt.Errorf("test"))
}
