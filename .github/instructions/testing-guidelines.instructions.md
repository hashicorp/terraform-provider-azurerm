---
applyTo: "internal/**/*.go"
description: This document outlines the testing guidelines for Go files in the Terraform AzureRM provider repository. It includes unit test organization, acceptance test patterns, naming conventions, and best practices for Azure API integration testing.
---

## Testing Guidelines
Given below are the testing guidelines for the Terraform AzureRM provider which **MUST** be followed.

### 🚨 **CRITICAL: Test Execution Protocol** 🚨

**⚠️ DO NOT RUN TESTS AUTOMATICALLY ⚠️**

**MANDATORY TEST EXECUTION PROTOCOL:**

#### **Rule #1: NEVER RUN TESTS WITHOUT EXPLICIT USER PERMISSION**
- **DO NOT** execute `make testacc`, `go test`, or any test commands automatically
- **DO NOT** run tests "to verify the implementation"
- **DO NOT** run tests "to check if everything works"
- **ALWAYS** provide the exact commands for the user to run manually
- **ALWAYS** explain what the tests will do before suggesting commands

#### **Rule #2: Test Command Provision Protocol**
When providing test commands to users:
1. **Explain the purpose**: What will this test verify?
2. **Provide exact commands**: Full command with proper arguments
3. **Set expectations**: How long will it take? What resources will be created?
4. **Include cleanup verification**: How to confirm resources are cleaned up

#### **Rule #3: Azure Resource Testing Considerations**
- Tests create **REAL AZURE RESOURCES** that incur costs
- Tests require **VALID AZURE CREDENTIALS** to be configured
- Some tests take **10+ MINUTES** to complete due to Azure provisioning times
- Failed tests may leave **ORPHANED RESOURCES** that continue to incur costs

#### **Example of Correct Test Command Provision:**
```powershell
# This will test VMSS resiliency policy backward compatibility
# Expected duration: 5-10 minutes, creates test VMSS resources in Azure
# Requires: ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID

make testacc TEST=./internal/services/compute TESTARGS='-run=TestAccLinuxVirtualMachineScaleSet_fieldsNotSetInState'
```

#### **Rule #4: User Confirmation Required**
- Before providing test commands, confirm user wants to run tests
- Verify user has Azure credentials configured
- Confirm user understands costs and time requirements
- Provide cleanup verification steps

**This protocol prevents accidental resource creation, unexpected costs, and ensures users maintain control over their Azure environment.**

---

### Implementation Approach Testing Considerations

#### Testing Approach Consistency

While the core testing principles remain the same, there are some differences in testing patterns between Typed Resource and UnTyped Resource implementations:

**For detailed implementation approach information, see the main copilot instructions file.**

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
// When testing resources with CustomizeDiff, remember the dual import requirement:
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
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

