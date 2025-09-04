package odize

import (
	"testing"
)

// log without formatting
func log(t *testing.T, args ...any) {
	t.Helper()

	t.Error(args...)
	t.FailNow()
}
