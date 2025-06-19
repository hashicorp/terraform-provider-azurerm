---
applyTo: "internal/**/*.go"
description: This document outlines the testing guidelines for Go files in the Terraform AzureRM provider repository. It includes unit test organization, acceptance test patterns, naming conventions, and best practices for Azure API integration testing.
---

## Testing Guidelines
Given below are the testing guidelines for the Terraform AzureRM provider which **MUST** be followed.

### Test Types

#### Unit Tests
- Place unit tests in the same package as the source code with `_test.go` suffix
- Focus on utility functions, parsers, validators, and isolated business logic
- Use Go's built-in `testing` package and `testify` for assertions
- Follow table-driven test patterns for multiple scenarios
- Mock external dependencies using interfaces
- Test edge cases, error scenarios, and Azure-specific validation logic
- Unit tests should not require Azure credentials or make external API calls

#### Acceptance Tests
- Primary testing method for Terraform resource lifecycle (CRUD operations)
- Place acceptance tests in files ending with `_test.go` in the same package
- Test against real Azure APIs with live Azure credentials
- Use the acceptance testing framework provided by the provider
- Test resource creation, updates, imports, and deletion
- Test dependent resource relationships and state management

### Naming Conventions

#### Unit Test Functions
- Use `TestFunctionName_Scenario_ExpectedOutcome` format
- Example: `TestParseFrontDoorProfileID_ValidID_ReturnsCorrectComponents`
- Example: `TestValidateFrontDoorProfileName_TooLong_ReturnsError`

#### Acceptance Test Functions
- Use `TestAccResourceName_scenario` format
- Example: `TestAccCdnFrontDoorProfile_basic`
- Example: `TestAccCdnFrontDoorProfile_update`
- Example: `TestAccCdnFrontDoorProfile_requiresImport`

### Go Testing Patterns

#### Table-Driven Tests
```go
func TestParseResourceID(t *testing.T) {
    testCases := []struct {
        name        string
        input       string
        expected    *ResourceId
        expectError bool
    }{        {
            name:  "valid resource ID",
            input: "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Cdn/profiles/profile1",
            expected: &ResourceId{
                SubscriptionId: "12345",
                ResourceGroup:  "rg1", 
                Name:          "profile1",
            },
            expectError: false,
        },
        {
            name:        "invalid resource ID",
            input:       "invalid-id",
            expected:    nil,
            expectError: true,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := ParseResourceID(tc.input)
            if tc.expectError {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tc.expected, result)
        })
    }
}
```

#### Assertion Patterns
- Use `require.NoError(t, err)` for errors that should stop test execution
- Use `assert.Error(t, err)` for expected errors that shouldn't stop execution
- Use `require.Equal(t, expected, actual)` for value comparisons
- Use `require.NotNil(t, result)` or `require.Nil(t, result)` for nil checks
- Use `require.True(t, condition)` or `require.False(t, condition)` for boolean conditions

### Acceptance Testing Patterns

#### Basic Resource Test
```go
func TestAccCdnFrontDoorProfile_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-%d", data.RandomInteger)),
                check.That(data.ResourceName).Key("sku_name").HasValue("Standard_AzureFrontDoor"),
            ),
        },
        data.ImportStep(), // No sensitive fields to exclude for CDN profiles
    })
}
```

#### Resource Update Test
```go
func TestAccCdnFrontDoorProfile_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.updated(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("tags.Environment").HasValue("Test"),
            ),
        },
        data.ImportStep(),
    })
}
```

#### Resource Requires Import Test
```go
func TestAccCdnFrontDoorProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
	r := CdnFrontDoorProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}
```

### Azure-Specific Testing Guidelines

#### Resource Existence Checks
```go
func (CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := parse.FrontDoorProfileID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.Cdn.FrontDoorProfilesClient.Get(ctx, *id)
    if err != nil {
        return nil, fmt.Errorf("reading CDN Front Door Profile (%s): %+v", *id, err)
    }

    return utils.Bool(resp.Model != nil), nil
}
```

