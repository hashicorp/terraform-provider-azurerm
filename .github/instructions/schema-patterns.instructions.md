---
applyTo: "internal/**/*.go"
description: Schema design patterns and validation standards for the Terraform AzureRM provider including field types, validation patterns, and Azure-specific schema considerations.
---

# Schema Design Patterns

Schema design patterns and validation standards for the Terraform AzureRM provider including field types, validation patterns, and Azure-specific schema considerations.

**Quick navigation:** [📋 Schema Types](#📋-schema-type-patterns) | [✅ Validation](#✅-validation-patterns) | [⚙️ Azure Specific](#⚙️-azure-specific-schema-patterns) | [🚀 Breaking Changes](#🚀-fivepointoh-feature-flag-patterns) | [🏗️ Complex Schemas](#🏗️-complex-schema-patterns) | [🔍 Field Naming](#🔍-field-naming-standards) | [🧪 Testing](#🧪-testing-schema-patterns) | [🔧 Test Helpers](#🔧-test-configuration-helpers)

## 📋 Schema Type Patterns

### Basic Field Types

```go
// Required string field with validation
"name": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ForceNew:     true,
    ValidateFunc: validation.StringIsNotEmpty,
},

// Optional boolean with default
"enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Default:  true,
},

// Optional integer with range validation
"timeout_seconds": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    ValidateFunc: validation.IntBetween(1, 3600),
},

// Computed string (read-only)
"endpoint": {
    Type:     pluginsdk.TypeString,
    Computed: true,
},
```

### List and Set Patterns

```go
// Simple string list
"allowed_origins": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    Elem: &pluginsdk.Schema{
        Type: pluginsdk.TypeString,
    },
},

// Set of strings with validation
"enabled_features": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Schema{
        Type:         pluginsdk.TypeString,
        ValidateFunc: validation.StringInSlice([]string{
            "FeatureA",
            "FeatureB",
            "FeatureC",
        }, false),
    },
},

// Map of string to string
"settings": {
    Type:     pluginsdk.TypeMap,
    Optional: true,
    Elem: &pluginsdk.Schema{
        Type: pluginsdk.TypeString,
    },
},
```

### Complex Nested Block Patterns

```go
// Single nested block (MaxItems: 1)
"configuration": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "setting1": {
                Type:         pluginsdk.TypeString,
                Required:     true,
                ValidateFunc: validation.StringIsNotEmpty,
            },
            "setting2": {
                Type:     pluginsdk.TypeString,
                Optional: true,
                Default:  "default_value",
            },
        },
    },
},

// Multiple nested blocks
"network_configuration": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "subnet_id": {
                Type:         pluginsdk.TypeString,
                Required:     true,
                ValidateFunc: commonids.ValidateSubnetID,
            },
            "private_ip_address": {
                Type:         pluginsdk.TypeString,
                Optional:     true,
                ValidateFunc: validation.IsIPv4Address,
            },
        },
    },
},
```

---
[⬆️ Back to top](#schema-design-patterns)

## ✅ Validation Patterns

### 🚨 Schema Definition Verification Before Field Validation

**CRITICAL RULE: AI Must Check Schema Definition Before Suggesting Field Validation Code**

Before suggesting any empty/exists checks or validation logic for fields, the AI MUST first examine the field's schema definition to determine its type and suggest the appropriate validation approach.

**Mandatory AI Pre-Code-Suggestion Analysis:**

1. **AI Must Identify Field Schema Type**:
   ```go
   // AI should examine the schema definition first
   "field_name": {
       Type:     pluginsdk.TypeString,
       Required: true,    // OR Optional: true, OR both Optional+Computed
       Computed: false,   // May be true for Optional+Computed fields
   }
   ```

2. **AI Decision Tree for Validation Pattern Selection**:

   **Required Fields** → AI should suggest direct access:
   ```go
   // Typed Resource Implementation
   var model ServiceNameModel
   if err := metadata.Decode(&model); err != nil {
       return fmt.Errorf("decoding: %+v", err)
   }
   // Use model.FieldName directly - Required fields guaranteed to have values

   // Untyped Resource Implementation  
   value := diff.Get("field_name").(string)
   // Use value directly - Required fields guaranteed to have values
   ```

   **Optional Fields** → AI should suggest raw config checking:
   ```go
   // Typed Resource Implementation
   var model ServiceNameModel
   if err := metadata.Decode(&model); err != nil {
       return fmt.Errorf("decoding: %+v", err)
   }
   // Check if user explicitly configured the field
   if !metadata.ResourceData.GetRawConfig().GetAttr("field_name").IsNull() {
       // Only validate if user explicitly configured the field
       if model.FieldName == "" {
           return fmt.Errorf("field `field_name` cannot be empty when specified")
       }
   }

   // Untyped Resource Implementation
   if !diff.GetRawConfig().GetAttr("field_name").IsNull() {
       value := diff.Get("field_name").(string)
       // Only validate if user explicitly configured the field
       if value == "" {
           return fmt.Errorf("field `field_name` cannot be empty when specified")
       }
   }
   ```

   **Optional+Computed Fields** → AI should suggest user vs Azure value distinction:
   ```go
   // Typed Resource Implementation
   var model ServiceNameModel
   if err := metadata.Decode(&model); err != nil {
       return fmt.Errorf("decoding: %+v", err)
   }
   // Distinguish user-configured vs Azure-computed values
   if !metadata.ResourceData.GetRawConfig().GetAttr("field_name").IsNull() {
       // Validate user-configured values only
       if model.FieldName == "" {
           return fmt.Errorf("field `field_name` cannot be empty when explicitly set")
       }
   }
   // Skip validation for Azure-computed values

   // Untyped Resource Implementation
   if !diff.GetRawConfig().GetAttr("field_name").IsNull() {
       value := diff.Get("field_name").(string)
       // Validate user-configured values only
       if value == "" {
           return fmt.Errorf("field `field_name` cannot be empty when explicitly set")
       }
   }
   // Skip validation for Azure-computed values
   ```

**AI Schema Analysis Checklist Before Code Suggestions:**
- [ ] Examined field schema definition (Required/Optional/Optional+Computed)
- [ ] Identified appropriate validation method based on schema type
- [ ] Avoided suggesting `GetRawConfig()` for Required fields (unnecessary overhead)
- [ ] Avoided suggesting Go zero value validation for Optional fields without raw config check
- [ ] Distinguished between user-configured and Azure-computed values for Optional+Computed fields

**AI Common Mistakes to Avoid:**
- **Zero Value Confusion**: Suggesting validation of Go zero values (`0`, `false`, `""`) without checking if user actually configured them
- **Required Field Over-validation**: Suggesting `GetRawConfig()` checks for Required fields
- **Optional Field Under-validation**: Suggesting direct validation of Optional fields without checking user intent
- **Schema Type Ignorance**: Making validation suggestions without first examining the field's schema definition

### String Validation

```go
// Basic string validation
"name": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.All(
        validation.StringIsNotEmpty,
        validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "name can only contain alphanumeric characters and hyphens"),
        validation.StringLenBetween(1, 64),
    ),
},

// Enum validation using Azure SDK
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
        false,
    ),
},

// Manual enum validation (when SDK doesn't provide function)
"sku_name": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        "Standard_S1",
        "Standard_S2",
        "Premium_P1",
    }, false),
},
```

### Integer Validation

```go
// Range validation
"port": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    Default:      80,
    ValidateFunc: validation.IntBetween(1, 65535),
},

// List of valid values
"replica_count": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    Default:      1,
    ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 5, 10}),
},
```

### Custom Validation Functions

```go
// Azure resource name validation
func ValidateAzureResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Validate length
    if len(value) < 1 || len(value) > 64 {
        errors = append(errors, fmt.Errorf("property `%s` must be between 1 and 64 characters, got %d", k, len(value)))
        return warnings, errors
    }

    // Validate allowed characters only (prevent injection)
    allowedPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
    if !allowedPattern.MatchString(value) {
        errors = append(errors, fmt.Errorf("property `%s` can only contain alphanumeric characters, hyphens, and underscores", k))
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

// SQL injection prevention
func ValidateSQLResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Check for SQL injection patterns
    sqlInjectionPatterns := []string{
        "'", "\"", ";", "--", "/*", "*/", "xp_", "sp_", "exec", "execute",
        "select", "insert", "update", "delete", "drop", "create", "alter",
    }

    lowerValue := strings.ToLower(value)
    for _, pattern := range sqlInjectionPatterns {
        if strings.Contains(lowerValue, pattern) {
            errors = append(errors, fmt.Errorf("property `%s` cannot contain potentially unsafe characters or SQL keywords", k))
            return warnings, errors
        }
    }

    return warnings, errors
}
```

---
[⬆️ Back to top](#schema-design-patterns)

## ⚙️ Azure-Specific Schema Patterns

### Common Azure Schema Helpers

```go
// Use common Azure schema helpers
"location": commonschema.Location(),
"resource_group_name": commonschema.ResourceGroupName(),
"tags": commonschema.Tags(),

// For data sources
"location": commonschema.LocationComputed(),
"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
"tags": commonschema.TagsDataSource(),
```

### The "None" Value Pattern

**Modern Approach (Preferred) - Exclude "None" from validation:**

```go
// Schema excludes the "None" value - users omit the field instead
"shutdown_on_idle": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(azureapi.ShutdownOnIdleModeUserAbsence),
        string(azureapi.ShutdownOnIdleModeLowUsage),
        // Note: While the "None" value exists it is handled in the Create/Update and Read functions.
        // NOT exposed in validation
        // string(azureapi.ShutdownOnIdleModeNone),
    }, false),
},
```

**Implementation in Create/Update:**
```go
// Convert Terraform null to Azure "None"
shutdownOnIdle := string(azureapi.ShutdownOnIdleModeNone)
if model.ShutdownOnIdle != "" {
    shutdownOnIdle = model.ShutdownOnIdle
}

properties := azureapi.ServiceProperties{
    ShutdownOnIdle: &shutdownOnIdle,
}
```

**Implementation in Read:**
```go
// Convert Azure "None" back to Terraform null
shutdownOnIdle := ""
if props.ShutdownOnIdle != nil && *props.ShutdownOnIdle != string(azureapi.ShutdownOnIdleModeNone) {
    shutdownOnIdle = *props.ShutdownOnIdle
}
model.ShutdownOnIdle = shutdownOnIdle
```

### Enhanced "None" Value Conversion Pattern

**Advanced Conversion Logic for Complex Azure APIs:**
When Azure APIs use complex enum states beyond simple "None" values, implement comprehensive state mapping:

```go
// Complex Azure API state mapping
func ExpandAdvancedFeature(input []interface{}) *azureapi.AdvancedFeature {
    // Handle multiple "disabled" states from Azure API
    defaultState := &azureapi.AdvancedFeature{
        State: pointer.To(azureapi.FeatureStateDisabled),
        Mode:  pointer.To(azureapi.FeatureModeNone),
        Level: pointer.To(azureapi.FeatureLevelOff),
    }

    if len(input) == 0 || input[0] == nil {
        return defaultState
    }

    raw := input[0].(map[string]interface{})

    return &azureapi.AdvancedFeature{
        State: pointer.To(azureapi.FeatureStateEnabled),
        Mode:  pointer.To(azureapi.FeatureMode(raw["mode"].(string))),
        Level: pointer.To(azureapi.FeatureLevel(raw["level"].(string))),
    }
}

// Complex flatten with multiple "None" equivalents
func FlattenAdvancedFeature(input *azureapi.AdvancedFeature) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    // Multiple conditions that represent "disabled/none" state
    isDisabled := pointer.From(input.State) == azureapi.FeatureStateDisabled ||
                  pointer.From(input.Mode) == azureapi.FeatureModeNone ||
                  pointer.From(input.Level) == azureapi.FeatureLevelOff

    if isDisabled {
        return make([]interface{}, 0)
    }

    return []interface{}{
        map[string]interface{}{
            "mode":  string(pointer.From(input.Mode)),
            "level": string(pointer.From(input.Level)),
        },
    }
}
```

**Conditional "None" Pattern with Dependencies:**
```go
// Complex conditional "None" handling
func ExpandConditionalFeature(input []interface{}, skuTier string) *azureapi.ConditionalFeature {
    // Default varies by SKU tier
    var defaultState azureapi.FeatureState
    switch skuTier {
    case "Premium":
        defaultState = azureapi.FeatureStateEnabled
    default:
        defaultState = azureapi.FeatureStateDisabled
    }

    if len(input) == 0 || input[0] == nil {
        return &azureapi.ConditionalFeature{
            State: pointer.To(defaultState),
        }
    }

    raw := input[0].(map[string]interface{})
    return &azureapi.ConditionalFeature{
        State:  pointer.To(azureapi.FeatureStateEnabled),
        Config: expandFeatureConfig(raw["config"].([]interface{})),
    }
}
```

### Advanced "Enabled" Property Handling Patterns

**Dynamic Enable/Disable Based on Azure Service Capabilities:**
```go
// SKU-dependent enable patterns
"advanced_security": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "threat_detection_enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
                Default:  false,
            },
            "vulnerability_assessment_enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
                Default:  false,
            },
            "advanced_data_security_enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
                Default:  false,
            },
        },
    },
},
```

**Expand Function with SKU Validation:**
```go
func ExpandAdvancedSecurity(input []interface{}, skuName string) *azureapi.SecurityConfiguration {
    // Base configuration with all features disabled
    config := &azureapi.SecurityConfiguration{
        ThreatDetection:      &azureapi.ThreatDetectionConfig{Enabled: pointer.To(false)},
        VulnerabilityAssess:  &azureapi.VulnerabilityConfig{Enabled: pointer.To(false)},
        AdvancedDataSecurity: &azureapi.AdvancedDataConfig{Enabled: pointer.To(false)},
    }

    if len(input) == 0 || input[0] == nil {
        return config
    }

    raw := input[0].(map[string]interface{})

    // Enable features based on configuration and SKU capabilities
    if skuName == "Premium" || skuName == "BusinessCritical" {
        if threatEnabled := raw["threat_detection_enabled"].(bool); threatEnabled {
            config.ThreatDetection.Enabled = pointer.To(true)
        }

        if vulnEnabled := raw["vulnerability_assessment_enabled"].(bool); vulnEnabled {
            config.VulnerabilityAssess.Enabled = pointer.To(true)
        }

        if advEnabled := raw["advanced_data_security_enabled"].(bool); advEnabled {
            config.AdvancedDataSecurity.Enabled = pointer.To(true)
        }
    }

    return config
}
```

**Hierarchical Enable Pattern:**
```go
// Parent-child enable relationships
"monitoring": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "enabled": {
                Type:     pluginsdk.TypeBool,
                Required: true,
            },
            "metrics": {
                Type:     pluginsdk.TypeList,
                Optional: true,
                MaxItems: 1,
                Elem: &pluginsdk.Resource{
                    Schema: map[string]*pluginsdk.Schema{
                        "enabled": {
                            Type:     pluginsdk.TypeBool,
                            Optional: true,
                            Default:  false,
                        },
                        "retention_days": {
                            Type:         pluginsdk.TypeInt,
                            Optional:     true,
                            Default:      30,
                            ValidateFunc: validation.IntBetween(1, 365),
                        },
                    },
                },
            },
            "logs": {
                Type:     pluginsdk.TypeList,
                Optional: true,
                MaxItems: 1,
                Elem: &pluginsdk.Resource{
                    Schema: map[string]*pluginsdk.Schema{
                        "enabled": {
                            Type:     pluginsdk.TypeBool,
                            Optional: true,
                            Default:  false,
                        },
                        "categories": {
                            Type:     pluginsdk.TypeSet,
                            Optional: true,
                            Elem: &pluginsdk.Schema{
                                Type: pluginsdk.TypeString,
                                ValidateFunc: validation.StringInSlice([]string{
                                    "Audit",
                                    "Performance",
                                    "Security",
                                }, false),
                            },
                        },
                    },
                },
            },
        },
    },
},
```

**Hierarchical Expand Function:**
```go
func ExpandMonitoringConfiguration(input []interface{}) *azureapi.MonitoringConfiguration {
    if len(input) == 0 || input[0] == nil {
        return &azureapi.MonitoringConfiguration{
            Enabled: pointer.To(false),
            Metrics: &azureapi.MetricsConfig{Enabled: pointer.To(false)},
            Logs:    &azureapi.LogsConfig{Enabled: pointer.To(false)},
        }
    }

    raw := input[0].(map[string]interface{})

    config := &azureapi.MonitoringConfiguration{
        Enabled: pointer.To(raw["enabled"].(bool)),
        Metrics: &azureapi.MetricsConfig{Enabled: pointer.To(false)},
        Logs:    &azureapi.LogsConfig{Enabled: pointer.To(false)},
    }

    // Only process child configurations if parent is enabled
    if config.Enabled != nil && *config.Enabled {
        if metricsRaw, exists := raw["metrics"]; exists {
            metricsList := metricsRaw.([]interface{})
            if len(metricsList) > 0 && metricsList[0] != nil {
                metrics := metricsList[0].(map[string]interface{})
                config.Metrics = &azureapi.MetricsConfig{
                    Enabled:       pointer.To(metrics["enabled"].(bool)),
                    RetentionDays: pointer.To(int32(metrics["retention_days"].(int))),
                }
            }
        }

        if logsRaw, exists := raw["logs"]; exists {
            logsList := logsRaw.([]interface{})
            if len(logsList) > 0 && logsList[0] != nil {
                logs := logsList[0].(map[string]interface{})
                config.Logs = &azureapi.LogsConfig{
                    Enabled:    pointer.To(logs["enabled"].(bool)),
                    Categories: expandStringSet(logs["categories"].(*pluginsdk.Set).List()),
                }
            }
        }
    }

    return config
}
```

### FivePointOh Feature Flag Patterns

**Breaking Changes and Deprecation Management:**

The `FivePointOh` feature flag system allows controlled introduction of breaking changes during development of major provider versions. This system is essential for managing deprecations and breaking changes in a controlled way.

**Key Functions:**
```go
import "github.com/hashicorp/terraform-provider-azurerm/internal/features"

// Check if running in 5.0 mode
if features.FivePointOh() {
    // Breaking change behavior for 5.0
}
```

**Untyped Resource Schema Deprecation Pattern:**
```go
// Schema with conditional deprecated fields (untyped resources)
func resourceServiceSchema() map[string]*pluginsdk.Schema {
    schema := map[string]*pluginsdk.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "new_field": {
            Type:     pluginsdk.TypeString,
            Optional: true,
        },
    }

    // Add deprecated fields conditionally - placed after main schema
    if !features.FivePointOh() {
        schema["legacy_field"] = &pluginsdk.Schema{
            Type:       pluginsdk.TypeString,
            Optional:   true,
            Deprecated: "This field will be removed in v5.0. Use `new_field` instead.",
        }
    }

    return schema
}
```

**Typed Resource Schema Deprecation Pattern:**
```go
// Schema with conditional deprecated fields (typed resources)
func (r ServiceNameResource) Arguments() map[string]*pluginsdk.Schema {
    arguments := map[string]*pluginsdk.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "new_field": {
            Type:     pluginsdk.TypeString,
            Optional: true,
        },
    }

    // Add deprecated fields conditionally - placed after main schema
    if !features.FivePointOh() {
        arguments["legacy_field"] = &pluginsdk.Schema{
            Type:       pluginsdk.TypeString,
            Optional:   true,
            Deprecated: "This field will be removed in v5.0. Use `new_field` instead.",
        }
    }

    return arguments
}
```

**Untyped Resource Implementation Patterns:**
```go
// Conditional behavior in resource operations (untyped)
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    properties := &serviceapi.Properties{
        Name: d.Get("name").(string),
    }

    // Use new field in 5.0, legacy field in 4.x
    if features.FivePointOh() {
        if newField := d.Get("new_field").(string); newField != "" {
            properties.NewProperty = pointer.To(newField)
        }
    } else {
        if legacyField := d.Get("legacy_field").(string); legacyField != "" {
            properties.LegacyProperty = pointer.To(legacyField)
        }
    }

    return nil
}

// Conditional state management (untyped)
func resourceServiceNameRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    resp, err := client.Get(ctx, id)
    if err != nil {
        return err
    }

    d.Set("name", resp.Properties.Name)
    d.Set("new_field", pointer.From(resp.Properties.NewProperty))

    // Only set legacy field in 4.x mode
    if !features.FivePointOh() {
        d.Set("legacy_field", pointer.From(resp.Properties.LegacyProperty))
    }

    return nil
}
```

**Typed Resource Implementation Patterns:**
```go
// Typed resource model with conditional fields
type ServiceNameResourceModel struct {
    Name      string `tfschema:"name"`
    NewField  string `tfschema:"new_field"`

    // Legacy field handled conditionally in decode/encode
    LegacyField string `tfschema:"legacy_field,removedInNextMajorVersion"`
}

// Conditional behavior in typed resource operations
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceNameResourceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            properties := &serviceapi.Properties{
                Name: model.Name,
            }

            // Use new field in 5.0, legacy field in 4.x
            if features.FivePointOh() {
                if model.NewField != "" {
                    properties.NewProperty = pointer.To(model.NewField)
                }
            } else {
                if model.LegacyField != "" {
                    properties.LegacyProperty = pointer.To(model.LegacyField)
                }
            }

            return nil
        },
    }
}
```

**Testing with FivePointOh:**
```go
// Test both 4.x and 5.0 behaviors
func TestAccResource_deprecationBehavior(t *testing.T) {
    // Skip test when using deprecated functionality in 5.0 mode
    if features.FivePointOh() {
        t.Skip("Test skipped in 5.0 mode: `legacy_field` was deprecated and removed in v5.0")
    }

    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    // Test current behavior (4.x mode)
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.withLegacyField(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("legacy_field").HasValue("test"),
            ),
        },
        data.ImportStep(),
    })
}

// Test configuration using legacy patterns
func (r ResourceResource) withLegacyField(data acceptance.TestData) string {
    return fmt.Sprintf(`
resource "azurerm_resource" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  # Legacy field - will show deprecation warning in 5.0 mode
  legacy_field = "test"
}
`, data.RandomInteger)
}
```

**Important Considerations:**
- **Environment Variable**: The `ARM_FIVEPOINTZERO_BETA` environment variable enables 5.0 mode for testing
- **Development Only**: This flag is for development and testing - NOT for production use
- **State Changes**: Enabling 5.0 mode may cause irreversible state changes
- **Documentation**: Always document the migration path in breaking change guides

**📚 Official Breaking Change Reference:**
For authoritative breaking change procedures, see: [Contributing Guide - Breaking Changes](../../../contributing/topics/guide-breaking-changes.md)

### Advanced Schema Validation Patterns

**Multi-Field Dependency Validation:**
```go
// CustomizeDiff for complex field relationships
func validateAdvancedConfiguration(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
    skuName := diff.Get("sku_name").(string)

    // Validate advanced security is only available on Premium SKUs
    if advSecurityRaw := diff.Get("advanced_security").([]interface{}); len(advSecurityRaw) > 0 {
        if skuName != "Premium" && skuName != "BusinessCritical" {
            return fmt.Errorf("`advanced_security` is only available with Premium or BusinessCritical SKUs")
        }

        advSecurity := advSecurityRaw[0].(map[string]interface{})

        // Validate threat detection requires vulnerability assessment
        if advSecurity["threat_detection_enabled"].(bool) && !advSecurity["vulnerability_assessment_enabled"].(bool) {
            return fmt.Errorf("`vulnerability_assessment_enabled` must be true when `threat_detection_enabled` is true")
        }
    }

    // Validate monitoring configuration hierarchy
    if monitoringRaw := diff.Get("monitoring").([]interface{}); len(monitoringRaw) > 0 {
        monitoring := monitoringRaw[0].(map[string]interface{})

        if !monitoring["enabled"].(bool) {
            // If monitoring is disabled, child features should also be disabled
            if metricsRaw := monitoring["metrics"].([]interface{}); len(metricsRaw) > 0 {
                metrics := metricsRaw[0].(map[string]interface{})
                if metrics["enabled"].(bool) {
                    return fmt.Errorf("`metrics.enabled` cannot be true when `monitoring.enabled` is false")
                }
            }
        }
    }

    return nil
}
```
---
[⬆️ Back to top](#schema-design-patterns)

## 🏗️ Complex Schema Patterns

### Advanced Schema Patterns

```go
// Single configuration block
"configuration": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: configurationSchema(),
    },
},

// Limited multiple items with Azure service constraints
"log_scrubbing_rule": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    MaxItems: 3, // Azure service limitation
    Elem: &pluginsdk.Resource{
        Schema: scrubbingRuleSchema(),
    },
},

// Unlimited items (don't specify MaxItems)
"network_interface": {
    Type:     pluginsdk.TypeSet,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: networkInterfaceSchema(),
    },
},
```

### TypeSet vs TypeList Guidelines

**Use TypeSet when:**
- Order doesn't matter
- Duplicates should be prevented
- Items are independent configurations

```go
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
```

**Use TypeList when:**
- Order matters
- Duplicates are allowed/meaningful
- Sequential processing is required

```go
"custom_rule": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    Elem: &pluginsdk.Resource{
        Schema: customRuleSchema(),
    },
},
```

### Sensitive Field Patterns

```go
// Sensitive strings (passwords, keys, secrets)
"password": {
    Type:      pluginsdk.TypeString,
    Required:  true,
    Sensitive: true,
},

// Sensitive in state but not in logs
"connection_string": {
    Type:      pluginsdk.TypeString,
    Computed:  true,
    Sensitive: true,
},
```
---
[⬆️ Back to top](#schema-design-patterns)

## 🔍 Field Naming Standards

### Descriptive Field Names

```go
// PREFERRED - Descriptive and clear
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

// AVOID - Generic or ambiguous names
"scrubbing_rule": {  // Too generic - what kind of scrubbing?
    Type:     pluginsdk.TypeSet,
    Optional: true,
},
"rule": {  // Too vague - what kind of rule?
    Type:     pluginsdk.TypeSet,
    Optional: true,
},
```

### Consistency Across Resources

```go
// Consistent naming for similar functionality
"response_timeout_seconds": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    ValidateFunc: validation.IntBetween(16, 240),
},

// Use same pattern across resources
"connection_timeout_seconds": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    ValidateFunc: validation.IntBetween(1, 300),
},
```

### Field Ordering Standards

**Arguments Section Ordering:**
- Required fields first (alphabetical within category)
- Optional fields second (alphabetical within category)
- Consistent structure across all resources

**Attributes Section Ordering:**
- Standard Azure fields first (`id`, `location`, etc.)
- Service-specific computed fields
- Complex nested computed structures

```go
func (r ServiceResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        // Required fields first (alphabetical)
        "location": commonschema.Location(),
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },
        "resource_group_name": commonschema.ResourceGroupName(),

        // Optional fields second (alphabetical)
        "enabled": {
            Type:     pluginsdk.TypeBool,
            Optional: true,
            Default:  true,
        },
        "tags": tags.Schema(),
        "timeout_seconds": {
            Type:         pluginsdk.TypeInt,
            Optional:     true,
            ValidateFunc: validation.IntBetween(1, 3600),
        },
    }
}
```

---
[⬆️ Back to top](#schema-design-patterns)

## 🧪 Testing Schema Patterns

### Advanced Configuration Testing

```go
func TestAccResource_advancedConfiguration(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.advancedConfigurationComplete(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("monitoring.0.enabled").HasValue("true"),
                check.That(data.ResourceName).Key("monitoring.0.metrics.0.enabled").HasValue("true"),
                check.That(data.ResourceName).Key("monitoring.0.logs.0.enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
        {
            Config: r.advancedConfigurationDisabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("monitoring.0.enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
    })
}
```

### Field Removal ForceNew Testing

```go
func TestAccResource_advancedConfiguration_fieldRemovalForceNew(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.withIrreversibleFeature(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("irreversible_feature_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
        {
            Config:             r.withoutIrreversibleFeature(data),
            ExpectNonEmptyPlan: true,
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("irreversible_feature_enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
    })
}
```

### Optional+Computed Field Testing

```go
func TestAccResource_optionalComputedDefaultRestoration(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.withExplicitOCValues(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("cookie_expiration_in_minutes").HasValue("60"),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
        {
            Config: r.withoutOCFields(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("cookie_expiration_in_minutes").HasValue("30"),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
        {
            Config: r.withExplicitOCValues(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("cookie_expiration_in_minutes").HasValue("60"),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
    })
}
```

## 🔧 Test Configuration Helpers

### Advanced Pattern Helpers

```go
func (r ResourceResource) withExplicitOCValues(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource" "test" {
  name                             = "acctest-%d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  sku_name                        = "Premium"

  # Explicit O+C values
  cookie_expiration_in_minutes    = 60
  resilient_vm_creation_enabled   = true

  encryption_settings {
    enabled          = true
    key_vault_key_id = azurerm_key_vault_key.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ResourceResource) withoutOCFields(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name           = "Premium"

  # O+C fields removed - should show Azure service defaults in state
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ResourceResource) withIrreversibleFeature(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource" "test" {
  name                           = "acctest-%d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  sku_name                      = "Premium"

  # Irreversible feature that cannot be disabled once enabled
  irreversible_feature_enabled  = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ResourceResource) withoutIrreversibleFeature(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource" "test" {
  name                = "acctest-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name           = "Premium"

  # irreversible_feature_enabled field removed - should trigger ForceNew
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
```

## 📚 Related Implementation Guidance (On-Demand)

### **Core Implementation**
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md) - Complete implementation patterns with schema integration
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md) - Azure-specific schema behaviors and the "None" pattern

### **Testing & Validation**
- 🧪 **Testing Guidelines**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md) - Schema validation testing patterns
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md) - Schema validation error handling

### **Documentation & Standards**
- 📝 **Documentation Guidelines**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md) - Documenting schema fields and validation
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md) - Azure provider schema standards

---
[⬆️ Back to top](#schema-design-patterns)
