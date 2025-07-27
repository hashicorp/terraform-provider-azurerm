---
applyTo: "internal/**/*.go"
description: Schema design patterns and validation standards for the Terraform AzureRM provider including field types, validation patterns, and Azure-specific schema considerations.
---

# Schema Design Patterns

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

### Optional+Computed Schema Pattern

**When to Use Optional+Computed (O+C):**
- Azure-managed defaults that users can override
- Irreversible settings that cannot be disabled once enabled
- Service-managed state that users can influence but not fully control

**Documentation Requirements for O+C:**
```go
// This field is Optional+Computed for two reasons:
// 1. Once resilient VM creation policy is enabled (set to true), it cannot be disabled (reverted to false)
// 2. Backward compatibility - existing scale sets won't show diffs when upgrading the provider
// The Computed attribute ensures Terraform reflects the actual Azure state.
"resilient_vm_creation_enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Computed: true,
},
```

### ForceNew Patterns

```go
// Properties that require resource recreation
"location": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ForceNew: true, // Changing location requires new resource
},

"name": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ForceNew:     true, // Names cannot be changed
    ValidateFunc: validation.StringIsNotEmpty,
},

// Properties that support in-place updates
"enabled": {
    Type:     pluginsdk.TypeBool,
    Optional: true,
    Default:  true,
    // No ForceNew - can be updated in place
},
```

### SKU Validation Patterns

```go
// Use Azure SDK PossibleValues when available
"sku_name": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForProfileSku(),
        false,
    ),
},

// Manual validation when SDK doesn't provide function
"performance_tier": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        "Standard_S1",
        "Standard_S2",
        "Premium_P1",
        "Premium_P2",
    }, false),
},
```

---
[⬆️ Back to top](#schema-design-patterns)

## 🏗️ Complex Schema Patterns

### Schema Flattening Pattern

**When to Apply Schema Flattening:**
- Remove unnecessary wrapper structures from Azure APIs
- Simplify user experience without losing functionality
- Eliminate single-purpose wrapper blocks

**Before Flattening (Complex):**
```go
"log_scrubbing": {
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
            "scrubbing_rule": {
                Type:     pluginsdk.TypeSet,
                MaxItems: 3,
                Optional: true,
                Elem: &pluginsdk.Resource{
                    Schema: map[string]*pluginsdk.Schema{
                        "match_variable": {
                            Type:     pluginsdk.TypeString,
                            Required: true,
                            ValidateFunc: validation.StringInSlice([]string{
                                "QueryStringArgNames",
                                "RequestIPAddress",
                                "RequestUri",
                            }, false),
                        },
                        "enabled": {
                            Type:     pluginsdk.TypeBool,
                            Optional: true,
                            Default:  true,
                        },
                        "operator": {
                            Type:     pluginsdk.TypeString,
                            Optional: true,
                            Default:  "EqualsAny",
                            ValidateFunc: validation.StringInSlice([]string{
                                "EqualsAny",
                                "Equals",
                            }, false),
                        },
                        "selector": {
                            Type:     pluginsdk.TypeString,
                            Optional: true,
                        },
                    },
                },
            },
        },
    },
},
```

**After Flattening (Simplified):**
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

### Conditional Field Requirements

**Use Notes for Complex Dependencies:**
```go
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
        false),
},

// ~> **Note:** The `match_variable` determines what type of request data will be scrubbed from logs.
```

### MaxItems Usage Patterns

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

- 🏗️ **Main Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- ⚡ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 🚨 **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
