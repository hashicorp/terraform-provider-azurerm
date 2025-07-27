---
applyTo: "internal/**/*.go"
description: Azure-specific implementation patterns for the Terraform AzureRM provider including PATCH operations, CustomizeDiff patterns, and Azure SDK integration patterns.
---

# Azure-Specific Implementation Patterns

Quick navigation: [🔄 PATCH Operations](#-patch-operations) | [✅ CustomizeDiff](#-customizediff-validation) | [🎯 Schema Flattening](#-schema-flattening) | [🔐 Security](#-security-patterns)

## 🔄 PATCH Operations

### Critical PATCH Behavior Understanding

**Azure Resource Manager PATCH Operations:**
Many Azure services use PATCH operations for resource updates, which have fundamentally different behavior from PUT operations:

- **PATCH preserves existing values** when fields are omitted from the request
- **PUT replaces the entire resource** with the provided configuration
- **Azure SDK nil filtering** removes nil values before sending requests to Azure
- **Residual state persistence** means previously enabled features remain active unless explicitly disabled

### PATCH Operation Pattern

```go
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    // PATCH Operations Requirement: Always return a complete structure
    // with explicit enabled=false for disabled features to clear residual state

    // Define complete structure with all features disabled by default
    result := &azuretype.Policy{
        AutomaticFeature: &azuretype.AutomaticFeature{
            Enabled: pointer.To(false), // Explicit disable for PATCH
            // Include all required fields even when disabled
            RequiredSetting: pointer.To(azuretype.DefaultValue),
        },
        OptionalFeature: &azuretype.OptionalFeature{
            Enabled: pointer.To(false), // Explicit disable for PATCH
        },
    }

    // If no configuration, return everything disabled (clears residual state)
    if len(input) == 0 || input[0] == nil {
        return result
    }

    raw := input[0].(map[string]interface{})

    // Enable only explicitly configured features
    if automaticRaw, exists := raw["automatic_feature"]; exists {
        automaticList := automaticRaw.([]interface{})
        if len(automaticList) > 0 && automaticList[0] != nil {
            // Enable the feature and apply user configuration
            result.AutomaticFeature.Enabled = pointer.To(true)

            automatic := automaticList[0].(map[string]interface{})
            if setting := automatic["required_setting"].(string); setting != "" {
                result.AutomaticFeature.RequiredSetting = pointer.To(azuretype.Setting(setting))
            }
        }
        // If exists but empty block, feature remains disabled
    }
    // If not exists, feature remains disabled

    return result
}
```

### Documentation Requirements for PATCH Operations

```go
// PATCH Behavior Note: Azure VMSS uses PATCH operations which preserve existing values
// when fields are omitted. This means previously enabled policies will remain active
// unless explicitly disabled with enabled=false. Sending nil values gets filtered out
// by the Azure SDK, so Azure never receives disable commands. We must explicitly
// send enabled=false for all policies that should be disabled.
```

## ✅ CustomizeDiff Validation

### Standard CustomizeDiff Pattern

**Important**: When using CustomizeDiff functions, you must import both packages:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)
```

### CustomizeDiff Implementation

```go
func resourceAzureServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceAzureServiceNameCreate,
        Read:   resourceAzureServiceNameRead,
        Update: resourceAzureServiceNameUpdate,
        Delete: resourceAzureServiceNameDelete,

        Schema: map[string]*pluginsdk.Schema{
            // Azure resource schema
        },

        CustomizeDiff: pluginsdk.CustomDiffWithAll(
            // Azure-specific validation with CustomizeDiffShim wrapper
            pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Validate Azure resource dependencies
                if diff.Get("sku_name").(string) == "Premium" &&
                   diff.Get("zone_redundant").(bool) == false {
                    return fmt.Errorf("`zone_redundant` must be true for Premium SKU")
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

### Boolean Comparison Best Practices in CustomizeDiff

**Simplified Boolean Expressions:**
```go
// PREFERRED - Simplified boolean expressions
pluginsdk.ForceNewIfChange("resilient_vm_creation_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
    fieldExists := !d.GetRawConfig().GetAttr("resilient_vm_creation_enabled").IsNull()
    return fieldExists && old.(bool) && !new.(bool)
}),

// AVOID - Verbose expressions that trigger gosimple linting errors
return fieldExists && old.(bool) == true && new.(bool) == false
```

### Azure-Specific CustomizeDiff Use Cases

- **SKU validation**: Ensure Azure SKU combinations are valid
- **Location constraints**: Validate region-specific feature availability
- **Resource dependencies**: Check Azure resource prerequisite relationships
- **API version compatibility**: Ensure feature combinations match Azure API versions
- **Performance tier validation**: Validate Azure performance tier constraints
- **Field conditional validation**: Validate field combinations based on Azure API constraints

## 🎯 Schema Flattening

### When to Apply Schema Flattening

Schema flattening should be considered when Azure APIs contain unnecessary wrapper structures that don't provide value to Terraform users:

- **Single-purpose wrappers**: Remove intermediate blocks that only contain a single array or enable flag
- **Azure API convenience structures**: Eliminate wrapper objects that exist purely for API organization
- **User experience improvement**: Flatten when it simplifies configuration without losing functionality
- **Logical grouping preservation**: Maintain nested structures when they provide logical organization

### Schema Flattening Example

**Before Flattening (Complex Structure):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  log_scrubbing {
    enabled = true

    scrubbing_rule {
      match_variable = "QueryStringArgNames"
    }
  }
}
```

**After Flattening (Simplified Structure):**
```hcl
resource "azurerm_cdn_frontdoor_profile" "example" {
  name = "example"

  log_scrubbing_rule {
    match_variable = "QueryStringArgNames"
  }
}
```

### Implementation Pattern for Schema Flattening

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

## 🚫 "None" Value Pattern

### The "None" Value Pattern

Many Azure APIs accept values like None, Off, or Default as default values. The provider is moving away from exposing these values directly to users, instead leveraging Terraform's native null handling by allowing fields to be omitted.

**Modern Approach (Preferred):**
```go
// Schema excludes the "None" value - users omit the field instead
"shutdown_on_idle": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(azureapi.ShutdownOnIdleModeUserAbsence),
        string(azureapi.ShutdownOnIdleModeLowUsage),
        // Note: "None" value exists but is handled in Create/Update and Read functions
        // NOT exposed in validation
    }, false),
},
```

**Typed Resource Implementation:**
```go
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Default to "None" if user did not specify a value
            shutdownOnIdle := string(azureapi.ShutdownOnIdleModeNone)
            if model.ShutdownOnIdle != "" {
                shutdownOnIdle = model.ShutdownOnIdle
            }

            properties := azureapi.ServiceProperties{
                ShutdownOnIdle: &shutdownOnIdle,
            }

            // ...continue with resource creation
            return nil
        },
    }
}