#### Test Configuration Templates
```go
func (CdnFrontDoorProfileResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = "Standard_AzureFrontDoor"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) complete(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = "Premium_AzureFrontDoor"

  response_timeout_seconds = 120

  tags = {
    environment = "Production"
    project     = "AccTest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (CdnFrontDoorProfileResource) updated(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfd-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = "Standard_AzureFrontDoor"

  response_timeout_seconds = 240

  tags = {
    environment = "Test"
    project     = "AccTest"
    updated     = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
```

### Environment Setup

#### Required Environment Variables
```bash
export ARM_SUBSCRIPTION_ID="your-azure-subscription-id"
export ARM_CLIENT_ID="your-service-principal-client-id"
export ARM_CLIENT_SECRET="your-service-principal-client-secret"
export ARM_TENANT_ID="your-azure-tenant-id"
export ARM_TEST_LOCATION=WestEurope
export ARM_TEST_LOCATION_ALT=EastUS2
export ARM_TEST_LOCATION_ALT2=WestUS2
```

#### Running Tests
```bash
# Run unit tests
go test ./internal/services/cdn/...

# Run specific unit test
go test -run TestParseFrontDoorProfileID ./internal/services/cdn/

# Run acceptance tests for a service
make testacc TEST=./internal/services/cdn TESTARGS='-run=TestAccCdnFrontDoorProfile'

# Run specific acceptance test
make testacc TEST=./internal/services/cdn TESTARGS='-run=TestAccCdnFrontDoorProfile_basic'
```

### Best Practices

#### Test Organization
- Group related tests using subtests with `t.Run()`
- Use `data.ResourceTest()` for acceptance tests to ensure proper cleanup
- Always include import tests for resources (`data.ImportStep()`)
- Test both successful operations and error conditions

#### Azure Resource Management
- Use `acceptance.BuildTestData()` for consistent test data generation
- Include dependency resources (resource groups, networks) in test configurations
- Use unique naming with random integers to avoid resource conflicts
- Test in multiple Azure regions when relevant

#### Error Handling
- Test validation functions with both valid and invalid inputs
- Verify specific error messages and error types
- Test Azure API error scenarios (resource not found, throttling, etc.)
- Ensure proper error wrapping and context information

#### Performance Considerations
- Use parallel test execution where possible (`t.Parallel()`)
- Minimize Azure API calls in unit tests
- Use appropriate timeouts for long-running operations
- Clean up test resources properly to avoid quota issues

### CDN Front Door-Specific Testing Guidelines

#### SKU Validation Testing
```go
func TestValidateFrontDoorProfileSku(t *testing.T) {
    validSkus := []string{
        "Standard_AzureFrontDoor",
        "Premium_AzureFrontDoor",
    }
    invalidSkus := []string{
        "Basic",
        "Standard",
        "Premium",
        "Invalid_SKU",
    }
    
    for _, sku := range validSkus {
        t.Run("valid_"+sku, func(t *testing.T) {
            warnings, errors := ValidateFrontDoorProfileSku(sku, "test")
            require.Empty(t, warnings)
            require.Empty(t, errors)
        })
    }
    
    for _, sku := range invalidSkus {
        t.Run("invalid_"+sku, func(t *testing.T) {
            _, errors := ValidateFrontDoorProfileSku(sku, "test")
            require.NotEmpty(t, errors)
        })
    }
}
```

#### CDN Front Door Testing Considerations
- **Global Resource**: CDN Front Door profiles are global resources, not tied to specific Azure regions
- **SKU Limitations**: Different features are available in Standard vs Premium SKUs
- **Naming Constraints**: Profile names must be globally unique across Azure
- **Propagation Time**: Changes may take time to propagate globally
- **Response Timeout**: Valid range is 16-240 seconds
