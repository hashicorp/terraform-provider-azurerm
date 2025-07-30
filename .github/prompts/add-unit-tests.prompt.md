Add unit tests for this Go file in the Terraform AzureRM provider. Make sure you cover the important use cases including resource lifecycle operations, Azure API integration, and edge cases. The code must build and pass tests without errors.

## Unit Tests Guidelines for Terraform Provider

- All unit tests must be placed in the same package as the source file with _test.go suffix
- **Package Naming**: Use external test package convention with _test suffix (e.g., package cdn_test for cdn package)
- Use Go's built-in **testing** framework and **testify** assertions when available
- For Terraform provider resources, focus on **acceptance tests** rather than unit tests
- Unit tests are primarily for utility functions, parsers, validators, and client code
- Follow Go testing conventions and best practices
- Use descriptive test function names: TestFunctionName_Scenario_ExpectedOutcome
- Mock external dependencies using interfaces and test doubles
- Test edge cases, error scenarios, and Azure API error conditions
- Follow the same code style and conventions as the main codebase
- If a function or method is not accessible (unexported), consider if it needs testing or if the public API coverage is sufficient

## Implementation Approach Considerations

This provider supports two implementation approaches. **Unit testing patterns should be appropriate for the implementation approach being tested.**

### Typed Resource Implementation Testing
**For resources using the internal/sdk framework**

- Test type-safe model structures with 	fschema tags
- Test receiver methods on resource struct types
- Test metadata.Decode() and metadata.Encode() patterns
- Test structured error handling with metadata methods
- Test resource interfaces (sdk.Resource, sdk.ResourceWithUpdate, etc.)
- Test IDValidationFunc() method implementations

### Untyped Resource Implementation Testing
**For resources using traditional Plugin SDK patterns**

- Test function-based CRUD patterns
- Test direct schema manipulation patterns
- Test traditional client initialization
- Test f.ImportAsExistsError() and state manipulation patterns
- Test d.Set() and d.Get() state management

### Common Testing Patterns (Both Approaches)
- Resource ID parsing and validation remain consistent
- Azure API integration patterns are the same
- Validation functions follow the same patterns
- Error handling standards apply to both approaches

## Acceptance Tests vs Unit Tests

- **Unit Tests**: Use for utility functions, parsers, validators, and isolated logic
- **Acceptance Tests**: Use for full resource lifecycle testing (Create, Read, Update, Delete)
- Acceptance tests interact with real Azure APIs and require Azure credentials
- Unit tests should not require Azure credentials or make external API calls
- **Testing consistency**: Both implementation approaches use identical acceptance test patterns

## Go Testing Patterns

- Use .Run() for subtests to group related test cases
- Use require for assertions that should stop test execution on failure
- Use Assert for assertions that should continue test execution
- Use table-driven tests for testing multiple scenarios
- Use testify/mock for mocking complex dependencies
- Test both success and failure paths

## Go Testing Assertions and Patterns

- Use if got != want or Require.Equal(t, want, got) for value comparisons
- Use Require.NoError(t, err) for error checks that should stop test execution
- Use Assert.Error(t, err) for expected errors that shouldn't stop execution
- Use Require.NotNil(t, result) or Require.Nil(t, result) for nil checks
- Use Require.True(t, condition) or Require.False(t, condition) for boolean conditions

## Implementation-Specific Testing Patterns

### typed resource Resource Testing Patterns

### Model Structure Testing
```go
package servicename_test

import (
    "testing"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename"
    "github.com/stretchr/testify/require"
)

func TestServiceNameResourceModel_TfschemaValidation(t *testing.T) {
    model := servicename.ServiceNameResourceModel{
        Name:          "test-resource",
        ResourceGroup: "test-rg",
        Location:      "West Europe",
        Enabled:       true,
    }

    // Test that model fields have correct tfschema tags
    // This would typically be tested through the resource's Arguments() and Attributes() methods
    require.Equal(t, "test-resource", model.Name)
    require.Equal(t, "test-rg", model.ResourceGroup)
}
```

### Resource Interface Testing
```go
func TestServiceNameResource_Interfaces(t *testing.T) {
    resource := servicename.ServiceNameResource{}

    // Test that resource implements required interfaces
    var _ sdk.Resource = resource
    var _ sdk.ResourceWithUpdate = resource

    // Test ResourceType method
    require.Equal(t, "azurerm_service_name", resource.ResourceType())

    // Test ModelObject method
    model := resource.ModelObject()
    require.IsType(t, &servicename.ServiceNameResourceModel{}, model)
}
```

