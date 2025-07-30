---
applyTo: "internal/**/*.go"
description: Testing guidelines for Terraform AzureRM provider Go files - test execution protocols, patterns, and Azure-specific considerations.
---

# 🧪 Testing Guidelines

## 🚨 CRITICAL: Test Execution Protocol

**⚠️ NEVER RUN TESTS AUTOMATICALLY ⚠️**

**Rules:**
- **DO NOT** execute `make testacc`, `go test`, or any test commands automatically
- **ALWAYS** provide exact commands for users to run manually
- **ALWAYS** explain test purpose, duration, and Azure resource costs
- Tests create **REAL AZURE RESOURCES** and require **VALID CREDENTIALS**

**Example Command Format:**
```bash
# Purpose: Test VMSS resiliency policy backward compatibility
# Duration: 5-10 minutes, creates test VMSS resources in Azure
# Requires: ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID

make testacc TEST=./internal/services/compute TESTARGS='-run=TestAccLinuxVirtualMachineScaleSet_fieldsNotSetInState'
```

## 🚨 ENFORCEMENT RULES FOR TERMINAL TOOL USAGE

**BEFORE using any terminal tool (run_in_terminal, get_terminal_output, etc.):**
- **STOP** and ask yourself: "Am I about to run a test or build command?"
- If **YES**: Provide manual command instead
- If **NO**: Proceed with tool

**run_in_terminal tool restrictions:**
- **FORBIDDEN for**: make testacc, go test, go build, any compilation commands
- **ALLOWED for**: file operations, directory listing, git status only
- **WHEN IN DOUBT**: Provide manual command

**Before providing ANY terminal command, answer:**
1. Does this create Azure resources? → Manual command required
2. Does this run tests? → Manual command required
3. Does this build/compile code? → Manual command required
4. Is this just file inspection? → Tool may be acceptable

**If you catch yourself about to run a build/test command:**
1. **STOP immediately**
2. Provide this exact format: "Please run this command manually: [command]"
3. Explain purpose, duration, and requirements
4. **DO NOT** use run_in_terminal

## Test Types

**Unit Tests:**
- Place in same package with `_test.go` suffix
- Test utility functions, parsers, validators
- Use table-driven patterns
- No Azure credentials required

**Acceptance Tests:**
- Test against real Azure APIs with live credentials
- Package naming: `package servicename_test` (external test package)
- Test CRUD operations, imports, and state management
- Use acceptance testing framework

## Naming Conventions

**Unit Tests:** `TestFunctionName_Scenario_ExpectedOutcome`
- Example: `TestParseFrontDoorProfileID_ValidID_ReturnsCorrectComponents`

**Acceptance Tests:** `TestAccResourceName_scenario`
- Example: `TestAccCdnFrontDoorProfile_basic`
- Example: `TestAccCdnFrontDoorProfile_requiresImport`
- Use underscores to separate logical components: `TestAccResourceName_featureGroup_specificScenario`
- Example: `TestAccWindowsVirtualMachineScaleSet_skuProfile_Prioritized`

**Test Helper Functions:** Use camelCase (Go convention for unexported functions)
- Example: `skuProfilePrioritized(data acceptance.TestData) string`
- Example: `withLogScrubbingRule(data acceptance.TestData) string`
- Example: `basicConfiguration(data acceptance.TestData) string`

**Key Distinction:**
- **Test function names**: Use underscores for logical separation (`_featureGroup_scenario`)
- **Helper function names**: Use camelCase following Go naming conventions for unexported functions

## Essential Test Patterns

**Basic Resource Test:**
```go
func TestAccResourceName_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource_name", "test")
    r := ResourceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(), // Validates all field values automatically
    })
}
```

**RequiresImport Test:**
```go
func TestAccResourceName_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource_name", "test")
    r := ResourceNameResource{}
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

#### **Rule #4: User Confirmation Required**
- Before providing test commands, confirm user wants to run tests
- Verify user has Azure credentials configured
- Confirm user understands costs and time requirements
- Provide cleanup verification steps

**This protocol prevents accidental resource creation, unexpected costs, and ensures users maintain control over their Azure environment.**

---

## CustomizeDiff Testing (MANDATORY)

**Why Critical:**
- CustomizeDiff prevents invalid Azure API calls
- Enforces Azure service field combination requirements
- Provides clear error messages before resource operations

**Required Test Coverage:**
- **Error scenarios**: Test invalid field combinations with `ExpectError: regexp.MustCompile()`
- **Success scenarios**: Test valid configurations that pass validation
- **Edge cases**: Test boundary conditions and Azure service constraints

**CustomizeDiff Test Pattern:**
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
        data.ImportStep(),
    })
}
```

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

