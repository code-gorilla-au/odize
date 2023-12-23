# odize

Testing, [supercharged](https://www.yourdictionary.com/odize)! odize is a lightweight wrapper over the standard testing lib that enables some additional features.

## Motivation

Bringing the JS ecosystem to golang! Jokes aside, I wanted to remove boilerplate code from the tests. 

The golang testing standard lib is more than capable for most cases, it's preferable to default to the standard lib where possible.

If what you're working on needs to be able to filter tests by tag, have more granular setup / teardown code, please consider odize 


[toc]

## Features

- Lightweight wrapper on the go standard library
    - Remove boilerplate code
    - Same reports / output
    - Same flags
- Lifecycle hooks
    - BeforeAll - run before all tests
    - BeforeEach - run before each test
    - AfterAll - run after all test
    - AfterEach - run after each test
- Test grouping
    - Group tests by tag to enable test filtering
- Assertions built in
    - AssertEqual
    - AssertTrue
    - AssertNoError
    - AssertNil

## Basic usage

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
		Run(t)

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
--- PASS: TestScenarioOne (0.00s)
    --- PASS: TestScenarioOne/user_age_should_equal_2 (0.00s)
    --- PASS: TestScenarioOne/user_age_should_equal_3 (0.00s)
    --- PASS: TestScenarioOne/user_age_should_equal_4 (0.00s)
    --- PASS: TestScenarioOne/user_age_should_equal_5 (0.00s)
    --- PASS: TestScenarioOne/user_age_should_equal_6 (0.00s)
PASS
ok      github.com/code-gorilla-au/odize     0.118s

```

## Filtering tets

Provide the specific environment variable with values `ODIZE_TAGS="unit"`. 

Multiple tags can be passed with a comma `,` delimiter `ODIZE_TAGS="unit,system"`

### Create group

create filtered group

```golang
func TestScenarioTwo(t *testing.T) {
	group := odize.NewGroup(t, &[]string{"integration"})

/** omit rest of the code **/

```

### Run test

```bash
# only run unit tests
ODIZE_TAGS="unit" go test --short -v -cover  -failfast ./... 

```

```bash
--- PASS: TestScenarioOne/should_equal_1 (0.00s)
=== RUN   TestScenarioTwo
    odize.go:55: Skipping test group  TestScenarioTwo
--- SKIP: TestScenarioTwo (0.00s)
```




