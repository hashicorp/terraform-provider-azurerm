---
applyTo: "internal/**/*.go"
description: Azure-specific guidelines for Go files in the Terraform Azure Provider repository. It includes best practices for Azure Resource Manager integration, Terraform provider patterns, and resource implementation.
---

# 🏢 Provider Guidelines

Azure-specific guidelines for Go files in the Terraform Azure Provider repository. It includes best practices for Azure Resource Manager integration, Terraform provider patterns, and resource implementation.

**Quick navigation:** [☁️ ARM Integration](#☁️-azure-resource-manager-arm-integration) | [⚙️ CustomizeDiff](#⚙️-customizediff-implementation-for-azure-resources) | [📐 Schema Design](#📐-azure-schema-design-and-flattening-guidelines) | [✅ API Validation](#✅-azure-api-value-validation)

## ☁️ Azure Resource Manager (ARM) Integration
- Use the HashiCorp Go Azure SDK as the primary SDK for Azure integrations
- Implement proper error handling for Azure API responses
- Use appropriate polling for long-running operations (LROs)
- Implement proper retry logic with exponential backoff
- Handle Azure API rate limits and throttling gracefully
- Use managed identity authentication when possible
- Always validate resource IDs using the proper parsing utilities

### Implementation Approach Guidelines

**Typed Resource Implementation (Preferred):**
- Use type-safe model structures with `tfschema` tags for Azure resource properties
- Leverage `metadata.Client` for Azure SDK client access and structured logging
- Implement proper resource ID validation with Azure-specific parsing utilities
- Use `metadata.Decode()` and `metadata.Encode()` for state management
- Handle Azure resource lifecycle with proper import conflict detection
- Implement Azure resource existence checks using `metadata.MarkAsGone()`

**UnTyped Resource Implementation (Maintenance Only):**
- Use direct schema manipulation with proper Azure resource validation
- Implement Azure API client initialization through `clients.Client`
- Handle Azure resource state using `d.Set()` and `d.Get()` patterns
- Use `tf.ImportAsExistsError()` for Azure resource import conflicts
- Implement Azure resource cleanup with proper error handling

**Common Azure Requirements (Both Approaches):**
- Follow the standard CRUD lifecycle: Create, Read, Update, Delete
- Implement proper state management and drift detection for Azure resources
- Use `ForceNew` for Azure properties that require resource recreation
- Implement proper timeout configurations for Azure long-running operations
- Use appropriate validation functions for Azure resource properties
- Handle nested Azure resource configurations properly using `TypeSet`, `TypeList`, and `TypeMap`

---
[⬆️ Back to top](#🏢-provider-guidelines)

## ⚙️ CustomizeDiff Implementation for Azure Resources

### Standard CustomizeDiff Pattern

**Note:** CustomizeDiff implementation differs between typed and untyped resources. Choose the appropriate pattern based on your resource implementation approach.

**Untyped Resource CustomizeDiff Pattern:**
```go
import (
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // Only this import needed
)

func resourceAzureServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceAzureServiceNameCreate,
        Read:   resourceAzureServiceNameRead,
        Update: resourceAzureServiceNameUpdate,
        Delete: resourceAzureServiceNameDelete,

        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },

        Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
            _, err := parse.ServiceNameID(id)
            return err
        }),

        Schema: map[string]*pluginsdk.Schema{
            // Azure resource schema
        },

        CustomizeDiff: pluginsdk.CustomDiffWithAll(
            // Azure-specific validation with CustomizeDiffShim wrapper
            pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
                // Validate Azure resource dependencies
                if diff.Get("sku_name").(string) == "Premium" &&
                   diff.Get("zone_redundant").(bool) == false {
                    return fmt.Errorf("`zone_redundant` must be `true` for `Premium` SKU")
                }
                return nil
            }),
            // Force recreation for Azure resource properties that require it
            pluginsdk.ForceNewIfChange("location", func(ctx context.Context, old, new, meta interface{}) bool {
                return old.(string) != new.(string)
            }),
        ),
    }
}
```

**Typed Resource CustomizeDiff Pattern:**
```go
// Typed resources use receiver methods with sdk.ResourceFunc
func (r ServiceNameResource) CustomizeDiff() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Azure SKU validation for typed resources
            if model.SkuName == "Premium" && !model.ZoneRedundant {
                return fmt.Errorf("`zone_redundant` must be `true` for `Premium` SKU")
            }

            return nil
        },
    }
}
```

### Azure-Specific CustomizeDiff Use Cases
- **SKU validation**: Ensure Azure SKU combinations are valid
- **Location constraints**: Validate region-specific feature availability
- **Resource dependencies**: Check Azure resource prerequisite relationships
- **API version compatibility**: Ensure feature combinations match Azure API versions
- **Performance tier validation**: Validate Azure performance tier constraints
- **Field conditional validation**: Validate field combinations based on Azure API constraints

### Boolean Comparison Best Practices in CustomizeDiff

**Simplified Boolean Expressions:**
```go
// PREFERRED - Simplified boolean expressions
pluginsdk.ForceNewIfChange("resilient_vm_creation_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
    fieldExists := !d.GetRawConfig().GetAttr("resilient_vm_creation_enabled").IsNull()
    return fieldExists && old.(bool) && !new.(bool)
}),
```

**FORBIDDEN - Verbose boolean comparisons:**
```go
// FORBIDDEN - Verbose expressions that trigger gosimple linting errors
return fieldExists && old.(bool) == true && new.(bool) == false
```

**Key Principles:**
- Use direct boolean expressions: `old.(bool) && !new.(bool)`
- Leverage Go's boolean semantics: `bool` values can be used directly in logical expressions
- Comply with linting standards: Simplified expressions pass gosimple and other Go linting tools
- Maintain readability: Shorter expressions are easier to understand and maintain
- `CustomizeDiff` validations should be thoroughly tested with acceptance tests to ensure they work correctly:
- Test invalid configurations that should trigger validation errors
- Test valid configurations that should pass validation
- Test edge cases and boundary conditions
- Use `ExpectError: regexp.MustCompile()` to verify specific error messages
- Test all possible combinations of field values that trigger different validation paths

For comprehensive `CustomizeDiff` testing examples, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

---
[⬆️ Back to top](#🏢-provider-guidelines)

## 📐 Azure Schema Design and Flattening Guidelines

### Schema Flattening for Improved User Experience

**When to Apply Schema Flattening:**
Schema flattening should be considered when Azure APIs contain unnecessary wrapper structures that don't provide value to Terraform users. The goal is to create intuitive, user-friendly schemas that reflect the logical structure of the Azure resource configuration.

**Flattening Decision Criteria:**
- **Single-purpose wrappers**: Remove intermediate blocks that only contain a single array or enable flag
- **Azure API convenience structures**: Eliminate wrapper objects that exist purely for API organization
- **User experience improvement**: Flatten when it simplifies configuration without losing functionality
- **Logical grouping preservation**: Maintain nested structures when they provide logical organization

**Example: CDN Front Door Profile Log Scrubbing Flattening**

**Before Flattening (Complex Structure):**
```go
# Complex nested structure with unnecessary wrapper
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }

    scrubbing_rule {
      match_variable = "RequestIPAddress"
    }
  }
}
```

**After Flattening (Simplified Structure):**
```go
# Flattened structure - direct access to scrubbing rules
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  log_scrubbing_rule {
    match_variable = "QueryStringArgNames"
  }

  log_scrubbing_rule {
    match_variable = "RequestIPAddress"
  }
}
```

**Implementation Pattern for Schema Flattening:**

```go
// Schema definition - direct access to the meaningful configuration
"log_scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    MaxItems: 3,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "match_variable": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc: validation.StringInSlice(
                    profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
                    false),
            },
        },
    },
},

// Expand function - handle the wrapper structure internally
func expandCdnFrontDoorProfileLogScrubbing(input []interface{}) *profiles.ProfileLogScrubbing {
    if len(input) == 0 {
        // When no rules configured, set to disabled (following "None" pattern)
        policyDisabled := profiles.ProfileScrubbingStateDisabled
        return &profiles.ProfileLogScrubbing{
            State:          &policyDisabled,
            ScrubbingRules: nil,
        }
    }

    // When rules are present, always enable the feature
    policyEnabled := profiles.ProfileScrubbingStateEnabled
    scrubbingRules := expandScrubbingRules(input)

    return &profiles.ProfileLogScrubbing{
        State:          &policyEnabled,
        ScrubbingRules: scrubbingRules,
    }
}

// Flatten function - hide wrapper complexity from users
func flattenCdnFrontDoorProfileLogScrubbing(input *profiles.ProfileLogScrubbing) []interface{} {
    if input == nil || pointer.From(input.State) == profiles.ProfileScrubbingStateDisabled {
        // When disabled, return empty list (following "None" pattern)
        return make([]interface{}, 0)
    }

    // Return only the meaningful rules, hiding the wrapper
    return flattenScrubbingRules(input.ScrubbingRules)
}
```

**Schema Flattening Best Practices:**
- **Preserve Azure API contract**: Handle wrapper structures in expand/flatten functions
- **Follow "None" pattern**: Use empty arrays/blocks to represent disabled features
- **Maintain backward compatibility**: Ensure flattening doesn't break existing configurations
- **Document the simplification**: Clearly explain the flattened structure in documentation
- **Test thoroughly**: Verify that all Azure API combinations work correctly
- **Consider feature flags**: Use provider features for gradual rollout if needed

**Anti-Patterns to Avoid:**
- **Over-flattening**: Don't remove meaningful logical groupings
- **Breaking changes**: Avoid flattening that would break existing user configurations
- **Lost context**: Ensure flattening doesn't lose important Azure resource context
- **Validation complexity**: Don't flatten if it makes validation significantly more complex

**When NOT to Flatten:**
- **Complex validation dependencies**: When nested structure enables clearer validation rules
- **Multiple configuration modes**: When wrapper provides distinct configuration patterns
- **Azure service evolution**: When Azure might add more properties to the wrapper
- **Breaking change risk**: When existing users would be significantly impacted

---
[⬆️ Back to top](#🏢-provider-guidelines)

## ✅ Azure API Value Validation

### Valid Value Documentation Standards
- **Azure SDK Alignment**: Always verify valid values against Azure SDK enum constants before documenting
- **API Documentation Cross-Reference**: Check Azure REST API documentation to confirm supported values
- **Service-Specific Validation**: Different Azure services may support different subsets of common values
- **Version-Specific Features**: Ensure documented values are available in the API version being used

### Implementation Requirements
- **Schema Validation**: Only include actually supported values in schema validation functions
- **Documentation Accuracy**: Documentation must only show values that work with the specific Azure service
- **Test Configuration**: Test configurations must use only valid Azure service values
- **Error Messages**: Validation errors should reference only valid values for the specific service

Example of proper Azure value validation:
```go
// PREFERRED - Use Azure SDK PossibleValues function when available
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
        false,
    ),
},

// ONLY when SDK doesn't expose PossibleValues function - Manual validation
"custom_field": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        // Only include values from Azure SDK constants that work with this specific service
        string(azureapi.CustomFieldValueOptionA),
        string(azureapi.CustomFieldValueOptionB),
    }, false),
},
```

**AVOID - Including values from other services:**
```go
// Don't include SDK constants from unrelated Azure services
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(profiles.ScrubbingRuleEntryMatchVariableQueryStringArgNames), // Correct for this service
        string(compute.VirtualMachineEvictionPolicyDeallocate),              // Wrong - from different service
    }, false),
},
```

## 📚 Related Implementation Guidance (On-Demand)

### **Schema & Testing**
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md) - Azure schema design and validation

### **Quality & Evolution**
- 📋 **Code Clarity**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md) - Code quality standards
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md) - Provider evolution patterns
- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md) - Azure API versioning

### **Specialized Areas**
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md) - Azure security patterns
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md) - Azure provider optimization

---
[⬆️ Back to top](#🏢-provider-guidelines)