func (r ServiceResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // ...retrieve resource from Azure

            model := ServiceModel{}

            // Only set value in state if it is not "None"
            shutdownOnIdle := ""
            if props.ShutdownOnIdle != nil && *props.ShutdownOnIdle != string(azureapi.ShutdownOnIdleModeNone) {
                shutdownOnIdle = string(*props.ShutdownOnIdle)
            }
            model.ShutdownOnIdle = shutdownOnIdle
            // If Azure returns "None", field remains empty in Terraform state

            return metadata.Encode(&model)
        },
    }
}
```

## 🔐 Security Patterns

### Credential and Secret Management

**Never Log Sensitive Information:**
```go
// GOOD - No sensitive data in logs
metadata.Logger.Infof("Creating Storage Account %s", id.StorageAccountName)
log.Printf("[DEBUG] Configuring network rules for %s", id)

// BAD - Sensitive data in logs
log.Printf("[DEBUG] Connection string: %s", connectionString) // Never log connection strings
metadata.Logger.Debugf("Client secret: %s", clientSecret)     // Never log secrets
log.Printf("[DEBUG] SAS token: %s", sasToken)                 // Never log tokens
```

**Secure Environment Variable Handling:**
```go
func validateTestCredentials() error {
    requiredVars := []string{
        "ARM_SUBSCRIPTION_ID",
        "ARM_CLIENT_ID",
        "ARM_CLIENT_SECRET",
        "ARM_TENANT_ID",
    }

    for _, envVar := range requiredVars {
        if value := os.Getenv(envVar); value == "" {
            return fmt.Errorf("required environment variable %s is not set", envVar)
        }
    }
    return nil
}
```

### Input Validation and Sanitization

**Prevent Injection Attacks:**
```go
func ValidateAzureResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Validate length
    if len(value) < 1 || len(value) > 64 {
        errors = append(errors, fmt.Errorf("property %s must be between 1 and 64 characters, got %d", k, len(value)))
        return warnings, errors
    }

    // Validate allowed characters only (prevent injection)
    allowedPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
    if !allowedPattern.MatchString(value) {
        errors = append(errors, fmt.Errorf("property %s can only contain alphanumeric characters, hyphens, and underscores", k))
        return warnings, errors
    }

    // Azure Storage Account specific reserved names
    reservedNames := []string{"admin", "root", "system", "default"}
    for _, reserved := range reservedNames {
        if strings.EqualFold(value, reserved) {
            errors = append(errors, fmt.Errorf("property `%s` cannot use reserved name `%s`", k, reserved))
            return warnings, errors
        }
    }

    return warnings, errors
}
```

## 🔄 State Management with d.GetRawConfig()

### When to Use `d.GetRawConfig()` vs `d.Get()`

**IMPORTANT**: This pattern is only available in untyped Plugin SDK resource implementations.

`d.GetRawConfig()` should be used in specific scenarios where you need to distinguish between user-configured values and computed/default values.

**Appropriate Use Cases:**
```go
// 1. Detecting if a user explicitly set a value vs using a default
func resourceServiceNameUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient

    id, err := parse.ServiceNameID(d.Id())
    if err != nil {
        return err
    }

    parameters := serviceapi.UpdateParameters{
        Name: d.Get("name").(string),
    }

    // Check if user explicitly configured the setting
    if raw := d.GetRawConfig().GetAttr("timeout_seconds"); !raw.IsNull() {
        // User explicitly set this value, use it
        timeoutValue := d.Get("timeout_seconds").(int)
        parameters.TimeoutSeconds = &timeoutValue
    }
    // If raw is null, don't send timeout_seconds parameter to Azure API
    // This allows Azure to use its service default

    if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
        return fmt.Errorf("updating %s: %+v", *id, err)
    }

    return nil
}

