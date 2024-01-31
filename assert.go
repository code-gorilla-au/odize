package odize

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Assert value is nil
//
// Example:
//
//	AssertNil(t, myValue)
func AssertNil(t *testing.T, value any) {
	t.Helper()

	if !isNil(value) {
		log(t, decorateDiff("<nil>", value))
	}
}

// AssertTrue checks if value is true
//
// Example:
//
//	AssertTrue(t, methodReturnsTrue())
func AssertTrue(t *testing.T, value bool) {
	t.Helper()

	if !value {
		log(t, decorateDiff(true, value))
	}
}

// AssertFalse checks if value is true
//
// Example:
//
//	AssertFalse(t, methodReturnsFalse())
func AssertFalse(t *testing.T, value bool) {
	t.Helper()

	if value {
		log(t, decorateDiff(false, value))
	}
}

// AssertNoError checks if error is nil
//
// Example:
//
//	AssertNoError(t, err)
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		log(t, decorateDiff("<nil>", err))
	}
}

// AssertError checks if error is not nil
//
// Example:
//
//	AssertError(t, err)
func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		log(t, decorateDiff("<error>", err))
	}
}

// AssertEqual checks if two values are equal
//
// Example:
//
//	AssertEqual(t, "a", "b")
func AssertEqual(t *testing.T, expected any, actual any) {
	t.Helper()

	if !isEqual(expected, actual) {
		logf(t, decorateDiff(expected, actual))
	}
}

// isNil - Check if a value is nil
func isNil(value any) bool {

	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()
	// check chan, func, interface, map, pointer, or slice value is nil
	if (kind >= reflect.Chan && kind <= reflect.Slice) && v.IsNil() {
		return true
	}

	return false
}

// isEqual - Check if two values are equal
func isEqual(expected any, actual any) bool {
	if isNil(expected) && isNil(actual) {
		return true
	}

	if isNil(expected) || isNil(actual) {
		return false
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	valueExpected := reflect.ValueOf(expected)
	valueActual := reflect.ValueOf(actual)

	return valueExpected == valueActual
}

func decorateDiff(expected, actual any) string {

	buf := new(bytes.Buffer)

	exp := decorateBlock("Expected", fmt.Sprintf("%v", expected), "+")
	got := decorateBlock("Got", fmt.Sprintf("%v", actual), "-")

	buf.WriteString(exp)
	buf.WriteString(got)

	return buf.String()

}

func decorateBlock(label string, content string, lineDecoration string) string {
	buf := new(bytes.Buffer)

	buf.WriteByte('\n')

	buf.WriteString(fmt.Sprintf("%s:\n", label))
	lines := strings.Split(content, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}

	for i, line := range lines {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("%s\t%v", lineDecoration, line))
	}

	buf.WriteByte('\n')

	return buf.String()
}