### Typed Resource Error Handling Testing
```go
func TestServiceNameResource_ErrorHandling(t *testing.T) {
    // Test metadata error patterns
    testCases := []struct {
        name          string
        setupMock     func() *servicename.MockMetadata
        expectedError string
    }{
        {
            name: "decode error handling",
            setupMock: func() *servicename.MockMetadata {
                mock := &servicename.MockMetadata{}
                mock.On("Decode", mock.Anything).Return(errors.New("decode failed"))
                return mock
            },
            expectedError: "decoding:",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test error handling patterns specific to typed resource
        })
    }
}
```

### untyped Plugin SDK Testing Patterns

### Function-Based CRUD Testing
```go
package servicename_test

func TestResourceServiceNameCreate_ValidationPatterns(t *testing.T) {
    // Test traditional function-based patterns
    testCases := []struct {
        name     string
        input    map[string]interface{}
        wantErr  bool
        errorMsg string
    }{
        {
            name: "valid input",
            input: map[string]interface{}{
                "name":                "test-resource",
                "resource_group_name": "test-rg",
                "location":           "West Europe",
            },
            wantErr: false,
        },
        {
            name: "missing required field",
            input: map[string]interface{}{
                "resource_group_name": "test-rg",
                "location":           "West Europe",
            },
            wantErr:  true,
            errorMsg: "name is required",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test untyped validation patterns
        })
    }
}
```

### Untyped Error Handling Testing
```go
func TestResourceServiceName_untypedErrorHandling(t *testing.T) {
    testCases := []struct {
        name           string
        mockResponse   *http.Response
        expectedAction string
    }{
        {
            name:           "resource not found",
            mockResponse:   &http.Response{StatusCode: 404},
            expectedAction: "remove_from_state",
        },
        {
            name:           "throttled request",
            mockResponse:   &http.Response{StatusCode: 429},
            expectedAction: "retry",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test untyped error handling patterns
        })
    }
}
```

## Common Testing Patterns (Both Approaches)

### Resource ID Parsing Tests
```go
func TestParseServiceNameResourceID(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected *servicename.ServiceNameResourceId
        wantErr  bool
    }{
        {
            name:  "valid resource ID",
            input: "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Service/resources/resource1",
            expected: &servicename.ServiceNameResourceId{
                SubscriptionId:    "12345",
                ResourceGroupName: "rg1",
                ResourceName:      "resource1",
            },
            wantErr: false,
        },
        {
            name:    "invalid resource ID",
            input:   "invalid-id",
            wantErr: true,
        },
        {
            name:    "empty resource ID",
            input:   "",
            wantErr: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := servicename.ParseServiceNameResourceID(tc.input)
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
func TestValidateServiceNameResourceName(t *testing.T) {
    validNames := []string{
        "valid-name",
        "valid123",
        "name-with-numbers-123",
        "a", // minimum length if applicable
    }

    invalidNames := []string{
        "",                                    // empty
        "toolong" + strings.Repeat("a", 100), // too long
        "-invalid",                           // invalid start
        "invalid-",                           // invalid end
        "invalid_underscore",                 // invalid characters
        "Invalid-Uppercase",                  // invalid case
    }

    for _, name := range validNames {
        t.Run("valid_"+name, func(t *testing.T) {
            warnings, errors := servicename.ValidateServiceNameResourceName(name, "test")
            require.Empty(t, warnings)
            require.Empty(t, errors)
        })
    }

    for _, name := range invalidNames {
        t.Run("invalid_"+name, func(t *testing.T) {
            _, errors := servicename.ValidateServiceNameResourceName(name, "test")
            require.NotEmpty(t, errors)
        })
    }
}
```

### Azure-Specific Validation Tests
```go
func TestValidateAzureLocation(t *testing.T) {
    validLocations := []string{
        "East US",
        "West Europe",
        "Southeast Asia",
        "eastus",      // normalized form
        "westeurope",  // normalized form
    }

    invalidLocations := []string{
        "",
        "Invalid Location",
        "123-invalid",
    }

    for _, location := range validLocations {
        t.Run("valid_"+location, func(t *testing.T) {
            warnings, errors := servicename.ValidateAzureLocation(location, "location")
            require.Empty(t, warnings)
            require.Empty(t, errors)
        })
    }

    for _, location := range invalidLocations {
        t.Run("invalid_"+location, func(t *testing.T) {
            _, errors := servicename.ValidateAzureLocation(location, "location")
            require.NotEmpty(t, errors)
        })
    }
}
```