// 2. Handling optional complex blocks that should be omitted when not configured
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // ... standard setup code

    parameters := serviceapi.CreateParameters{
        Name:     name,
        Location: location,
    }

    // Only include advanced_config if user explicitly configured it
    if raw := d.GetRawConfig().GetAttr("advanced_config"); !raw.IsNull() {
        advancedConfig := expandAdvancedConfig(d.Get("advanced_config").([]interface{}))
        parameters.AdvancedConfig = &advancedConfig
    }
    // If raw is null, don't include AdvancedConfig in API call

    // ... continue with resource creation
}
```

**When NOT to Use `d.GetRawConfig()`:**
- Required fields (always use `d.Get()`)
- When you always need the value regardless of how it was set
- In typed resource implementations (use metadata patterns instead)
- When simple default handling in schema is sufficient

## 🏗️ Progressive Code Simplification

### When Complex Logic Needs Simplification

When implementing Azure resource expand/flatten functions, especially for PATCH operations, complex conditional logic can often be simplified through strategic refactoring.

**Step 2: Define Complete Disabled Structure (Recommended)**
```go
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    // Define complete result with all features disabled by default
    result := &azuretype.Policy{
        FeatureA: &azuretype.FeatureA{
            Enabled: pointer.To(false),                        // Disabled by default
            RequiredField: pointer.To(azuretype.DefaultValue), // Required even when disabled
        },
        FeatureB: &azuretype.FeatureB{
            Enabled: pointer.To(false), // Disabled by default
        },
    }

    // If no input, return everything disabled
    if len(input) == 0 || input[0] == nil {
        return result
    }

    raw := input[0].(map[string]interface{})

    // Simple field flipping logic - enable only what's configured
    if featureARaw, exists := raw["feature_a"]; exists {
        featureAList := featureARaw.([]interface{})
        if len(featureAList) > 0 && featureAList[0] != nil {
            // Enable the feature and apply user configuration
            result.FeatureA.Enabled = pointer.To(true)
            // Apply other configuration...
        }
    }

    if featureBEnabled, exists := raw["feature_b_enabled"]; exists {
        result.FeatureB.Enabled = pointer.To(featureBEnabled.(bool))
    }

    return result
}
```

**Key Simplification Principles:**
1. **Define end state first** - Create complete structure with desired defaults
2. **Use simple field flipping** - Change only what needs to change based on input
3. **Eliminate conditional returns** - Single return path reduces complexity
4. **Extract common patterns** - Use variables for repeated structures
5. **Start with working code** - Simplify incrementally, don't rewrite from scratch

---

## Quick Reference Links

- 🏗️ **Main Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ⚡ **Quick References**: [reference/](./reference/)
