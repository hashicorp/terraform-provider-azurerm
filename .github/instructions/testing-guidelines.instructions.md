---
applyTo: "internal/**/*.go"
description: Testing guidelines for Terraform AzureRM provider Go files - test execution protocols, patterns, and Azure-specific considerations.
---

# 🧪 Testing Guidelines

Testing guidelines for Terraform AzureRM provider Go files - test execution protocols, patterns, and Azure-specific considerations.

**Quick navigation:** [🚨 Test Execution Policy](#🚨-critical-test-execution-policy-🚨) | [🚨 Redundant Validation Policy](#🚨-critical-redundant-validation-checks-with-import-step-policy-🚨) | [🧪 Test Types](#🧪-test-types) | [⚡ Essential Patterns](#⚡-essential-test-patterns) | [✅ CustomizeDiff Testing](#✅-customizediff-testing-mandatory) | [📊 Data Source Testing](#📊-data-source-testing-patterns) | [🏗️ Test Organization](#🏗️-test-organization-and-placement-rules) | [☁️ Azure-Specific Testing](#☁️-azure-specific-testing-guidelines) | [🔧 Environment Setup](#🔧-environment-setup)

## 🚨 CRITICAL: TEST EXECUTION POLICY 🚨

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

**MANDATORY SELF-CHECK BEFORE ANY TERMINAL COMMAND:**

**Before using any terminal tool, AI MUST answer these questions:**
1. "Does this command run tests?" → If YES: **AUTOMATIC VIOLATION - PROVIDE MANUAL COMMAND**
2. "Does this command build/compile?" → If YES: **AUTOMATIC VIOLATION - PROVIDE MANUAL COMMAND**
3. "Does this create Azure resources?" → If YES: **AUTOMATIC VIOLATION - PROVIDE MANUAL COMMAND**
4. "Is this only file inspection?" → If YES: Tool may be acceptable

**🚫 AUTOMATIC VIOLATIONS - NEVER USE run_in_terminal FOR:**
- `make testacc` - Azure resource creation
- `go test` - Test execution
- `go build` - Compilation
- `terraform plan/apply` - Infrastructure changes
- Any command creating billable Azure resources

**✅ ACCEPTABLE run_in_terminal USAGE:**
- `ls`, `dir` - Directory listing
- `cat`, `Get-Content` - File reading
- `git status` - Repository status
- File operations that don't execute code

**🔄 VIOLATION RESPONSE PROTOCOL:**
If you catch yourself about to run a forbidden command:
1. **STOP immediately**
2. **NEVER** use run_in_terminal
3. Provide manual command with this exact format:
   ```
   Please run this command manually:
   [command]

   Purpose: [what this does]
   Duration: [expected time]
   Requirements: [prerequisites]
   ```

---
[⬆️ Back to top](#🧪-testing-guidelines)

## 🚨 CRITICAL: REDUNDANT VALIDATION CHECKS WITH IMPORT STEP POLICY 🚨

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

**Key Principles:**
- **ImportStep handles field validation**: Don't duplicate validation of configured field values
- **Keep only ExistsInAzure**: Essential for verifying resource creation/existence
- **Add checks sparingly**: Only for behavior that ImportStep cannot verify
- **Document rationale**: Comment why additional checks are needed when used

---
[⬆️ Back to top](#🧪-testing-guidelines)

## 🧪 Test Types

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

### Naming Conventions

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

### Go Testing Patterns

**Table-Driven Tests:**
```go
func TestParseResourceID(t *testing.T) {
    testCases := []struct {
        name        string
        input       string
        expected    ResourceID
        shouldError bool
    }{
        {
            name:     "valid resource ID",
            input:    "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Service/resources/resource1",
            expected: ResourceID{SubscriptionID: "12345", ResourceGroup: "rg1", Name: "resource1"},
            shouldError: false,
        },
        {
            name:        "invalid resource ID",
            input:       "invalid-id",
            expected:    ResourceID{},
            shouldError: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := ParseResourceID(tc.input)

            if tc.shouldError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }

            if !reflect.DeepEqual(result, tc.expected) {
                t.Errorf("expected %+v, got %+v", tc.expected, result)
            }
        })
    }
}
```

**Assertion Patterns:**
```go
// Use testify assertions for cleaner test code
func TestResourceValidation(t *testing.T) {
    require := require.New(t)
    assert := assert.New(t)

    // Test setup
    resource := createTestResource()

    // Assertions
    require.NotNil(resource)
    assert.Equal("expected-value", resource.Name)
    assert.True(resource.Enabled)
    assert.Contains(resource.Tags, "environment")
}
```

---
[⬆️ Back to top](#🧪-testing-guidelines)

## ⚡ Essential Test Patterns

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

### **Rule #4: User Confirmation Required**
- Before providing test commands, confirm user wants to run tests
- Verify user has Azure credentials configured
- Confirm user understands costs and time requirements
- Provide cleanup verification steps

**This protocol prevents accidental resource creation, unexpected costs, and ensures users maintain control over their Azure environment.**

---
[⬆️ Back to top](#🧪-testing-guidelines)

## ✅ CustomizeDiff Testing (MANDATORY)

**Why Critical:**
- CustomizeDiff prevents invalid Azure API calls
- Enforces Azure service field combination requirements
- Provides clear error messages before resource operations

**Required Test Coverage:**
- **Error scenarios**: Test invalid field combinations with `ExpectError: regexp.MustCompile()`
- **Success scenarios**: Not needed, they will be tested in the other test cases (e.g., `basic`, `update`, and `complete`)
- **Edge cases**: Test boundary conditions and Azure service constraints

**CustomizeDiff Test Pattern:**
```go
func TestAccServiceName_featureName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is `true`"),
        },
    })
}
```

CustomizeDiff validations are essential for enforcing Azure API constraints and preventing invalid configurations. **Testing these validations is mandatory** and requires comprehensive coverage of both success and failure scenarios.

### Why CustomizeDiff Testing is Critical

**Azure API Constraint Enforcement:**
- CustomizeDiff validations prevent invalid API calls that would fail at runtime
- They enforce Azure service-specific field combination requirements
- They validate complex resource dependencies before Azure API interaction
- They provide clear error messages to users before resource `creation`/`update`

**Testing Requirements:**
- **Error Scenarios**: Test all invalid field combinations that should trigger validation errors
- **Success Scenarios**: Not needed, they will be tested in the other test cases (e.g., `basic`, `update`, and `complete`)
- **Edge Cases**: Test boundary conditions and corner cases
- **Error Message Validation**: Verify specific error messages using `ExpectError: regexp.MustCompile()`
- **Field Path Accuracy**: Ensure error messages include correct field paths and constraints
- **Azure API Alignment**: Test that validations match actual Azure API behavior

### CustomizeDiff Testing Mandatory Practices

**Comprehensive Test Coverage:**
```go
func TestAccServiceName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        // REQUIRED: Test invalid configuration
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is `true`"),
        },
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
- **Success scenarios**: Not required, they will be tested in the other test cases (e.g., `basic`, `update`, and `complete`)
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

---
[⬆️ Back to top](#🧪-testing-guidelines)

## Acceptance Testing Patterns

### Basic Resource Test
```go
func TestAccCdnFrontDoorProfile_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(), // No sensitive fields to exclude for CDN profiles
    })
}
```

### Resource Update Test
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
            ),
        },
        data.ImportStep(),
    })
}
```

### Resource Requires Import Test
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
---
[⬆️ Back to top](#🧪-testing-guidelines)

## 📊 Data Source Testing Patterns

Data sources have different testing requirements than resources since they retrieve existing information rather than manage resource lifecycle.

**Basic Data Source Test:**
```go
func TestAccCdnFrontDoorProfileDataSource_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_profile", "test")
    r := CdnFrontDoorProfileDataSource{}

    data.DataSourceTest(t, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                // Data sources don't have ExistsInAzure checks - they retrieve existing resources
                check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestcdnfd-%d", data.RandomInteger)),
                check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-cdn-%d", data.RandomInteger)),
                check.That(data.ResourceName).Key("sku_name").HasValue("Standard_AzureFrontDoor"),
                check.That(data.ResourceName).Key("id").Exists(),
            ),
        },
    })
}
```

**Data Source Test Configuration Pattern:**
```go
func (CdnFrontDoorProfileDataSource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_profile" "test" {
  name                = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorProfileResource{}.basic(data))
}
```

**Data Source Key Validation Requirements:**
- **Field Verification**: Data sources MUST validate that expected fields are populated with correct values
- **Computed Field Verification**: Test that computed fields (like IDs, endpoints) are populated
- **Complex Structure Validation**: Use Key validation for nested data structures retrieved from Azure
- **No ImportStep**: Data sources don't support import, so all validation must be explicit

**Valid Data Source Key Validation Examples:**
```go
// VALID: Verifying data source retrieves correct values
check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
check.That(data.ResourceName).Key("tags.Environment").HasValue("Production"),

// VALID: Validating computed fields are populated
check.That(data.ResourceName).Key("id").Exists(),
check.That(data.ResourceName).Key("endpoint").Exists(),

// VALID: Complex structure validation for data sources
check.That(data.ResourceName).Key("log_scrubbing_rule.#").HasValue("2"),
check.That(data.ResourceName).Key("log_scrubbing_rule.0.match_variable").HasValue("QueryStringArgNames"),
```
---
[⬆️ Back to top](#🧪-testing-guidelines)

## 🏗️ Test Organization and Placement Rules

### Acceptance Test File Structure
- **Test function placement**: All test functions must be placed before the `Exists` function in the test file
- **Helper function placement**: Test configuration helper functions should be placed after the `Exists` function
- **No duplicate functions**: Remove any duplicate or old test functions to maintain clean file structure
- **Consistent ordering**: Place tests in logical order (basic, update, requires import, other scenarios)

### Test Case Consolidation Standards

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

### Cross-Implementation Consistency Requirements

When working with related Azure resources that have both Linux and Windows variants (like VMSS), ensure validation logic and behavior consistency:

**Validation Logic Consistency:**
- **Same validation rules**: Linux and Windows implementations must have identical CustomizeDiff validation logic
- **Field requirements**: If Windows requires field X for scenario Y, Linux must have the same requirement
- **Error messages**: Use identical error message patterns across related implementations
- **Default behavior**: Ensure both implementations handle defaults and omitted fields identically

---
[⬆️ Back to top](#🧪-testing-guidelines)

## ☁️ Azure-Specific Testing Guidelines

### Resource Existence Checks

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

**UnTyped Resource Existence Check:**
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

### Azure Test Cleanup Issues

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

---
[⬆️ Back to top](#🧪-testing-guidelines)

## 🔧 Environment Setup

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

**Common Azure Test Cleanup Issues:**
- `ResourceGroupBeingDeleted: Cannot perform operation while resource group is being deleted`
- Scale-down operations blocked due to health monitoring requirements
- Soft-delete conflicts preventing immediate recreation

---
[⬆️ Back to top](#🧪-testing-guidelines)

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 📋 **Code Clarity Enforcement**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)

### 🚀 Enhanced Guidance Files

- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md)
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md)
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md)
- 🔧 **Troubleshooting**: [troubleshooting-decision-trees.instructions.md](./troubleshooting-decision-trees.instructions.md)

---
[⬆️ Back to top](#🧪-testing-guidelines)
