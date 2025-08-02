---
applyTo: "internal/**/*.go"
description: Schema design patterns and validation standards for the Terraform AzureRM provider including field types, validation patterns, and Azure-specific schema considerations.
---

# Schema Design Patterns

Schema design patterns and validation standards for the Terraform AzureRM provider including field types, validation patterns, and Azure-specific schema considerations.

**Quick navigation:** [📋 Schema Types](#📋-schema-type-patterns) | [✅ Validation](#✅-validation-patterns) | [⚙️ Azure Specific](#⚙️-azure-specific-schema-patterns) | [🏗️ Complex Schemas](#🏗️-complex-schema-patterns)

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

### Field Removal ForceNew Pattern

**Critical Pattern for Fields Removed from Configuration Requiring Resource Recreation:**

When Azure resources have irreversible configuration changes (like enabling security policies that cannot be disabled), removing the field from Terraform configuration should trigger resource recreation. This requires using `CustomizeDiffShim` with both `SetNew()` and `ForceNew()` to work together.

**Why Both SetNew and ForceNew Are Required:**
- **SetNew()**: Creates a detectable state change in Terraform's plan showing the field going from `true` → `false`
- **ForceNew()**: Triggers resource recreation when this change occurs
- **Plan Visibility**: Terraform must show the field value change to justify the ForceNew action to users
- **Test Framework**: Acceptance tests require visible state changes to validate ForceNew behavior

**Implementation Pattern:**
```go
pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
    var featureExists, policyExists bool

    // Check if fields exist in the raw configuration (not computed/inferred values)
    if rawConfig := diff.GetRawConfig(); !rawConfig.IsNull() {
        featureExists = !rawConfig.AsValueMap()["irreversible_feature_enabled"].IsNull()
        policyExists = !rawConfig.AsValueMap()["security_policy_enabled"].IsNull()
    }

    // Only apply ForceNew logic during updates (not during initial creation)
    if diff.Id() != "" {
        // Handle irreversible_feature_enabled field removal
        if !featureExists {
            // Check if field was previously enabled in state
            if old, _ := diff.GetChange("irreversible_feature_enabled"); old.(bool) {
                // CRITICAL: SetNew makes the change visible in Terraform plan
                // This shows users: irreversible_feature_enabled: true → false
                if err := diff.SetNew("irreversible_feature_enabled", false); err != nil {
                    return fmt.Errorf("setting `irreversible_feature_enabled` to `false`: %+v", err)
                }
                // ForceNew triggers resource recreation since Azure cannot disable this feature
                return diff.ForceNew("irreversible_feature_enabled")
            }
        }

        // Handle security_policy_enabled field removal (same pattern)
        if !policyExists {
            if old, _ := diff.GetChange("security_policy_enabled"); old.(bool) {
                // Same pattern: make change visible then force recreation
                if err := diff.SetNew("security_policy_enabled", false); err != nil {
                    return fmt.Errorf("setting `security_policy_enabled` to `false`: %+v", err)
                }
                return diff.ForceNew("security_policy_enabled")
            }
        }
    }

    return nil
}),
```

**Azure Use Cases:**
- **VM Scale Set Resiliency Policies**: Cannot be disabled once enabled
- **Security Features**: Irreversible security configurations
- **Compliance Settings**: Audit policies that cannot be downgraded
- **Performance Tiers**: Service levels that require recreation to reduce

**Key Requirements:**
- **Irreversible Changes**: Only use for Azure features that cannot be disabled once enabled
- **Raw Config Detection**: Use `GetRawConfig().AsValueMap()` to detect field presence vs absence in configuration
- **Update-Only Logic**: Check `diff.Id() != ""` to ensure logic only applies to existing resources, not during creation
- **State Visibility**: SetNew must be called before ForceNew to create visible plan entry
- **Error Handling**: SetNew errors should be caught and wrapped with descriptive context
- **Test Validation**: Tests must verify both the state change and ForceNew trigger

**Common Mistakes to Avoid:**
- **ForceNew without SetNew**: Plan won't show why recreation is needed - users will be confused by ForceNew without visible changes
- **SetNew without ForceNew**: State changes but resource doesn't recreate when Azure constraints require it
- **Missing Error Handling**: SetNew failures can break plan generation if not properly handled
- **Wrong Field Detection**: Use `GetRawConfig().AsValueMap()[field].IsNull()` to detect field removal, not `diff.Get()`
- **Creation vs Update**: Apply logic only during updates (`diff.Id() != ""`), not during initial resource creation

### Optional+Computed Field Patterns

**Modern Optional+Computed Implementation for Azure Service Defaults:**

Optional+Computed (O+C) fields handle Azure service defaults gracefully while maintaining user control. This pattern is essential for Azure resources where the service provides intelligent defaults that users may want to override.

**Schema Definition Pattern:**
```go
"resilient_vm_creation_enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Computed: true,
},

"cookie_expiration_in_minutes": {
    Type:         pluginsdk.TypeInt,
    Optional:     true,
    Computed:     true,
    ValidateFunc: validation.IntBetween(1, 43200), // 1 minute to 30 days
},

"encryption_settings": {
    Type:     pluginsdk.TypeList,
    Optional: true,
    Computed: true,
    MaxItems: 1,
    Elem: &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            "enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
                Computed: true,
            },
            "key_vault_key_id": {
                Type:         pluginsdk.TypeString,
                Optional:     true,
                ValidateFunc: keyVaultValidate.NestedItemId,
            },
        },
    },
},
```

**Three-Function Implementation Pattern for O+C Fields:**

**1. Read Function - Simple Azure Value Reading:**
```go
func (r ServiceResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // ...existing Azure API call logic...

            state := ServiceResourceModel{}

            // Simple value assignment - no config presence logic needed
            if props := resp.Model.Properties; props != nil {
                state.ResilientVmCreationEnabled = pointer.FromBool(props.ResilientVmCreationEnabled, false)
                state.CookieExpirationInMinutes = pointer.FromInt32(props.CookieExpirationInMinutes, 30)

                if encryptionSettings := props.EncryptionSettings; encryptionSettings != nil {
                    state.EncryptionSettings = []EncryptionSettingsModel{{
                        Enabled:       pointer.FromBool(encryptionSettings.Enabled, false),
                        KeyVaultKeyId: pointer.FromString(encryptionSettings.KeyVaultKeyId, ""),
                    }}
                }
            }

            return metadata.Encode(&state)
        },
    }
}
```

**2. Update Function - Default Application Using GetOk() Pattern:**
```go
func (r ServiceResource) Update() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceResourceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Read existing resource to get current values
            existing, err := client.Get(ctx, id)
            if err != nil {
                return fmt.Errorf("retrieving existing %s: %+v", id, err)
            }

            properties := existing.Model.Properties

            // Use GetOk() to detect field presence vs absence in configuration
            d := metadata.ResourceData

            // Apply user values when configured, preserve existing when removed from config
            if v, ok := d.GetOk("resilient_vm_creation_enabled"); ok {
                properties.ResilientVmCreationEnabled = pointer.To(v.(bool))
            } else {
                // Field removed from config - apply Azure service default
                properties.ResilientVmCreationEnabled = pointer.To(false) // Azure default
            }

            if v, ok := d.GetOk("cookie_expiration_in_minutes"); ok {
                properties.CookieExpirationInMinutes = pointer.To(int32(v.(int)))
            } else {
                // Field removed from config - apply Azure service default
                properties.CookieExpirationInMinutes = pointer.To(int32(30)) // Azure default
            }

            if encryptionRaw, ok := d.GetOk("encryption_settings"); ok {
                encryptionList := encryptionRaw.([]interface{})
                if len(encryptionList) > 0 && encryptionList[0] != nil {
                    encryption := encryptionList[0].(map[string]interface{})
                    properties.EncryptionSettings = &azureapi.EncryptionSettings{
                        Enabled:       pointer.To(encryption["enabled"].(bool)),
                        KeyVaultKeyId: pointer.To(encryption["key_vault_key_id"].(string)),
                    }
                }
            } else {
                // Block removed from config - apply Azure service defaults
                properties.EncryptionSettings = &azureapi.EncryptionSettings{
                    Enabled:       pointer.To(false), // Azure default
                    KeyVaultKeyId: nil,               // Azure default
                }
            }

            // Update the resource with new properties
            if err := client.UpdateThenPoll(ctx, id, azureapi.UpdateParameters{
                Properties: properties,
            }); err != nil {
                return fmt.Errorf("updating %s: %+v", id, err)
            }

            return nil
        },
    }
}
```

**3. CustomizeDiff Function - O+C Persistence Handling:**
```go
func (r ServiceResource) CustomizeDiff() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            d := metadata.ResourceData

            // Handle O+C field persistence when removed from configuration
            if d.Id() != "" { // Only for existing resources
                // Check if O+C fields were removed from configuration
                if !d.GetRawConfig().GetAttr("resilient_vm_creation_enabled").IsNull() {
                    // Field is present in config - normal validation applies
                } else {
                    // Field removed from config - ensure it persists in state with default
                    oldVal, _ := d.GetChange("resilient_vm_creation_enabled")
                    if oldVal != nil {
                        // Keep the field in state showing the default value
                        d.SetNewComputed("resilient_vm_creation_enabled")
                    }
                }

                // Same pattern for other O+C fields
                if !d.GetRawConfig().GetAttr("cookie_expiration_in_minutes").IsNull() {
                    // Field is present in config
                } else {
                    // Field removed from config
                    oldVal, _ := d.GetChange("cookie_expiration_in_minutes")
                    if oldVal != nil {
                        d.SetNewComputed("cookie_expiration_in_minutes")
                    }
                }

                if !d.GetRawConfig().GetAttr("encryption_settings").IsNull() {
                    // Block is present in config
                } else {
                    // Block removed from config
                    oldVal, _ := d.GetChange("encryption_settings")
                    if oldVal != nil {
                        d.SetNewComputed("encryption_settings")
                    }
                }
            }

            return nil
        },
    }
}
```

**Key Principles for O+C Implementation:**
- **Read Function**: Simple Azure value reading, no configuration presence logic
- **Update Function**: Use `GetOk()` to detect field presence and apply defaults when removed
- **CustomizeDiff**: Ensure O+C fields persist in state when removed from configuration
- **State Presence**: O+C fields remain in state forever once set, showing current values
- **Default Restoration**: When fields are removed from config, they show Azure service defaults in state
- **User Override**: When fields are re-added to config, user values take precedence

**Advanced O+C Testing Requirements:**
- **Initial State Testing**: Verify defaults appear when fields not configured
- **Explicit Configuration Testing**: Verify user values override defaults
- **Default Restoration Testing**: Verify defaults appear when fields removed from config
- **Import State Testing**: Verify import works with both user-set and default values
- **State Persistence Testing**: Verify O+C fields never disappear from state once set

### Advanced Testing Patterns for Complex Schemas

**Comprehensive Schema Testing:**
```go
func TestAccResource_advancedConfiguration(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.advancedConfigurationComplete(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // Test hierarchical enable patterns
                check.That(data.ResourceName).Key("monitoring.0.enabled").HasValue("true"),
                check.That(data.ResourceName).Key("monitoring.0.metrics.0.enabled").HasValue("true"),
                check.That(data.ResourceName).Key("monitoring.0.logs.0.enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
        {
            // Test disabling parent disables children
            Config: r.advancedConfigurationDisabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("monitoring.0.enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
    })
}

// Test SKU dependency validation
func TestAccResource_advancedConfiguration_skuValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.advancedConfigurationInvalidSku(data),
            ExpectError: regexp.MustCompile("`advanced_security` is only available with Premium or BusinessCritical SKUs"),
        },
    })
}

// Test field removal ForceNew behavior
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
            // Removing irreversible feature should force new resource
            Config:             r.withoutIrreversibleFeature(data),
            ExpectNonEmptyPlan: true, // ForceNew creates a plan
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // After ForceNew, the field should show false in new resource
                check.That(data.ResourceName).Key("irreversible_feature_enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
    })
}

// Test O+C field default restoration
func TestAccResource_optionalComputedDefaultRestoration(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_resource", "test")
    r := ResourceResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            // Step 1: Create with explicit values
            Config: r.withExplicitOCValues(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("cookie_expiration_in_minutes").HasValue("60"),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("true"),
            ),
        },
        data.ImportStep(),
        {
            // Step 2: Remove from config (should show defaults)
            Config: r.withoutOCFields(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // O+C fields should show Azure service defaults when removed from config
                check.That(data.ResourceName).Key("cookie_expiration_in_minutes").HasValue("30"),
                check.That(data.ResourceName).Key("resilient_vm_creation_enabled").HasValue("false"),
            ),
        },
        data.ImportStep(),
        {
            // Step 3: Re-add explicit values
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

**Test Configuration Helpers for Advanced Patterns:**
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

---

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 📋 **Code Clarity Enforcement**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)

### 🚀 Enhanced Guidance Files

- 🔄 **API Evolution**: [api-evolution-patterns.instructions.md](./api-evolution-patterns.instructions.md)
- ⚡ **Performance**: [performance-optimization.instructions.md](./performance-optimization.instructions.md)
- 🔐 **Security**: [security-compliance.instructions.md](./security-compliance.instructions.md)
- 🔧 **Troubleshooting**: [troubleshooting-decision-trees.instructions.md](./troubleshooting-decision-trees.instructions.md)

---
[⬆️ Back to top](#schema-design-patterns)
