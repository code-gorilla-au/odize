package odize

// Skip - Skip this test
func Skip() TestFuncOpts {
	return func(to *TestOpts) {
		to.Skip = true
	}
}

// Only - Only run this test.
// If multiple tests are marked as only,
// the group will run only the tests marked as only
func Only() TestFuncOpts {
	return func(to *TestOpts) {
		to.Only = true
	}
}