**Testing Azure-Specific CustomizeDiff Validation (CDN Front Door Firewall Policy Log Scrubbing):**
```go
func TestAccCdnFrontDoorFirewallPolicy_logScrubbingValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_firewall_policy", "test")
    r := CdnFrontDoorFirewallPolicyResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            // Test that selector cannot be set when using EqualsAny operator
            Config:      r.logScrubbingInvalidSelectorWithEqualsAny(data),
            ExpectError: regexp.MustCompile(`the 'selector' field cannot be set when the "EqualsAny" 'operator' is used`),
        },
        {
            // Test that selector is required when using Equals operator
            Config:      r.logScrubbingInvalidMissingSelectorWithEquals(data),
            ExpectError: regexp.MustCompile(`the 'selector' field must be set when the "Equals" 'operator' is used`),
        },
        {
            // Test that RequestIPAddress must use EqualsAny operator
            Config:      r.logScrubbingInvalidRequestIPAddressWithEquals(data),
            ExpectError: regexp.MustCompile(`the "RequestIPAddress" 'match_variable' must use the "EqualsAny" 'operator'`),
        },
        {
            // Test that RequestUri must use EqualsAny operator
            Config:      r.logScrubbingInvalidRequestUriWithEquals(data),
            ExpectError: regexp.MustCompile(`the "RequestUri" 'match_variable' must use the "EqualsAny" 'operator'`),
        },
        {
            // Test valid configuration with QueryStringArgNames and Equals operator with selector
            Config: r.logScrubbingValidQueryStringWithEquals(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("QueryStringArgNames"),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.operator").HasValue("Equals"),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.selector").HasValue("custom_param"),
            ),
        },
        {
            // Test valid configuration with RequestIPAddress and EqualsAny operator
            Config: r.logScrubbingValidRequestIPAddress(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("RequestIPAddress"),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.operator").HasValue("EqualsAny"),
            ),
        },
        {
            // Test valid configuration with RequestUri and EqualsAny operator
            Config: r.logScrubbingValidRequestUri(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.match_variable").HasValue("RequestUri"),
                check.That(data.ResourceName).Key("log_scrubbing.0.scrubbing_rule.0.operator").HasValue("EqualsAny"),
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
// Test configurations for CDN Front Door Firewall Policy log scrubbing validation
func (r CdnFrontDoorFirewallPolicyResource) logScrubbingInvalidSelectorWithEqualsAny(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestBodyJsonArgNames"
      operator       = "EqualsAny"
      selector       = "invalid_selector"  # This should trigger validation error
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingInvalidMissingSelectorWithEquals(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      # No selector specified - this should trigger validation error
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingInvalidRequestIPAddressWithEquals(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestIPAddress"
      operator       = "Equals"  # This should trigger validation error - must be EqualsAny
      selector       = "invalid_selector"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingInvalidRequestUriWithEquals(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestUri"
      operator       = "Equals"  # This should trigger validation error - must be EqualsAny
      selector       = "invalid_selector"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingValidQueryStringWithEquals(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "QueryStringArgNames"
      operator       = "Equals"
      selector       = "custom_param"  # This is valid for QueryStringArgNames with Equals
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingValidRequestIPAddress(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestIPAddress"
      operator       = "EqualsAny"  # This is required for RequestIPAddress
      # No selector - this is required for EqualsAny operator
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorFirewallPolicyResource) logScrubbingValidRequestUri(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_firewall_policy" "test" {
  name                = "accTestWAF%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name           = azurerm_cdn_frontdoor_profile.test.sku_name
  mode               = "Prevention"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      enabled        = true
      match_variable = "RequestUri"
      operator       = "EqualsAny"  # This is required for RequestUri
      # No selector - this is required for EqualsAny operator
    }
  }
}
`, r.template(data), data.RandomInteger)
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
                check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-cdn-%d", data.RandomInteger)),
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

### ImportStep Validation Guidelines

#### Redundant Validation Checks with ImportStep

When using `data.ImportStep()` in acceptance tests, most field validation checks are **redundant** because ImportStep automatically validates that the resource can be imported and that all field values match between the configuration and the imported state.

**Standard Pattern - Minimal Checks with ImportStep:**
```go
func TestAccCdnFrontDoorProfile_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r), // Essential - verifies resource exists
                // AVOID: check.That(data.ResourceName).Key("name").HasValue(...) - ImportStep validates this
                // AVOID: check.That(data.ResourceName).Key("sku_name").HasValue(...) - ImportStep validates this
            ),
        },
        data.ImportStep(), // Automatically validates all configured field values
    })
}
```

**Legitimate Validation Checks with ImportStep:**
In rare circumstances, additional validation checks may be appropriate alongside ImportStep:

1. **Computed Field Verification**: Testing computed values that aren't in the configuration
2. **Complex Behavior Validation**: Testing TypeSet deduplication, ordering, or transformation logic
3. **Azure API-Specific Behavior**: Testing Azure service-specific transformations

