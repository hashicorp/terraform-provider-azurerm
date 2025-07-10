---
applyTo: "internal/**/*.go"
description: This document outlines the testing guidelines for Go files in the Terraform AzureRM provider repository. It includes unit test organization, acceptance test patterns, naming conventions, and best practices for Azure API integration testing.
---

## Testing Guidelines
Given below are the testing guidelines for the Terraform AzureRM provider which **MUST** be followed.

### Implementation Approach Testing Considerations

#### Typed vs UnTyped Resource Testing

While the core testing principles remain the same, there are some differences in testing patterns between Typed Resource and UnTyped Resource (Plugin SDK) resource implementations:

**Typed Resource Testing:**
- Resources use the `internal/sdk` framework with type-safe models
- Testing patterns leverage `sdk.ResourceFunc` return types with `metadata` access
- Resource existence checks access clients through `metadata.Client`
- Error handling and logging use structured `metadata` patterns
- State management uses `metadata.Decode()` and `metadata.Encode()`

**UnTyped Resource Testing:**
- Resources use traditional Plugin SDK patterns with function-based CRUD
- Testing patterns use direct `*pluginsdk.ResourceData` and `clients.Client` access
- Resource existence checks use traditional client initialization patterns
- Error handling uses traditional error patterns and direct state manipulation
- State management uses `d.Set()` and `d.Get()` patterns

**Testing Consistency Requirements:**
- **User Experience**: Acceptance tests should be identical regardless of implementation approach
- **Test Framework**: Both approaches use the same acceptance testing framework
- **Resource Lifecycle**: Both approaches test the same CRUD operations and import functionality
- **Azure Integration**: Both approaches test the same Azure API interactions and behaviors

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
- Place acceptance tests in files ending with `_test.go` in the same directory as the resource
- **Package Naming**: Acceptance tests must use the resource package name with `_test` appended
  - Example: If the resource is in package `cdn`, tests should be in package `cdn_test`
  - Example: If the resource is in package `compute`, tests should be in package `compute_test`
  - This follows Go's external test package convention for testing exported APIs
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

### Testing CustomizeDiff Validations: Critical Requirements

CustomizeDiff validations are essential for enforcing Azure API constraints and preventing invalid configurations. **Testing these validations is mandatory** and requires comprehensive coverage of both success and failure scenarios.

#### Why CustomizeDiff Testing is Critical

**Azure API Constraint Enforcement:**
- CustomizeDiff validations prevent invalid API calls that would fail at runtime
- They enforce Azure service-specific field combination requirements
- They validate complex resource dependencies before Azure API interaction
- They provide clear error messages to users before resource creation/update

**Testing Requirements:**
- **Error Scenarios**: Test all invalid field combinations that should trigger validation errors
- **Success Scenarios**: Test all valid field combinations that should pass validation
- **Edge Cases**: Test boundary conditions and corner cases
- **Error Message Validation**: Verify specific error messages using `ExpectError: regexp.MustCompile()`
- **Field Path Accuracy**: Ensure error messages include correct field paths and constraints
- **Azure API Alignment**: Test that validations match actual Azure API behavior

#### CustomizeDiff Testing Mandatory Practices

**Comprehensive Test Coverage:**
```go
func TestAccServiceName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        // REQUIRED: Test invalid configuration
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is true"),
        },
        // REQUIRED: Test valid configuration  
        {
            Config: r.validConfiguration(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        // REQUIRED: Test import for valid configuration
        data.ImportStep(),
    })
}
```

**Azure-Specific Validation Testing:**
- Test Azure service-specific constraints (SKU dependencies, region limitations, etc.)
- Validate Azure API field combination requirements
- Test Azure resource lifecycle constraints
- Verify Azure service version-specific validations

### CustomizeDiff Testing Patterns

