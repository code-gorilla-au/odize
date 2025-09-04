package odize

import "errors"

var (
	ErrTestAlreadyExists = errors.New("test already exists")
)

// Error - return error string
func (e *ListError) Error() string {
	result := ""
	for _, item := range e.errors {
		result += item.Error() + ", "
	}

	return result
}

// Unwrap - returns list of errors
func (e *ListError) Unwrap() []error {
	return e.errors
}

// Append - append error to list
func (e *ListError) Append(err error) {
	e.errors = append(e.errors, err)
}

func (e *ListError) Len() int {
	return len(e.errors)
}

// Pop - remove the first error from the list and return
func (e *ListError) Pop() error {
	if len(e.errors) == 0 {
		return nil
	}

	first := e.errors[0]
	rest := e.errors[1:]
	e.errors = rest
	return first
}