### Schema Expand/Flatten Function Tests
```go
func TestExpandServiceNameResourceConfig(t *testing.T) {
    input := []interface{}{
        map[string]interface{}{
            "name":    "test",
            "enabled": true,
            "settings": map[string]interface{}{
                "key1": "value1",
                "key2": "value2",
            },
            "tags": map[string]interface{}{
                "Environment": "test",
                "Project":     "terraform",
            },
        },
    }

    result := servicename.ExpandServiceNameResourceConfig(input)

    require.NotNil(t, result)
    require.Equal(t, "test", *result.Name)
    require.True(t, *result.Enabled)
    require.Equal(t, "value1", (*result.Settings)["key1"])
    require.Equal(t, "test", (*result.Tags)["Environment"])
}

func TestFlattenServiceNameResourceConfig(t *testing.T) {
    input := &servicename.ServiceNameResourceConfig{
        Name:    utils.String("test"),
        Enabled: utils.Bool(true),
        Settings: &map[string]string{
            "key1": "value1",
            "key2": "value2",
        },
        Tags: &map[string]string{
            "Environment": "test",
            "Project":     "terraform",
        },
    }

    result := servicename.FlattenServiceNameResourceConfig(input)

    require.Len(t, result, 1)
    config := result[0].(map[string]interface{})
    require.Equal(t, "test", config["name"])
    require.True(t, config["enabled"].(bool))

    settings := config["settings"].(map[string]interface{})
    require.Equal(t, "value1", settings["key1"])

    tags := config["tags"].(map[string]interface{})
    require.Equal(t, "test", tags["Environment"])
}
```

### Azure API Error Handling Tests
```go
func TestHandleAzureAPIErrors(t *testing.T) {
    testCases := []struct {
        name         string
        statusCode   int
        errorBody    string
        expectedType string
    }{
        {
            name:         "resource not found",
            statusCode:   404,
            errorBody:    {"error":{"code":"ResourceNotFound"}},
            expectedType: "not_found",
        },
        {
            name:         "resource conflict",
            statusCode:   409,
            errorBody:    {"error":{"code":"ResourceExists"}},
            expectedType: "conflict",
        },
        {
            name:         "throttled request",
            statusCode:   429,
            errorBody:    {"error":{"code":"TooManyRequests"}},
            expectedType: "throttled",
        },
        {
            name:         "server error",
            statusCode:   500,
            errorBody:    {"error":{"code":"InternalServerError"}},
            expectedType: "server_error",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test Azure API error handling patterns
            resp := &http.Response{
                StatusCode: tc.statusCode,
                Body:       ioutil.NopCloser(strings.NewReader(tc.errorBody)),
            }

            errorType := servicename.ClassifyAzureAPIError(resp)
            require.Equal(t, tc.expectedType, errorType)
        })
    }
}
```

## Azure Service-Specific Testing Considerations

### Resource Naming and Validation
- Test Azure resource naming conventions
- Test globally unique name requirements (where applicable)
- Test region-specific naming rules
- Test character restrictions and length limits

### Location and Region Handling
```go
func TestAzureLocationNormalization(t *testing.T) {
    testCases := []struct {
        input    string
        expected string
    }{
        {"East US", "eastus"},
        {"West Europe", "westeurope"},
        {"Southeast Asia", "southeastasia"},
        {"eastus", "eastus"}, // already normalized
    }

    for _, tc := range testCases {
        t.Run(tc.input, func(t *testing.T) {
            result := servicename.NormalizeAzureLocation(tc.input)
            require.Equal(t, tc.expected, result)
        })
    }
}
```

### Resource Tags Testing
```go
func TestAzureResourceTags(t *testing.T) {
    input := map[string]interface{}{
        "Environment": "production",
        "Project":     "terraform-provider",
        "Owner":       "platform-team",
    }

    result := servicename.ExpandTags(input)

    require.NotNil(t, result)
    require.Equal(t, "production", (*result)["Environment"])
    require.Equal(t, "terraform-provider", (*result)["Project"])
    require.Equal(t, "platform-team", (*result)["Owner"])
}
```

## Running Tests

- Run unit tests: go test ./internal/services/servicename/...
- Run specific test: go test -run TestFunctionName ./internal/services/servicename/
- Run tests with coverage: go test -cover ./internal/services/servicename/...
- Run acceptance tests: make testacc TEST=./internal/services/servicename TESTARGS='-run=TestAccResourceName'
- Run tests for typed resource resources: Focus on receiver method testing and model validation
- Run tests for untyped resources: Focus on function-based patterns and traditional state management

## Best Practices Summary

1. **Implementation Awareness**: Choose appropriate testing patterns based on whether the code uses typed resource or untyped approaches
2. **Package Naming**: Use external test package convention (package servicename_test)
3. **Test Coverage**: Focus on utility functions, parsers, and validators rather than full resource lifecycle
4. **Azure Integration**: Test Azure-specific patterns like resource ID parsing, location handling, and API error responses
5. **Error Handling**: Test both typed resource error patterns (metadata methods) and untyped patterns (traditional error handling)
6. **Consistency**: Ensure tests follow the same patterns regardless of implementation approach for common functionality
7. **Mocking**: Use appropriate mocking for Azure API calls and external dependencies
8. **Edge Cases**: Test boundary conditions, invalid inputs, and error scenarios thoroughly
