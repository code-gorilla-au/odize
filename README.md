# odize

Testing, [supercharged](https://www.yourdictionary.com/odize)! odize is a lightweight wrapper over the standard testing lib that enables some additional features.

[![Go Report Card](https://goreportcard.com/badge/github.com/code-gorilla-au/odize)](https://goreportcard.com/report/github.com/code-gorilla-au/odize)
[![Go Reference](https://pkg.go.dev/badge/github.com/code-gorilla-au/odize.svg)](https://pkg.go.dev/github.com/code-gorilla-au/odize)

## Motivation

Bringing the JS ecosystem to golang! Jokes aside, I wanted to remove boilerplate code from the tests while still using the core testing library. Heavy inspiration from [vitest](https://vitest.dev/) and [jest](https://jestjs.io/), odize aims to be a light weight, easy to use test framework on top of the standard library.

The golang testing standard lib is more than capable for most cases, it's preferable to default to the standard lib where possible.

If what you're working on needs to be able to filter tests by tag, have more granular setup / teardown code, please consider odize.


## Features

| Feature | Description |
| ------- | ----------- |
| Powered by std lib |  Lightweight wrapper over the standard testing library, easy plug and play, no need to update your test commands. | 
| Lifecycle hooks | Have granular control in the setup / teardown tests with helper functions: `BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` |
| Test filtering | Run a subset of tests based off either `group tags`, or via `test options`. |
| Assertions | Built in core assertions `AssertEqual`, `AssertTrue`, `AssertFalse`, `AssertNoError`, `AssertError`, `AssertNil` | 

## Basic usage

### Install

```bash
go get github.com/code-gorilla-au/odize@latest
```

### Create your group

Create a test group 

```golang
// Note you can add test tags to filter tests
func TestScenarioOne(t *testing.T) {
	group := odize.NewGroup(t, nil)

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
			AssertEqual(t, 2, user.Age)
		}).
		Test("user age should equal 3", func(t *testing.T) {
			AssertEqual(t, 3, user.Age)
		}).
		Test("user age should equal 4", func(t *testing.T) {
			AssertEqual(t, 4, user.Age)
		}).
		Test("user age should equal 5", func(t *testing.T) {
			AssertEqual(t, 5, user.Age)
		}).
		Test("user age should equal 6", func(t *testing.T) {
			AssertEqual(t, 6, user.Age)
		}).
		Run()

	AssertNoError(t, err)

}

```

### Run test

Run the test command with your normal flags

```bash
go test --short -cover -v -failfast ./...
```
Terminal output

```bash
go test -v --short -cover -failfast ./...
=== RUN   TestDecorateBlock
=== RUN   TestDecorateBlock/should_contain_label
=== RUN   TestDecorateBlock/should_contain_content
=== RUN   TestDecorateBlock/should_contain_line_decorator
--- PASS: TestDecorateBlock (0.00s)
    --- PASS: TestDecorateBlock/should_contain_label (0.00s)
    --- PASS: TestDecorateBlock/should_contain_content (0.00s)
    --- PASS: TestDecorateBlock/should_contain

```

## Lifecycle hooks 

Odize has helper functions that help provide granular setup / teardown helpers for each test within the group.

| Hook | Description |
| ---- | ----------- |
| BeforeAll | Invoke before all tests within a group |
| BeforeEach | Invoke before each test within a group |
| AfterEach | Invoke after each test within a group |
| AfterAll | Invoke after all tests within a group | 

## Test options

Optionally, you are able to provide some test options to a test within a group. This provides fine grain control over the test group, especially when you need to isolate a singular test within a group to debug.

| Option | Description |
| ------ | ----------- |
| Skip	 |	Skip specified test |
| Only   | Within the test group, only run the specified test |


### Providing options to a test

Skip example 

```golang
func TestSkipExample(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should equal 2", func(t *testing.T) {
			result := add(1,1)
			AssertEqual(t, 2, result)
		}).
		Test("should equal 4", func(t *testing.T) {
			result := add(2,2)
			AssertEqual(t, 4, result)
		}).
		Test("should equal 3", func(t *testing.T) {
			// Note this test will be skipped
			result := add(1,2)
			AssertEqual(t, 3, result)
		}, Skip()).
		Run(t)

	AssertNoError(t, err)
}
```


```golang
func TestOnlyExample(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should equal 2", func(t *testing.T) {
			result := add(1,1)
			AssertEqual(t, 2, result)
		}).
		Test("should equal 3", func(t *testing.T) {
			// Note, only this test will be run within this group
			result := add(1,2)
			AssertEqual(t, 3, result)
		}, Only()).
		Run(t)

	AssertNoError(t, err)
}
```

## Filtering tests

Provide the specific environment variable with values `ODIZE_TAGS="unit"`. 

Multiple tags can be passed with a comma `,` delimiter `ODIZE_TAGS="unit,system"`

### Create group

create filtered group

```golang
func TestScenarioTwo(t *testing.T) {
	group := odize.NewGroup(t, &[]string{"integration"})

/** omit rest of the code **/
}

```

### Run test

```bash
# only run unit tests
ODIZE_TAGS="unit" go test --short -v -cover  -failfast ./... 

```

```bash
=== RUN   TestSkipGroup
    unit_test.go:159: Skipping test group  TestSkipGroup
--- SKIP: TestSkipGroup (0.00s)
```



## Examples

See [examples provided](./examples/examples_test.go) for more details.