**Example - Valid Additional Checks:**
```go
func TestAccCdnFrontDoorProfile_logScrubbing_withDuplicateScrubbingRules(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.logScrubbingDuplicateScrubbingRules(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // VALID: Testing TypeSet deduplication behavior - not simple field validation
                check.That(data.ResourceName).Key("log_scrubbing_rule.#").HasValue("2"),
            ),
        },
        data.ImportStep(), // Still validates all other field values
    })
}
```

**Key Principles:**
- **ImportStep handles field validation**: Don't duplicate validation of configured field values
- **Keep only ExistsInAzure**: Essential for verifying resource creation/existence
- **Add checks sparingly**: Only for behavior that ImportStep cannot verify
- **Document rationale**: Comment why additional checks are needed when used

### Test Organization and Placement Rules

#### Acceptance Test File Structure
- **Test function placement**: All test functions must be placed before the `Exists` function in the test file. This ensures a consistent structure across test files, making it easier to locate and understand test cases. Additionally, placing test functions first improves readability by prioritizing the main test logic over helper functions.
- **Helper function placement**: Test configuration helper functions should be placed after the `Exists` function to separate them from the main test logic and maintain a clean structure.
- **No duplicate functions**: Remove any duplicate or old test functions to maintain clean file structure.
- **Consistent ordering**: Place tests in logical order (basic, update, requires import, other scenarios).
- **Exceptions for complex scenarios**: In cases where helper functions or dependencies are required earlier in the file for complex test scenarios, it is acceptable to deviate from this rule. Developers should document the rationale for such deviations within the code to ensure clarity for future maintainers.

#### Test Case Consolidation Standards

**HashiCorp Standard - Essential Tests Only:**
- **Basic Test**: Core functionality with minimal configuration
- **Update Test**: Resource update scenarios (only if resource supports updates)
- **RequiresImport Test**: Import conflict detection
- **Complete Test**: Full feature demonstration (optional, for complex resources)

**AVOID Excessive Test Cases:**
- Multiple basic tests with minor variations
- Separate tests for each individual field
- Redundant validation tests that don't add value
- Over-testing obvious functionality

**Example of Proper Test Organization:**
```go
package cdn_test

// ESSENTIAL TESTS - Keep these
func TestAccCdnFrontDoorProfile_basic(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_update(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_requiresImport(t *testing.T) { ... }

// FEATURE-SPECIFIC TESTS - Only if testing complex features
func TestAccCdnFrontDoorProfile_logScrubbing(t *testing.T) { ... }

// Exists function - SEPARATOR between tests and helpers
func (CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) { ... }

// Helper functions - AFTER Exists function
func (r CdnFrontDoorProfileResource) basic(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) requiresImport(data acceptance.TestData) string { ... }
```

#### Test Configuration Consolidation

**Consolidate Multiple Examples into Single, Comprehensive Examples:**

**PREFERRED - Single Comprehensive Example:**
```go
func (r CdnFrontDoorProfileResource) complete(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                     = "acctestcdnfd-%d"
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                = "Premium_AzureFrontDoor"
  response_timeout_seconds = 120

  scrubbing_rule {
    match_variable = "QueryStringArgNames"
    operator       = "Equals"
    selector       = "secret"
  }

  scrubbing_rule {
    match_variable = "RequestIPAddress"
    operator       = "EqualsAny"
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
```

**AVOID - Multiple Separate Examples:**
```go
// Don't create separate functions for minor variations
func (r CdnFrontDoorProfileResource) withTimeout(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) withTags(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) withScrubbing(data acceptance.TestData) string { ... }
```

#### Test Configuration Standards
- **Azure value validation**: Only use valid Azure service values in test configurations
- **SDK constant alignment**: Test values should match Azure SDK constants and API documentation
- **Cross-resource consistency**: When testing similar features across resources, use consistent value patterns

#### Field Rename Testing Requirements

When implementing field renames for better descriptive naming, comprehensive testing across all affected files is mandatory:

**Field Rename Testing Checklist:**
- **Resource implementation**: Update field name in schema definition and all references
- **Data source implementation**: Ensure consistency between resource and data source field names
- **Test configurations**: Update all test helper functions to use the new field name
- **Documentation**: Update website documentation with new field name and examples
- **Import functionality**: Verify that resource import still works correctly after field rename
- **State compatibility**: Ensure that existing Terraform state remains compatible

**Example Field Rename Pattern:**
```go
// BEFORE - Generic field name
"scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "match_variable": {
                Type:     pluginsdk.TypeString,
                Required: true,
            },
        },
    },
},

// AFTER - Descriptive field name
"log_scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "match_variable": {
                Type:     pluginsdk.TypeString,
                Required: true,
            },
        },
    },
},
```

**Field Rename Validation Testing:**
```go
func TestAccCdnFrontDoorProfile_logScrubbingRuleFieldRename(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.withLogScrubbingRule(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(), // Verify import works with new field name
    })
}
```

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
  name     = "acctestRG-servicename-%d"
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
  name     = "acctestRG-cdn-%d"
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
  name     = "acctestRG-cdn-%d"
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
  name     = "acctestRG-cdn-%d"
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

### SKU Validation Testing

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

### Azure Provider Feature Flags for Test Cleanup

**Critical Insight**: Some Azure resources with protective features (like VMSS resiliency, database high availability, etc.) may block Terraform's default cleanup behavior during test teardown. The test framework automatically tries to scale down resources to 0 instances before deletion, but Azure may prevent this due to service-specific requirements.

**Root Cause Understanding**:
When tests fail with errors like `OperationNotAllowed: Cannot update [resource] when [protective feature] is enabled`, the issue is typically the **auto scale-down behavior** during cleanup. The Terraform test framework tries to gracefully scale resources to 0 instances before deletion, but Azure blocks this operation to maintain service guarantees (health monitoring, backup retention, etc.).

**Force Delete Pattern for Protected Resources:**
```go
// Template function with force delete configuration
func (r ServiceNameResource) templateWithForceDelete(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    virtual_machine_scale_set {
      force_delete = true
    }
    # Add other service-specific force delete flags as needed
    # key_vault {
    #   purge_soft_delete_on_destroy = true
    # }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-service-%d"
  location = "%s"
}
# ... other base resources
`, r.templatePublicKey(), data.RandomInteger, data.Locations.Primary)
}

// Usage in test configurations
func (r ServiceNameResource) resiliencyTestWithCleanup(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_service_resource" "test" {
  # ... resource configuration with protective features
  resiliency {
    enabled = true
  }
}
`, r.templateWithForceDelete(data))
}
```

**When to Use Force Delete Flags:**
- **VMSS with resiliency enabled**: `virtual_machine_scale_set.force_delete = true`
- **Key Vault with soft delete**: `key_vault.purge_soft_delete_on_destroy = true`
- **SQL databases with backup protection**: Appropriate force delete flags
- **Any resource where Azure blocks normal Terraform cleanup operations**

**Common Cleanup Error Patterns to Watch For:**
- `OperationNotAllowed: Cannot update [resource] when [protective feature] is enabled`
- `ResourceGroupBeingDeleted: Cannot perform operation while resource group is being deleted`
- Scale-down operations blocked due to health monitoring requirements
- Soft-delete conflicts preventing immediate recreation

**Template Organization Best Practices:**
- Create semantic template functions: `templateWithForceDelete`, `templateWithBackupRetention`, etc.
- Encapsulate provider feature flags within template functions rather than inline
- Use descriptive names that clearly indicate the template's cleanup behavior
- Follow existing patterns like `templateWithLocation` for consistency

**Debugging Test Cleanup Issues:**
1. **Identify the blocking operation**: Look for Azure API errors mentioning protective features
2. **Understand the root cause**: Auto scale-down during cleanup is often the culprit
3. **Apply appropriate force delete flags**: Use service-specific provider feature flags
4. **Create semantic templates**: Organize force delete configurations in reusable template functions
5. **Test the fix**: Verify that tests can create, update, and **successfully clean up** resources