When testing resources that use CustomizeDiff, remember the dual import requirement:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
```

**Testing CustomizeDiff Validation:**
```go
func TestAccServiceName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is true"),
        },
        {
            Config: r.validConfiguration(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
    })
}
```

**Testing Azure-Specific CustomizeDiff Validation (CDN Front Door Log Scrubbing):**
```go
func TestAccCdnFrontDoorProfile_logScrubbingValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            // Test that selector cannot be set for RequestIPAddress
            Config:      r.logScrubbingInvalidRequestIPAddress(data),
            ExpectError: regexp.MustCompile(`log_scrubbing\.0\.scrubbing_rule\.0: ` + "`selector`" + ` cannot be set when ` + "`match_variable`" + ` is ` + "`RequestIPAddress`"),
        },
        {
            // Test that selector cannot be set for RequestUri  
            Config:      r.logScrubbingInvalidRequestUri(data),
            ExpectError: regexp.MustCompile(`log_scrubbing\.0\.scrubbing_rule\.0: ` + "`selector`" + ` cannot be set when ` + "`match_variable`" + ` is ` + "`RequestUri`"),
        },
        {
            // Test valid configuration with QueryStringArgNames and selector
            Config: r.logScrubbingValidQueryStringWithSelector(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rules.0.match_variable").HasValue("QueryStringArgNames"),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rules.0.selector").HasValue("custom_param"),
            ),
        },
        {
            // Test invalid configuration with QueryStringArgNames without selector
            Config:      r.logScrubbingInvalidQueryStringWithoutSelector(data),
            ExpectError: regexp.MustCompile(`log_scrubbing\.0\.scrubbing_rules\.0: ` + "`selector`" + ` is required when ` + "`match_variable`" + ` is ` + "`QueryStringArgNames`"),
        },
        {
            // Test valid configuration with RequestIPAddress without selector
            Config: r.logScrubbingValidRequestIPAddress(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rules.0.match_variable").HasValue("RequestIPAddress"),
            ),
        },
        {
            // Test valid configuration with RequestUri without selector
            Config: r.logScrubbingValidRequestUri(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rules.0.match_variable").HasValue("RequestUri"),
            ),
        },
    })
}
```

**Testing ForceNew Behavior with CustomizeDiff:**
```go
func TestAccServiceName_forceNewOnPropertyChange(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("property_name").HasValue("initial_value"),
            ),
        },
        {
            Config: r.updatedProperty(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("property_name").HasValue("updated_value"),
                // Verify the resource ID changed (ForceNew behavior)
                check.That(data.ResourceName).Key("id").Exists(),
            ),
        },
    })
}
```

**CustomizeDiff Test Configuration Helper Functions:**
```go
// Test configurations for CDN Front Door log scrubbing validation
func (r CdnFrontDoorProfileResource) logScrubbingInvalidRequestIPAddress(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      selector       = "invalid_selector"  # This should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingInvalidRequestUri(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      selector       = "invalid_selector"  # This should trigger validation error
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingValidQueryStringWithSelector(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      operator       = "EqualsAny"
      selector       = "custom_param"  # This is valid for QueryStringArgNames
      enabled        = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingInvalidQueryStringWithoutSelector(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "QueryStringArgNames"
      operator       = "EqualsAny"
      enabled        = true
      # No selector specified - this should trigger validation error
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingValidRequestIPAddress(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"
      enabled        = true
      # No selector specified - this is required for RequestIPAddress
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontDoorProfileResource) logScrubbingValidRequestUri(data acceptance.TestData) string {
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

  log_scrubbing {
    enabled = true
    scrubbing_rule {
      match_variable = "RequestUri"
      operator       = "EqualsAny"
      enabled        = true
      # No selector specified - this is required for RequestUri
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
```

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

### Test Organization and Placement Rules

#### Acceptance Test File Structure
- **Test function placement**: All test functions must be placed before the `Exists` function in the test file
- **Helper function placement**: Test configuration helper functions should be placed after the `Exists` function
- **No duplicate functions**: Remove any duplicate or old test functions to maintain clean file structure
- **Consistent ordering**: Place tests in logical order (basic, update, requires import, other scenarios)

Example of proper test file structure:
```go
package cdn_test

// Test functions - BEFORE Exists function
func TestAccCdnFrontDoorProfile_basic(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_update(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_requiresImport(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_logScrubbing(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_logScrubbingUpdate(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_logScrubbingDisabled(t *testing.T) { ... }

// Exists function - SEPARATOR between tests and helpers
func (CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) { ... }

// Helper functions - AFTER Exists function
func (r CdnFrontDoorProfileResource) basic(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) requiresImport(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) logScrubbing(data acceptance.TestData) string { ... }
```

#### Test Configuration Standards
- **Azure value validation**: Only use valid Azure service values in test configurations
- **SDK constant alignment**: Test values should match Azure SDK constants and API documentation
- **Cross-resource consistency**: When testing similar features across resources, use consistent value patterns

### Azure-Specific Testing Guidelines

#### Resource Existence Checks

The implementation of resource existence checks differs between typed and untyped approaches:

**Typed Resource Existence Check:**
```go
func (r ServiceNameResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := parse.ServiceNameID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.ServiceName.ResourceClient.Get(ctx, *id)
    if err != nil {
        return nil, fmt.Errorf("reading %s: %+v", *id, err)
    }

    return utils.Bool(resp.Model != nil), nil
}
```

**UnTyped Resource Resource Existence Check:**
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

Test configuration templates should be consistent regardless of implementation approach, but the underlying resource structure may differ:

**Typed Resource Test Configuration Example:**
```go
func (r ServiceNameResource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_name" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name           = "Standard"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
```

**UnTyped Resource Test Configuration Example:**
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

### CustomizeDiff Testing Anti-Patterns and Best Practices

#### Common Testing Mistakes to Avoid

**Insufficient Error Testing:**
```go
// AVOID: Only testing valid configurations
func TestAccServiceName_customizeDiff(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.validConfiguration(data), // Only testing success case
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
    })
}

// GOOD: Testing both valid and invalid configurations
func TestAccServiceName_customizeDiff(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("specific validation error"),
        },
        {
            Config: r.validConfiguration(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
    })
}
```

**Vague Error Message Testing:**
```go
// AVOID: Not validating specific error messages
{
    Config:      r.invalidConfiguration(data),
    ExpectError: regexp.MustCompile("error"), // Too vague
}

// GOOD: Validate specific error messages with field paths
{
    Config:      r.invalidLogScrubbing(data),
    ExpectError: regexp.MustCompile(`log_scrubbing\.0\.scrubbing_rule\.0: ` + "`selector`" + ` cannot be set when ` + "`match_variable`" + ` is ` + "`RequestIPAddress`"),
}
```

**Missing Edge Case Testing:**
```go
// AVOID: Only testing obvious invalid cases
func TestAccServiceName_validation(t *testing.T) {
    // Only test one invalid scenario
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.obviouslyInvalid(data),
            ExpectError: regexp.MustCompile("validation error"),
        },
    })
}

// GOOD: Test all boundary conditions and edge cases
func TestAccServiceName_validation(t *testing.T) {
    data.ResourceTest(t, r, []acceptance.TestStep{
        // Test each invalid field combination
        {
            Config:      r.invalidCase1(data),
            ExpectError: regexp.MustCompile("specific error 1"),
        },
        {
            Config:      r.invalidCase2(data),
            ExpectError: regexp.MustCompile("specific error 2"),
        },
        {
            Config:      r.edgeCase(data),
            ExpectError: regexp.MustCompile("edge case error"),
        },
        // Test valid configurations
        {
            Config: r.validCase1(data),
            Check:  acceptance.ComposeTestCheckFunc(...),
        },
        {
            Config: r.validCase2(data),
            Check:  acceptance.ComposeTestCheckFunc(...),
        },
    })
}
```

#### CustomizeDiff Testing Best Practices

**Comprehensive Scenario Coverage:**
- Test all possible field combinations that trigger validation
- Test boundary values and edge cases
- Test Azure service-specific constraints
- Verify error messages include field paths and clear guidance
- Test both creation and update scenarios with CustomizeDiff

**Azure API Alignment:**
- Ensure test values match Azure SDK constants
- Test configurations that would be rejected by Azure API
- Validate that CustomizeDiff catches errors before Azure API calls
- Test combinations that are Azure service version-specific

**Error Message Quality:**
- Verify error messages include specific field names in backticks
- Check that error messages explain the constraint clearly
- Ensure field paths are accurate (e.g., `log_scrubbing.0.scrubbing_rule.0`)
- Test that error messages provide actionable guidance

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


### CustomizeDiff Testing Patterns

When testing resources that use CustomizeDiff, remember the dual import requirement:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
```

**Testing CustomizeDiff Validation:**
```go
func TestAccServiceName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is true"),
        },
        {
            Config: r.validConfiguration(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
    })
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

```powershell
$env:ARM_SUBSCRIPTION_ID="your-azure-subscription-id"
$env:ARM_CLIENT_ID="your-service-principal-client-id"
$env:ARM_CLIENT_SECRET="your-service-principal-client-secret"
$env:ARM_TENANT_ID="your-azure-tenant-id"
$env:ARM_TEST_LOCATION="WestEurope"
$env:ARM_TEST_LOCATION_ALT="EastUS2"
$env:ARM_TEST_LOCATION_ALT2="WestUS2"
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

#### CDN Front Door Testing Considerations
- **Global Resource**: CDN Front Door profiles are global resources, not tied to specific Azure regions
- **SKU Limitations**: Different features are available in Standard vs Premium SKUs
- **Naming Constraints**: Profile names must be globally unique across Azure
- **Propagation Time**: Changes may take time to propagate globally
- **Response Timeout**: Valid range is 16-240 seconds
