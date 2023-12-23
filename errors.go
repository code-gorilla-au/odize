package odize

import "fmt"

// Error - return error string
func (e *ErrorList) Error() string {
	result := ""
	for _, item := range e.errors {
		result += fmt.Sprintf("%s, ", item.Error())
	}

	return result
}

// Unwrap - returns list of errors
func (e *ErrorList) Unwrap() []error {
	return e.errors
}

// Append - append error to list
func (e *ErrorList) Append(err error) {
	e.errors = append(e.errors, err)
}

func (e *ErrorList) Len() int {
	return len(e.errors)
}

// Pop - remove the first error from the list and return
func (e *ErrorList) Pop() error {
	if len(e.errors) == 0 {
		return nil
	}

	first := e.errors[0]
	rest := e.errors[1:]
	e.errors = rest
	return first
}
