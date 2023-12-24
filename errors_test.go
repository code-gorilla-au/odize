package odize

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	list := ErrorList{}
	list.Append(errors.New("expected"))

	if list.Error() != "expected, " {
		t.Errorf("expected [%s], got [%v ]", "foo", list.Error())
	}
}

func TestPop(t *testing.T) {
	list := ErrorList{}
	list.Append(errors.New("first"))
	list.Append(errors.New("second"))

	expected := list.Pop()
	if expected.Error() != "first" {
		t.Errorf("expected %s, got %v ", "first", expected.Error())
	}
}
func TestPop_no_errors(t *testing.T) {
	list := ErrorList{}

	expected := list.Pop()
	if expected != nil {
		t.Errorf("expected %s, got %v ", "first", expected.Error())
	}
}
