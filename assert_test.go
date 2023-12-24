package odize

import (
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