**For complete CustomizeDiff implementation patterns, import requirements, and detailed examples, see:** [Implementation Guide - CustomizeDiff Import Requirements](./implementation-guide.instructions.md#customizediff-import-requirements)

**Testing Azure-Specific CustomizeDiff Validation:**

**Essential Test Coverage:**
- **Error scenarios**: Test invalid field combinations with `ExpectError: regexp.MustCompile()`
- **Success scenarios**: Test valid configurations that pass validation
- **Edge cases**: Test boundary conditions and Azure service constraints

**Key Testing Requirements:**
- Test Azure service-specific constraints (SKU dependencies, region limitations, etc.)
- Validate Azure API field combination requirements
- Test Azure resource lifecycle constraints
- Verify Azure service version-specific validations

**Advanced Testing Patterns:**
- Use `ResourceTestIgnoreRecreate` for CustomizeDiff ForceNew validation
- Test plan verification with ConfigPlanChecks for complex state transitions
- Validate error messages with specific regexp patterns

**For Azure-specific CustomizeDiff behaviors and validation patterns, see:** [Azure Patterns - CustomizeDiff Validation](./azure-patterns.instructions.md#customizediff-validation)

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

**🚨 CRITICAL RULE: DO NOT ADD REDUNDANT FIELD VALIDATION CHECKS**

**MANDATORY Pattern - Only ExistsInAzure Check:**
```go
func TestAccCdnFrontDoorProfile_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r), // ONLY THIS CHECK - verifies resource exists
                // FORBIDDEN: check.That(data.ResourceName).Key("name").HasValue(...) - ImportStep validates this
                // FORBIDDEN: check.That(data.ResourceName).Key("sku_name").HasValue(...) - ImportStep validates this
                // FORBIDDEN: check.That(data.ResourceName).Key("field").HasValue(...) - ImportStep validates this
            ),
        },
        data.ImportStep(), // Automatically validates ALL configured field values
    })
}
```

**EXTREMELY RARE Exceptions - Only for Special Cases:**
Additional validation checks are **FORBIDDEN** except in these extremely rare circumstances:

1. **Computed Field Verification**: Testing computed values that aren't in the configuration
2. **Complex Behavior Validation**: Testing TypeSet deduplication, ordering, or transformation logic
3. **Azure API-Specific Behavior**: Testing Azure service-specific transformations

**⚠️ WARNING: If you add check.That(data.ResourceName).Key() for configured fields, you are violating the guidelines and creating redundant validation.**

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

### Cross-Implementation Consistency Requirements

When working with related Azure resources that have both Linux and Windows variants (like VMSS), ensure validation logic and behavior consistency:

**Validation Logic Consistency:**
- **Same validation rules**: Linux and Windows implementations must have identical CustomizeDiff validation logic
- **Field requirements**: If Windows requires field X for scenario Y, Linux must have the same requirement
- **Error messages**: Use identical error message patterns across related implementations
- **Default behavior**: Ensure both implementations handle defaults and omitted fields identically

**Documentation Consistency:**
- **Cross-reference checks**: When updating documentation for one variant, verify the other variant's documentation
- **Field descriptions**: Use identical descriptions for shared fields across resource variants
- **Note blocks**: Apply the same conditional logic notes to both implementations
- **Examples**: Ensure examples demonstrate the same patterns across variants

**Testing Consistency:**
- **Test coverage parity**: Both implementations should have equivalent test scenarios
- **Naming conventions**: Use parallel naming patterns (`Linux...` vs `Windows...`)
- **Helper function patterns**: Use consistent camelCase naming for test helper functions
- **Configuration templates**: Maintain similar test configuration structures

**Example Validation Consistency Pattern:**
```go
// Linux VMSS validation - must match Windows behavior
if diff.Get("sku_profile_allocation_strategy").(string) == string(compute.AllocationStrategyLowestPrice) ||
   diff.Get("sku_profile_allocation_strategy").(string) == string(compute.AllocationStrategyPrioritized) {
    // rank field validation logic must be identical
}

// Windows VMSS validation - must match Linux behavior
if diff.Get("sku_profile_allocation_strategy").(string) == string(compute.AllocationStrategyLowestPrice) ||
   diff.Get("sku_profile_allocation_strategy").(string) == string(compute.AllocationStrategyPrioritized) {
    // rank field validation logic must be identical
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

#### Optional+Computed Schema Fields Testing Requirements

When testing resources with Optional+Computed (O+C) schema fields, **comprehensive testing is mandatory** to ensure correct state management and Azure API integration. O+C fields have complex behavior that requires specific testing patterns.

**Testing Requirements for O+C Fields:**
- **Initial State Testing**: Verify the field behavior when not explicitly set by the user
- **Explicit Configuration Testing**: Test behavior when the field is explicitly configured by the user
- **Computed State Testing**: Verify that Azure-managed state changes are correctly reflected in Terraform
- **Import State Testing**: Ensure import functionality handles both user-set and Azure-managed states correctly
- **Backward Compatibility Testing**: Test that existing resources don't show unexpected diffs when upgrading the provider

**O+C Field Testing Pattern:**
```go
func TestAccServiceName_optionalComputedField(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            // Test initial state - field not explicitly set
            Config: r.withoutOptionalComputedField(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // Verify Azure default behavior is reflected
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
            ),
        },
        data.ImportStep(), // Verify import works with Azure-managed default
        {
            // Test explicit configuration by user
            Config: r.withOptionalComputedFieldEnabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(), // Verify import works with user-configured value
        {
            // Test that irreversible changes cannot be reverted (if applicable)
            Config: r.withOptionalComputedFieldDisabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // For irreversible Azure features, value should remain true
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(), // Verify final state import
    })
}
```

**Azure-Specific O+C Testing Considerations:**
- **Irreversible Features**: Test that Azure features that cannot be disabled maintain their enabled state
- **Service Defaults**: Verify that Azure service defaults are correctly detected and preserved
- **API Behavior Validation**: Ensure tests reflect actual Azure API behavior, not just Terraform logic
- **Cross-Resource Dependencies**: Test O+C fields that depend on other resource configurations

**Common O+C Testing Scenarios:**
- **VM Scale Set Resiliency Policies**: Test irreversible enablement behavior
- **Database High Availability**: Test Azure-managed redundancy settings
- **Network Security Features**: Test security policies that cannot be downgraded
- **Backup and Retention Policies**: Test Azure-managed retention configurations

### Environment Setup

## Test Organization and Placement Rules

**Acceptance Test File Structure:**
- **Test function placement**: All test functions must be placed before the `Exists` function in the test file
- **Helper function placement**: Test configuration helper functions should be placed after the `Exists` function
- **No duplicate functions**: Remove any duplicate or old test functions to maintain clean file structure
- **Consistent ordering**: Place tests in logical order (basic, update, requires import, other scenarios)

**Test Case Consolidation Standards:**
- **Basic Test**: Core functionality with minimal configuration
- **Update Test**: Resource update scenarios (only if resource supports updates)
- **RequiresImport Test**: Import conflict detection
- **Complete Test**: Full feature demonstration (optional, for complex resources)

**Example of Proper Test Organization:**
```go
package cdn_test

// ESSENTIAL TESTS - Keep these
func TestAccCdnFrontDoorProfile_basic(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_update(t *testing.T) { ... }
func TestAccCdnFrontDoorProfile_requiresImport(t *testing.T) { ... }

// Exists function - SEPARATOR between tests and helpers
func (CdnFrontDoorProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) { ... }

// Helper functions - AFTER Exists function
func (r CdnFrontDoorProfileResource) basic(data acceptance.TestData) string { ... }
func (r CdnFrontDoorProfileResource) requiresImport(data acceptance.TestData) string { ... }
```
## Azure Test Cleanup Issues

**Problem:** Azure resources with protective features block test cleanup.

**Solution:** Use provider feature flags to force deletion:
```go
provider "azurerm" {
  features {
    virtual_machine_scale_set {
      force_delete = true
    }
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}
```

**When to Use:**
- VMSS with resiliency enabled
- Key Vault with soft delete
- SQL databases with backup protection
- Any resource blocking normal cleanup

## Environment Setup

**Required Environment Variables:**
```bash
export ARM_SUBSCRIPTION_ID="your-azure-subscription-id"
export ARM_CLIENT_ID="your-service-principal-client-id"
export ARM_CLIENT_SECRET="your-service-principal-client-secret"
export ARM_TENANT_ID="your-azure-tenant-id"
export ARM_TEST_LOCATION=WestEurope
export ARM_TEST_LOCATION_ALT=EastUS2
```

**Running Tests:**
```bash
# Unit tests
go test ./internal/services/cdn/...

# Acceptance tests (MANUAL EXECUTION ONLY)
make testacc TEST=./internal/services/cdn TESTARGS='-run=TestAccCdnFrontDoorProfile_basic'
```
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

---
[⬆️ Back to top](#🧪-testing-guidelines)
