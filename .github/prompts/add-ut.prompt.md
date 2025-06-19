Add unit tests for this Go file in the Terraform AzureRM provider. Make sure you cover the important use cases including resource lifecycle operations, Azure API integration, and edge cases. The code must build and pass tests without errors.

## Unit Tests Guidelines for Terraform Provider

- All unit tests must be placed in the same package as the source file with `_test.go` suffix
- Use Go's built-in **testing** framework and **testify** assertions when available
- For Terraform provider resources, focus on **acceptance tests** rather than unit tests
- Unit tests are primarily for utility functions, parsers, validators, and client code
- Follow Go testing conventions and best practices
- Use descriptive test function names: `TestFunctionName_Scenario_ExpectedOutcome`
- Mock external dependencies using interfaces and test doubles
- Test edge cases, error scenarios, and Azure API error conditions
- Follow the same code style and conventions as the main codebase
- If a function or method is not accessible (unexported), consider if it needs testing or if the public API coverage is sufficient

## Acceptance Tests vs Unit Tests

- **Unit Tests**: Use for utility functions, parsers, validators, and isolated logic
- **Acceptance Tests**: Use for full resource lifecycle testing (Create, Read, Update, Delete)
- Acceptance tests interact with real Azure APIs and require Azure credentials
- Unit tests should not require Azure credentials or make external API calls

## Go Testing Patterns

- Use `t.Run()` for subtests to group related test cases
- Use `require` for assertions that should stop test execution on failure
- Use `assert` for assertions that should continue test execution
- Use table-driven tests for testing multiple scenarios
- Use `testify/mock` for mocking complex dependencies
- Test both success and failure paths

## Go Testing Assertions and Patterns

- Use `if got != want` or `require.Equal(t, want, got)` for value comparisons
- Use `require.NoError(t, err)` for error checks that should stop test execution
- Use `assert.Error(t, err)` for expected errors that shouldn't stop execution
- Use `require.NotNil(t, result)` or `require.Nil(t, result)` for nil checks
- Use `require.True(t, condition)` or `require.False(t, condition)` for boolean conditions

## Terraform Provider Specific Testing

### Resource ID Parsing Tests
```go
func TestParseResourceID(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected *ResourceId
        wantErr  bool
    }{
        {
            name:  "valid resource ID",
            input: "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1",
            expected: &ResourceId{
                SubscriptionId: "12345",
                ResourceGroup:  "rg1",
                Name:          "vm1",
            },
            wantErr: false,
        },
        // Add more test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := ParseResourceID(tc.input)
            if tc.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tc.expected, result)
        })
    }
}
```

### Validation Function Tests
```go
func TestValidateResourceName(t *testing.T) {
    validNames := []string{"valid-name", "valid123", "name-with-numbers-123"}
    invalidNames := []string{"", "toolong" + strings.Repeat("a", 100), "-invalid", "invalid-"}
    
    for _, name := range validNames {
        t.Run("valid_"+name, func(t *testing.T) {
            warnings, errors := ValidateResourceName(name, "test")
            require.Empty(t, warnings)
            require.Empty(t, errors)
        })
    }
    
    for _, name := range invalidNames {
        t.Run("invalid_"+name, func(t *testing.T) {
            _, errors := ValidateResourceName(name, "test")
            require.NotEmpty(t, errors)
        })
    }
}
```

### Schema Expand/Flatten Function Tests
```go
func TestExpandResourceConfig(t *testing.T) {
    input := []interface{}{
        map[string]interface{}{
            "name":    "test",
            "enabled": true,
            "tags": map[string]interface{}{
                "Environment": "test",
            },
        },
    }
    
    result := expandResourceConfig(input)
    
    require.NotNil(t, result)
    require.Equal(t, "test", *result.Name)
    require.True(t, *result.Enabled)
    require.Equal(t, "test", *result.Tags["Environment"])
}
```

## Running Tests

- Run unit tests: `go test ./internal/services/servicename/...`
- Run specific test: `go test -run TestFunctionName`
- Run tests with coverage: `go test -cover ./...`
- Run acceptance tests: `make testacc TEST=./internal/services/servicename TESTARGS='-run=TestAccResourceName'`
