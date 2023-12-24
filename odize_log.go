package odize

import (
	"fmt"
	"testing"
)

// log without formatting
func log(t *testing.T, args ...any) {
	t.Helper()

	t.Error(args...)
	t.FailNow()
}

// logf log with formatting
func logf(t *testing.T, format string, args ...any) {
	t.Helper()
	log(t, fmt.Sprintf(format, args...))
}
