---
applyTo: "internal/**/*.go"
description: This document outlines the coding standards for Go files in the Terraform AzureRM provider repository. It includes naming conventions, file organization, error handling patterns, resource implementation guidelines, and Azure SDK usage standards.
---

## Coding Standards
Given below are the coding standards for the Terraform AzureRM provider which **MUST** be followed.

### Typed vs Untyped Resource Implementation Standards

#### Typed Resource Implementation (Preferred)
The Typed Resource implementation uses the `internal/sdk` framework and should be used for new resources and data sources.

#### UnTyped Resource Implementation
The Untyped Resource Implementation uses traditional Plugin SDK patterns and should be maintained for existing resources but not used for new implementations.

### Naming Conventions

#### Package Names
- Use lowercase, single-word package names
- Match the service name (e.g., `compute`, `storage`, `network`)
- Avoid underscores or mixed caps

#### Function Names
- **Exported functions**: PascalCase (e.g., `CreateResource`, `ValidateInput`)
- **Unexported functions**: camelCase (e.g., `parseResourceID`, `buildParameters`)

**Typed Resource Implementation:**
- **Resource struct methods**: Use receiver methods on struct types
  - Examples: `(r ServiceNameResource) Create()`, `(r ServiceNameResource) Read()`
- **Data source struct methods**: Use receiver methods on struct types
  - Examples: `(r ServiceNameDataSource) Read()`
- **Model structs**: Use PascalCase with descriptive suffixes
  - Examples: `ServiceNameModel`, `ServiceNameDataSourceModel`

**UnTyped Resource Implementation:**
- **Resource CRUD functions**: `resource[ResourceType][Operation]`
  - Examples: `resourceVirtualMachineCreate`, `resourceStorageAccountRead`
- **Data source functions**: `dataSource[ResourceType]`
  - Examples: `dataSourceVirtualMachine`, `dataSourceResourceGroup`

#### Variable Names
- **Exported variables**: PascalCase
- **Unexported variables**: camelCase
- **Constants**: PascalCase for exported, camelCase for unexported
- **Acronyms**: Keep uppercase (e.g., `resourceGroupID`, `vmSSH`, `apiURL`)

#### Type Names
- **Exported types**: PascalCase
- **Unexported types**: camelCase
- **Interface names**: Often descriptive or end with 'er' (e.g., `ResourceProvider`, `Validator`)

### File Organization

#### Directory Structure
- **Resource files**: `internal/services/[service]/[resource_type]_resource.go`
- **Resource Test files**: Same directory and name as source with `_test.go` suffix
- **Resource documentation files**: `internal/website/docs/r/[resource_type].html.markdown`
- **Data source files**: `internal/services/[service]/[resource_type]_data_source.go`
- **Data source Test files**: Same directory and name as data source with `_test.go` suffix
- **Data source documentation files**: `internal/website/docs/d/[resource_type].html.markdown`
- **Utility files**: Group related functions (e.g., `validate.go`, `parse.go`, `flatten.go`, `expand.go`)
- **Registration**: Each service has a `registration.go` file

#### File Naming
- Use snake_case for file names
- Keep files focused on single responsibility
- Aim for files under 1000 lines when possible
- Separate complex logic into utility functions

### Resource Implementation Patterns

#### Typed Resource Implementation (Preferred)

**Model Structure:**
```go
package servicename

import (
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/commonschema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tags"
)

type ServiceNameModel struct {
    Name              string            `tfschema:"name"`
    ResourceGroup     string            `tfschema:"resource_group_name"`
    Location          string            `tfschema:"location"`
    Sku               string            `tfschema:"sku_name"`
    Enabled           bool              `tfschema:"enabled"`
    Tags              map[string]string `tfschema:"tags"`
    
    // Computed attributes
    Endpoint          string            `tfschema:"endpoint"`
    Status            string            `tfschema:"status"`
}

type ServiceNameResource struct{}

var (
    _ sdk.Resource           = ServiceNameResource{}
    _ sdk.ResourceWithUpdate = ServiceNameResource{}
)

func (r ServiceNameResource) ResourceType() string {
    return "azurerm_service_name"
}

func (r ServiceNameResource) ModelObject() interface{} {
    return &ServiceNameModel{}
}

func (r ServiceNameResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateServiceNameID
}

func (r ServiceNameResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Description:  "The name of the Service.",
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },
        "resource_group_name": commonschema.ResourceGroupName(),
        "location": commonschema.Location(),
        "tags": tags.Schema(),
    }
}

func (r ServiceNameResource) Attributes() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "endpoint": {
            Description: "The endpoint URL of the Service.",
            Type:        pluginsdk.TypeString,
            Computed:    true,
        },
        "status": {
            Description: "The current status of the Service.",
            Type:        pluginsdk.TypeString,
            Computed:    true,
        },
    }
}
```

**Typed CRUD Functions:**
```go
package servicename

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            id := parse.NewServiceNameID(subscriptionId, model.ResourceGroup, model.Name)

            metadata.Logger.Infof("Import check for %s", id)
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
            }

            if !response.WasNotFound(existing.HttpResponse) {
                return metadata.ResourceRequiresImport(r.ResourceType(), id)
            }

            metadata.Logger.Infof("Creating %s", id)
            
            properties := servicenametype.Resource{
                Location: model.Location,
                Properties: &servicenametype.ResourceProperties{
                    Enabled: pointer.To(model.Enabled),
                },
                Tags: &model.Tags,
            }

            if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
                return fmt.Errorf("creating %s: %+v", id, err)
            }

            metadata.SetID(id)
            return nil
        },
    }
}

func (r ServiceNameResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient

            id, err := parse.ServiceNameID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }

            metadata.Logger.Infof("Reading %s", id)
            resp, err := client.Get(ctx, *id)
            if err != nil {
                if response.WasNotFound(resp.HttpResponse) {
                    return metadata.MarkAsGone(id)
                }
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }

            model := resp.Model
            if model == nil {
                return fmt.Errorf("retrieving %s: model was nil", id)
            }

            state := ServiceNameModel{
                Name:          id.ServiceName,
                ResourceGroup: id.ResourceGroupName,
                Location:      model.Location,
                Sku:           string(model.Sku.Name),
                Tags:          model.Tags,
            }

            if props := model.Properties; props != nil {
                state.Enabled = pointer.FromBool(props.Enabled, false)
                state.Endpoint = pointer.FromString(props.Endpoint, "")
                state.Status = pointer.FromString(props.Status, "")
            }

            return metadata.Encode(&state)
        },
    }
}
```

#### UnTyped Resource Implementation

**Standard CRUD Functions:**
```go
package servicename

import (
    "context"

    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Implementation here
}
func resourceServiceNameRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Implementation here
}  
func resourceServiceNameUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Implementation here
}
func resourceServiceNameDelete(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Implementation here
}
```

**Resource Schema Patterns:**
```go
package servicename

import (
    "time"

    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
)

func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,
        
        Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
            _, err := parse.ServiceNameID(id)
            return err
        }),
        
        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },
        
        Schema: map[string]*pluginsdk.Schema{
            "name": {
                Type:         pluginsdk.TypeString,
                Required:     true,
                ForceNew:     true,
                ValidateFunc: validation.StringIsNotEmpty,
            },
            // ... more schema definitions
        },
    }
}
```

#### ValidateFunc Patterns

If the Azure SDK package offers a `PossibleValuesForFieldName` function, use that in the `validation.StringInSlice` function instead of hardcoding the possible values manually.

##### Example
```go
// AVOID - Hardcoded values that may become outdated
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice([]string{
        string(profiles.ScrubbingRuleEntryMatchVariableQueryStringArgNames),
        string(profiles.ScrubbingRuleEntryMatchVariableRequestIPAddress),
        string(profiles.ScrubbingRuleEntryMatchVariableRequestUri),
    }, false),
},

// PREFERRED - Use SDK-provided possible values function
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
        false,
    ),
},
```

**Benefits of using SDK-provided functions:**

- **Automatic updates**: When Azure adds new values, they're automatically available
- **Consistency**: Ensures validation matches what the Azure API actually accepts
- **Maintenance**: Reduces manual updates when Azure service capabilities change
- **Accuracy**: Eliminates risk of typos in hardcoded values

### Error Handling

#### typed resource Error Patterns
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

// Use metadata.Decode for model decoding errors
var model ServiceNameModel
if err := metadata.Decode(&model); err != nil {
    return fmt.Errorf("decoding: %+v", err)
}

// Use metadata.Logger for structured logging
metadata.Logger.Infof("Import check for %s", id)

// Use metadata.ResourceRequiresImport for import conflicts
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}

// Use metadata.MarkAsGone for deleted resources
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// Use metadata.SetID for resource ID management
metadata.SetID(id)

// Use metadata.Encode for state management
return metadata.Encode(&model)
```

#### untyped Error Patterns
```go
package servicename

import (
    "fmt"
    "log"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/helpers/resource"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

// Use consistent error formatting with context
if err != nil {
    return fmt.Errorf("creating Resource %q: %+v", name, err)
}

// Include resource information in error messages
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Resource %q was not found - removing from state", id.ResourceName)
    d.SetId("")
    return nil
}

// Handle Azure-specific errors
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled"))
}
```


#### CustomizeDiff Implementation Pattern

**Dual Import Requirement:**
When implementing CustomizeDiff functions, both packages must be imported:

```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
```

**Standard CustomizeDiff Resource Pattern:**
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,

        CustomizeDiff: pluginsdk.All(
            // Must use *schema.ResourceDiff from external package
            pluginsdk.ForceNewIfChange("property_name", func(ctx context.Context, old, new, meta interface{}) bool {
                return old.(string) != new.(string)
            }),
            func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Custom validation logic
                if diff.Get("enabled").(bool) && diff.Get("configuration") == nil {
                    return fmt.Errorf("configuration is required when enabled is true")
                }
                return nil
            },
        ),

        Schema: map[string]*pluginsdk.Schema{
            // Schema definitions use pluginsdk types
        },
    }
}
```

**Why This Pattern is Required:**
- The internal pluginsdk package provides aliases for most Plugin SDK types
- However, CustomizeDiff function signatures are **not** aliased and must use *schema.ResourceDiff
- The pluginsdk.All(), pluginsdk.ForceNewIfChange() helpers are available in the internal package
- Resource and schema definitions use pluginsdk types for consistency

#### Common CustomizeDiff Import Issues and Troubleshooting

**Compilation Error: `pluginsdk.ResourceDiff` is not defined**
```go
package servicename

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//  INCORRECT - Will cause compilation error
func(ctx context.Context, diff *pluginsdk.ResourceDiff, meta interface{}) error {
    // pluginsdk.ResourceDiff does not exist - this will fail to compile
}

//  CORRECT - Must use external schema package
func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
    // Custom validation logic using the correct type
}
```

**Missing Import Error**
If you see errors like `undefined: schema.ResourceDiff`, ensure you have both imports:
```go
package servicename

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"    // Required for *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"  // Required for helpers
)
```

**Helper Function Usage**
```go
package servicename

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Use pluginsdk helpers for common patterns
CustomizeDiff: pluginsdk.All(
    pluginsdk.ForceNewIfChange("location", func(ctx context.Context, old, new, meta interface{}) bool {
        return old.(string) != new.(string)
    }),
    // Custom validation functions use *schema.ResourceDiff
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        // Validation logic
        return nil
    },
),
```

#### Azure-Specific CustomizeDiff Examples

**Azure Storage Account SKU and Kind Validation:**
```go
package storage

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

CustomizeDiff: pluginsdk.All(
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        accountTier := diff.Get("account_tier").(string)
        accountKind := diff.Get("account_kind").(string)
        
        if accountTier == "Premium" && accountKind != "BlockBlobStorage" && accountKind != "FileStorage" {
            return fmt.Errorf("`account_kind` must be `BlockBlobStorage` or `FileStorage` when `account_tier` is `Premium`")
        }
        return nil
    },
    pluginsdk.ForceNewIfChange("location", func(ctx context.Context, old, new, meta interface{}) bool {
        return old.(string) != new.(string)
    }),
),
```

**Azure Virtual Machine SKU and Zone Validation:**
```go
package compute

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

CustomizeDiff: pluginsdk.All(
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        size := diff.Get("size").(string)
        zones := diff.Get("zones").([]interface{})
        
        // Check if VM size supports availability zones
        if len(zones) > 0 && !supportsAvailabilityZones(size) {
            return fmt.Errorf("VM size `%s` does not support availability zones", size)
        }
        return nil
    },
    pluginsdk.ForceNewIfChange("size", func(ctx context.Context, old, new, meta interface{}) bool {
        // Force recreation when changing VM size
        return old.(string) != new.(string)
    }),
),
```

**Azure Premium SKU with Zone Redundancy:**
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

CustomizeDiff: pluginsdk.All(
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        skuName := diff.Get("sku_name").(string)
        zoneRedundant := diff.Get("zone_redundant").(bool)
        
        if skuName == "Premium" && !zoneRedundant {
            return fmt.Errorf("`zone_redundant` must be true for Premium SKU")
        }
        return nil
    },
),
```

**Azure CDN Front Door Log Scrubbing Validation:**
```go
package cdn

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

CustomizeDiff: pluginsdk.All(
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        // Validate log scrubbing configuration
        if scrubbingRulesRaw, ok := diff.GetOk("log_scrubbing.0.scrubbing_rules"); ok {
            scrubbingRules := scrubbingRulesRaw.([]interface{})
            
            for i, ruleRaw := range scrubbingRules {
                rule := ruleRaw.(map[string]interface{})
                matchVariable := rule["match_variable"].(string)
                selector := rule["selector"].(string)
                
                // Azure API constraint: selector is required for QueryStringArgNames
                if matchVariable == "QueryStringArgNames" {
                    if selector == "" {
                        return fmt.Errorf("log_scrubbing.0.scrubbing_rules.%d: `selector` is required when `match_variable` is `%s`", i, matchVariable)
                    }
                } else if matchVariable == "RequestIPAddress" || matchVariable == "RequestUri" {
                    // Azure API constraint: selector cannot be set for RequestIPAddress and RequestUri
                    if selector != "" {
                        return fmt.Errorf("log_scrubbing.0.scrubbing_rules.%d: `selector` cannot be set when `match_variable` is `%s`", i, matchVariable)
                    }
                }
            }
        }
        return nil
    },
),
```

### Migration Guidelines

#### Migrating from Untyped to Typed Implementation

**Migration Decision Criteria:**
- **New Resources**: Always use Typed Resource Implementation
- **Major Feature Additions**: Consider migration opportunity when adding significant functionality
- **Bug Fixes**: Maintain existing untyped implementation for simple fixes

**Service Registration During Migration:**
When migrating resources from untyped to typed implementations, services often need to be registered in both lists temporarily:

```go
// In internal/provider/services.go

func SupportedTypedServices() []sdk.TypedServiceRegistration {
    services := []sdk.TypedServiceRegistration{
        // Add service here when it has any typed resources
        cdn.Registration{},
        // ...other services
    }
    return services
}

func SupportedUntypedServices() []sdk.UntypedServiceRegistration {
    return func() []sdk.UntypedServiceRegistration {
        out := []sdk.UntypedServiceRegistration{
            // Keep service here until all resources are migrated
            cdn.Registration{},
            // ...other services
        }
        return out
    }()
}
```

**Registration Rules:**
- **Dual Registration**: Services appear in both lists during migration period
- **Service Registration Interface**: The same service registration struct implements both `TypedServiceRegistration` and `UntypedServiceRegistration` interfaces
- **Resource Distribution**: Individual resources within the service can be either typed or untyped
- **Final Migration**: Remove from `SupportedUntypedServices()` only when all resources in the service are migrated

**Registration Implementation Requirements:**
- **Both Functions Required**: When a service implements both `TypedServiceRegistration` and `UntypedServiceRegistration` interfaces, you **MUST** implement both sets of functions even if no resources exist in one category
- **Empty Returns Allowed**: Return empty slices/maps for resource types that don't exist yet
- **Consistent Interface**: This ensures the registration struct satisfies both interfaces correctly

```go
var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/connections"
}

// REQUIRED: Always implement both typed and untyped functions
func (r Registration) DataSources() []sdk.DataSource {
    return []sdk.DataSource{
        // Typed data sources here, or empty slice if none exist
        ApiConnectionDataSource{},
    }
}

func (r Registration) Resources() []sdk.Resource {
    return []sdk.Resource{
        // Typed resources here, or empty slice if none exist
        // Add typed resources here when implemented
    }
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Untyped data sources here, or empty map if none exist
        "azurerm_managed_api": dataSourceManagedApi(),
    }
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Untyped resources here, or empty map if none exist
        "azurerm_api_connection": resourceApiConnection(),
    }
}
```
**Migration Decision Matrix:**

| Scenario | Action | Approach |
|----------|--------|----------|
| New Resource | Always use Typed Resource Implementation | Start with typed from day one |
| Bug Fix (< 5 lines) | Maintain Untyped Implementation | Quick fix in existing pattern |
| Feature Addition (< 50 lines) | Consider migration if touching >30% of resource | Evaluate cost/benefit |
| Major Refactor (> 50 lines) | Migrate to Typed Implementation | Plan migration with comprehensive testing |
| EOL/Deprecation Planning | Plan Typed Migration | Include in deprecation timeline |

**Migration Process:**
1. **Assessment Phase**
   - Analyze existing resource complexity and usage patterns
   - Identify breaking changes that may be required
   - Plan migration timeline and testing strategy

2. **Implementation Phase**
   ```go
   package servicename

   import (
       "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
   )

   // Step 1: Create typed model structure
   type ExistingResourceModel struct {
       Name          string            `tfschema:"name"`
       ResourceGroup string            `tfschema:"resource_group_name"`
       // Map all existing schema fields
    }
   
   // Step 2: Implement typed resource interface
   type ExistingResource struct{}
   
   var _ sdk.Resource = ExistingResource{}
   ```

3. **Testing Phase**
   - Comprehensive acceptance test coverage
   - State compatibility validation
   - Import functionality verification
   - Cross-version compatibility testing

**Migration Considerations:**
- **State Compatibility**: Ensure Terraform state remains compatible during migration
- **User Impact**: Migration should be transparent to end users
- **Feature Parity**: Maintain all existing functionality
- **Documentation**: Update examples and guides without breaking existing configurations
- **Rollback Plan**: Maintain ability to revert if critical issues arise

**State Management During Migration:**
```go
package servicename

import (
    "context"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
)

// Ensure consistent resource ID format
func (r ExistingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateExistingResourceID  // Use same parser as untyped version
}

// Maintain backward compatibility in Read operation
func (r ExistingResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            
            // Handle both old and new state formats if necessary
            id, err := parse.ExistingResourceID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }
                        
            resp, err := client.Get(ctx, *id)
            if err != nil {
                if response.WasNotFound(resp.HttpResponse) {
                    return metadata.MarkAsGone(id)
                }
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }
            
            // Map response to model and encode state
            model := ExistingResourceModel{
                Name:          id.ResourceName,
                ResourceGroup: id.ResourceGroupName,
                // Map other fields from response
            }
            
            return metadata.Encode(&model)
        },
    }
}
```

### Azure SDK Integration

#### Pointer Package Usage

**Use the `pointer` package instead of the `utils` package for pointer operations where applicable:**  

##### Migration Criteria  
- **New Code**: Always use the `pointer` package for new implementations.  
- **Existing Code**: Migrate from `utils` to `pointer` if:  
  - The code is being actively modified or refactored.  
  - The migration does not introduce significant risk or require extensive changes.  
- **Legacy Compatibility**: Retain `utils` usage if:  
  - The code is stable and not undergoing changes.  
  - The migration would disrupt backward compatibility or require substantial effort.  

Examples:  
- **Migrate**: Updating a resource implementation to use the `pointer` package for optional Azure API parameters.  
- **Do Not Migrate**: Stable legacy code that uses `utils` for pointer operations and is not being actively maintained. 

```go
package servicename

import (
    "github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// PREFERRED - Use pointer package
// Convert values to pointers
// Common scenarios: optional Azure API parameters, nullable fields
stringPtr := pointer.To("example")
intPtr := pointer.To(42)
boolPtr := pointer.To(true)

// Convert pointers to values with defaults
stringValue := pointer.From(stringPtr)
stringValueWithDefault := pointer.FromString(stringPtr, "default")
intValue := pointer.FromInt32(intPtr, 0)
boolValue := pointer.FromBool(boolPtr, false)

// Check if pointer is nil or has value
if pointer.IsNil(stringPtr) {
    // Handle nil pointer
}
```

```go
package servicename

import (
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)

// AVOID - Legacy utils package patterns (where pointer package can be used)
// Legacy patterns - use pointer package instead
stringPtr := utils.String("example")  // Use pointer.To("example")
intPtr := utils.Int32(42)             // Use pointer.To(42)
boolPtr := utils.Bool(true)           // Use pointer.To(true)

// Legacy dereference patterns
if stringPtr != nil {
    value := *stringPtr
}
// Use pointer.From(stringPtr) or pointer.FromString(stringPtr, "default")
```

**When to Use Each Package:**
- **pointer package**: For basic pointer operations, type conversions, and nil checking
- **utils package**: For Azure-specific utilities, complex transformations, and legacy compatibility where pointer package does not provide equivalent functionality

**Migration Guidelines:**
- **New Code**: Always use `pointer` package for pointer operations
- **Existing Code**: Gradually migrate to `pointer` package during refactoring
- **Legacy Compatibility**: Maintain `utils` package usage only where `pointer` package does not provide equivalent functionality

#### typed resource Client Usage
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// Use metadata.Client for accessing clients
client := metadata.Client.ServiceName.ResourceClient
subscriptionId := metadata.Client.Account.SubscriptionId

// Use pointer package for pointer operations
// Common scenarios: optional Azure API parameters, nullable fields
enabled := pointer.To(true)
name := pointer.To("example-resource")

// Use structured logging with metadata.Logger
metadata.Logger.Infof("Creating %s", id)

// Use proper error context with typed resource
if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
    return fmt.Errorf("creating %s: %+v", id, err)
}

// Use metadata for resource ID management
metadata.SetID(id)

// Use metadata for state encoding/decoding
var model ServiceNameModel
if err := metadata.Decode(&model); err != nil {
    return fmt.Errorf("decoding: %+v", err)
}
return metadata.Encode(&model)
```

#### untyped Client Usage
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
)

// Standard client initialization
client := meta.(*clients.Client).ServiceName.ResourceClient

// Use pointer package for pointer operations
// Common scenarios: optional Azure API parameters, nullable fields
enabled := pointer.To(d.Get("enabled").(bool))
timeout := pointer.To(d.Get("timeout_seconds").(int))

// Use resource ID parsing for type safety
id := parse.NewResourceID(subscriptionId, resourceGroupName, resourceName)

// Long-running operations
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating Resource %q: %+v", id.ResourceName, err)
}
```

#### Resource ID Management
```go
package servicename

import (
    "fmt"

    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Parse resource IDs consistently
id, err := parse.ResourceID(d.Id())
if err != nil {
    return fmt.Errorf("parsing Resource ID %q: %+v", d.Id(), err)
}

// Set resource ID after creation
d.SetId(id.ID())
```

### Schema Design Standards

#### Common Schema Patterns
```go
package servicename

import (
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/commonschema"
)

// Use common Azure schema helpers
"location": commonschema.Location(),
"resource_group_name": commonschema.ResourceGroupName(),
"tags": commonschema.Tags(),

// Consistent validation
"name": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ForceNew:     true,
    ValidateFunc: validation.StringIsNotEmpty,
},

// Proper ForceNew usage
"size": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    // ForceNew: false allows in-place updates
},
```

#### The "None" Value Pattern

Many Azure APIs accept values like None, Off, or Default as default values and expose them as constants in the API specification. The provider is moving away from exposing these values directly to users, instead leveraging Terraform's native null handling by allowing fields to be omitted. While it is not uncommon to find older resources in the provider that expose and accept these as valid values, the provider is moving away from this pattern, since Terraform has its own null type (i.e., by omitting the field). This ultimately means that the end user does not need to bloat their configuration with superfluous information that is implied through the omission of information. The resulting schema requires a conversion between the Terraform null value and "None" within the Create and Read functions.

**Modern Approach (Preferred):**
```go
package servicename

import (
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

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

// Example validation that excludes "None" - users cannot explicitly set it
"performance_level": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    ValidateFunc: validation.StringInSlice([]string{
        "Low",
        "Medium",
        "High",
        // Note: "None" is NOT included here - handled automatically
    }, false),
},
```

**Typed resource Implementation:**
```go
package servicename

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// Create/Update function - Convert Terraform null to Azure "None"
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Default to "None" if user did not specify a value
            // The resource property shutdown_on_idle maps to the attribute ShutdownOnIdle in the model
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

// Read function - Convert Azure "None" back to Terraform null
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

**Untyped Plugin SDK Implementation:**
```go
package servicename

import (
    "context"

    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Create function - Convert Terraform null to Azure "None"
func resourceServiceCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // Default to "None" if user did not specify a value
    shutdownOnIdle := string(azureapi.ShutdownOnIdleModeNone)
    if v := d.Get("shutdown_on_idle").(string); v != "" {
        shutdownOnIdle = v
    }

    properties := azureapi.ServiceProperties{
        ShutdownOnIdle: &shutdownOnIdle,
    }

    // ...continue with resource creation
    return nil
}

// Read function - Convert Azure "None" back to Terraform null
func resourceServiceRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    // ...retrieve resource from Azure

    // Only set value in state if it is not "None"
    shutdownOnIdle := ""
    if props.ShutdownOnIdle != nil && *props.ShutdownOnIdle != string(azureapi.ShutdownOnIdleModeNone) {
        shutdownOnIdle = *props.ShutdownOnIdle
    }
    d.Set("shutdown_on_idle", shutdownOnIdle)

    return nil
}
```

**Error Handling:**
```go
package servicename

import (
    "fmt"
)

// If user somehow attempts to set "None" explicitly
if model.ShutdownOnIdle == string(azureapi.ShutdownOnIdleModeNone) {
    return fmt.Errorf("property `shutdown_on_idle` cannot be set to `None` - omit the field to use default behavior")
}
```

**Key Principles:**
- **User Experience**: Users omit optional fields instead of explicitly setting "None" values
- **Validation**: Exclude "None", "Off", "Default" values from schema validation
- **Create/Update**: Convert empty/null Terraform values to appropriate Azure default constants
- **Read**: Convert Azure default constants back to empty values in Terraform state
- **Legacy Support**: Existing resources with exposed "None" values are planned for removal in v4.0

**Benefits:**
- Cleaner user configurations without superfluous default values
- Leverages Terraform's native null handling
- Consistent with Terraform best practices
- Reduces configuration bloat

### State Management Standards

For detailed state management patterns including when to use `d.GetRawConfig()` vs `d.Get()` in untyped Plugin SDK resources, see the State Management section in [`coding-patterns.instructions.md`](./coding-patterns.instructions.md).

### Console Output Interpretation

**CONSOLE LINE WRAPPING WARNING**: When reviewing git diff output in terminal/console, be aware that long lines may wrap and appear malformed. Always verify actual file content for syntax validation, especially for JSON, YAML, or structured data files. Console wrapping can make valid syntax appear broken.

**VERIFICATION PROTOCOL FOR SUSPECTED ISSUES**:
- **Before flagging malformed content**: Use `Get-Content filename` (PowerShell) or `cat filename` (bash) to verify file contents
- **JSON Validation**: For JSON files specifically, consider using `Get-Content file.json | ConvertFrom-Json` (PowerShell) or `jq "." file.json` (bash) to validate syntax
- **Console Wrapping Indicators**: 
  - Text breaks mid-sentence or mid-word
  - Missing closing quotes/brackets that don't make logical sense
  - Fragmented lines that appear to continue elsewhere
  - Content looks syntactically invalid but conceptually correct
- **Verification Rule**: If actual file content is valid, acknowledge console wrapping and do not flag as an issue

### Security Considerations

#### Credential and Secret Management

**Never Log Sensitive Information:**
```go
package servicename

import (
    "context"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

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
package servicename

import (
    "fmt"
    "os"
)

// GOOD - Proper environment variable validation
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

// BAD - Hardcoded credentials
const (
    subscriptionID = "12345678-1234-1234-1234-123456789012" // Never hardcode
    clientSecret   = "super-secret-value"                   // Never hardcode
)
```

**Azure Key Vault Integration:**
```go
package servicename

import (
    "context"
    "fmt"
    "strings"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Key Vault reference pattern for typed resource
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Validate Key Vault reference format
            if isKeyVaultReference(model.ConnectionString) {
                if err := validateKeyVaultReference(model.ConnectionString); err != nil {
                    return fmt.Errorf("invalid Key Vault reference format: %+v", err)
                }
            }

            // Use the reference without logging actual value
            properties := azureapi.ServiceProperties{
                ConnectionString: model.ConnectionString, // Azure will resolve the Key Vault reference
            }

            return nil
        },
    }
}

func isKeyVaultReference(value string) bool {
    return strings.HasPrefix(value, "@Microsoft.KeyVault(SecretUri=")
}

func validateKeyVaultReference(reference string) error {
    // Validate Key Vault reference format without logging the actual secret
    if !strings.Contains(reference, "vault.azure.net") {
        return fmt.Errorf("Key Vault reference must use vault.azure.net domain")
    }
    return nil
}
```

#### Input Validation and Sanitization

**Prevent Injection Attacks:**
```go
package servicename

import (
    "fmt"
    "regexp"
    "strings"
)

// GOOD - Proper input validation
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
            // Ensure all error messages follow the consistent format
            errors = append(errors, fmt.Errorf("property `%s` cannot use reserved name `%s`", k, reserved))
            return warnings, errors
        }
    }
    
    return warnings, errors
}

// BAD - No input validation
func ValidateNameUnsafe(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)
    // No validation - vulnerable to injection attacks
    return warnings, errors
}
```

**SQL Injection Prevention in Resource Names:**
```go
package servicename

import (
    "fmt"
    "strings"
)

// GOOD - Sanitize resource names that might be used in SQL contexts
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
            errors = append(errors, fmt.Errorf("property %s cannot contain potentially unsafe characters or SQL keywords", k))
            return warnings, errors
        }
    }
    
    return warnings, errors
}
```

#### Network Security

**TLS and Certificate Validation:**
```go
package servicename

import (
    "context"
    "crypto/tls"
    "fmt"
    "net/http"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Proper TLS configuration
func configureSecureHTTPClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12, // Enforce minimum TLS version
                CipherSuites: []uint16{
                    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
                    tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                },
            },
        },
    }
}

// BAD - Insecure TLS configuration
func configureInsecureHTTPClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true, // Never skip certificate verification
                MinVersion:         tls.VersionTLS10, // Too old
            },
        },
    }
}
```

#### Authentication and Authorization

**Service Principal Best Practices:**
```go
package servicename

import (
    "context"
    "fmt"
    "time"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Proper authentication handling
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            
            // Verify client authentication before proceeding
            if err := verifyClientAuthentication(ctx, client); err != nil {
                return fmt.Errorf("authentication verification failed: %+v", err)
            }
            
            // Use context with timeout for all operations
            ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
            defer cancel()
            
            // Proceed with authenticated operations
            return nil
        },
    }
}

func verifyClientAuthentication(ctx context.Context, client interface{}) error {
    // Implement authentication verification logic
    // This should be done by the Azure SDK, but can be verified if needed
    return nil
}
```

**Token Refresh and Lifecycle Management:**
```go
package servicename

import (
    "context"
    "fmt"
    "time"
)

// GOOD - Proper token lifecycle management
func handleTokenRefresh(ctx context.Context, operation func() error) error {
    const maxRetries = 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        // Check if error is due to token expiration
        if isTokenExpiredError(err) && attempt < maxRetries {
            // Wait before retry (exponential backoff)
            backoffDuration := time.Duration(attempt*attempt) * time.Second
            time.Sleep(backoffDuration)
            continue
        }
        
        return fmt.Errorf("operation failed after %d attempts: %+v", attempt, err)
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

func isTokenExpiredError(err error) bool {
    // Implement logic to detect token expiration errors
    return false // Placeholder
}
```

### Performance Considerations

#### Azure API Rate Limiting and Throttling

**Exponential Backoff Implementation:**
```go
package servicename

import (
    "context"
    "fmt"
    "math"
    "time"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

// GOOD - Proper rate limiting with exponential backoff
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            
            operation := func() error {
                resp, err := client.CreateOrUpdate(ctx, id, properties)
                if err != nil {
                    if response.WasThrottled(resp.HttpResponse) {
                        return &ThrottledError{Err: err}
                    }
                    return err
                }
                return nil
            }
            
            return retryWithExponentialBackoff(ctx, operation, metadata.Logger)
        },
    }
}

type ThrottledError struct {
    Err error
}

func (e *ThrottledError) Error() string {
    return e.Err.Error()
}

func retryWithExponentialBackoff(ctx context.Context, operation func() error, logger interface{}) error {
    const maxRetries = 5
    const baseDelay = 1 * time.Second
    const maxDelay = 32 * time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        // Check if it's a throttling error
        if throttleErr, ok := err.(*ThrottledError); ok {
            if attempt == maxRetries-1 {
                return fmt.Errorf("request throttled after %d attempts: %+v", maxRetries, throttleErr.Err)
            }
            
            // Calculate exponential backoff delay
            delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
            if delay > maxDelay {
                delay = maxDelay
            }
            
            // Log throttling and wait
            if logger != nil {
                // Note: Use appropriate logger based on implementation type
                // metadata.Logger.Infof() for typed resource
                // log.Printf() for untyped
            }
            
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(delay):
                continue
            }
        }
        
        // For non-throttling errors, return immediately
        return err
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
```

#### Efficient Resource Queries

**Batch Operations:**
```go
package servicename

import (
    "context"
    "fmt"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Use batch operations when available
func (r ServiceResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            
            // If Azure API supports batch operations, use them
            if batchClient, ok := client.(BatchOperationClient); ok {
                return r.readWithBatch(ctx, metadata, batchClient)
            }
            
            // Fall back to individual operations
            return r.readIndividual(ctx, metadata, client)
        },
    }
}

type BatchOperationClient interface {
    BatchGet(ctx context.Context, ids []string) ([]Resource, error)
}

func (r ServiceResource) readWithBatch(ctx context.Context, metadata sdk.ResourceMetaData, client BatchOperationClient) error {
    // Implement batch reading logic
    return nil
}

func (r ServiceResource) readIndividual(ctx context.Context, metadata sdk.ResourceMetaData, client interface{}) error {
    // Implement individual reading logic
    return nil
}
```

**Efficient Pagination:**
```go
package servicename

import (
    "context"
    "fmt"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Proper pagination handling
func listResourcesEfficiently(ctx context.Context, client interface{}) ([]Resource, error) {
    var allResources []Resource
    nextLink := ""
    
    for {
        // Use appropriate page size (usually 100-1000 depending on Azure service)
        pageSize := 100
        
        resp, err := client.List(ctx, ListOptions{
            PageSize: pageSize,
            NextLink: nextLink,
        })
        if err != nil {
            return nil, fmt.Errorf("listing resources: %+v", err)
        }
        
        allResources = append(allResources, resp.Resources...)
        
        // Check if there are more pages
        if resp.NextLink == "" {
            break
        }
        nextLink = resp.NextLink
        
        // Prevent infinite loops
        if len(allResources) > 10000 { // Reasonable upper limit
            return nil, fmt.Errorf("too many resources returned, possible infinite pagination")
        }
    }
    
    return allResources, nil
}
```

#### Caching Strategies

**Safe Resource Caching:**
```go
package servicename

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// GOOD - Implement caching for read-only or slowly changing data
type ResourceCache struct {
    cache map[string]CachedResource
    mutex sync.RWMutex
    ttl   time.Duration
}

type CachedResource struct {
    Resource  interface{}
    Timestamp time.Time
}

func NewResourceCache(ttl time.Duration) *ResourceCache {
    return &ResourceCache{
        cache: make(map[string]CachedResource),
        ttl:   ttl,
    }
}

func (c *ResourceCache) Get(ctx context.Context, key string, fetchFunc func(ctx context.Context) (interface{}, error)) (interface{}, error) {
    c.mutex.RLock()
    cached, exists := c.cache[key]
    c.mutex.RUnlock()
    
    // Check if cache hit and not expired
    if exists && time.Since(cached.Timestamp) < c.ttl {
        return cached.Resource, nil
    }
    
    // Cache miss or expired, fetch new data
    resource, err := fetchFunc(ctx)
    if err != nil {
        return nil, err
    }
    
    // Update cache
    c.mutex.Lock()
    c.cache[key] = CachedResource{
        Resource:  resource,
        Timestamp: time.Now(),
    }
    c.mutex.Unlock()
    
    return resource, nil
}

// CAUTION - Only cache immutable or slowly changing data
// Never cache: user data, secrets, frequently changing resources
// Safe to cache: Azure service endpoints, supported VM sizes, available regions (with appropriate TTL)
```

#### Memory Management

**Efficient Resource Processing:**
```go
package servicename

import (
    "context"
    "fmt"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Process large datasets efficiently
func (r ServiceResource) processLargeDataset(ctx context.Context, metadata sdk.ResourceMetaData) error {
    client := metadata.Client.ServiceName.ResourceClient
    
    // Use streaming/paging instead of loading all data at once
    processFunc := func(resource Resource) error {
        // Process individual resource
        return nil
    }
    
    return client.StreamResources(ctx, processFunc)
}

// GOOD - Proper memory cleanup
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Use defer for cleanup
            defer func() {
                // Clean up any temporary resources
                if tempData != nil {
                    tempData.Cleanup()
                }
            }()
            
            // Limit slice capacity for large datasets
            items := make([]Item, 0, 1000) // Set reasonable initial capacity
            
            // Process data
            return nil
        },
    }
}
```

#### Timeout Management

**Appropriate Timeout Configuration:**
```go
package servicename

import (
    "time"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// GOOD - Service-appropriate timeouts
func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Timeouts: &pluginsdk.ResourceTimeout{
            // Short operations (metadata only)
            Read: pluginsdk.DefaultTimeout(5 * time.Minute),
            
            // Medium operations (simple resources)
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
            
            // Long operations (complex resources like clusters)
            // Create: pluginsdk.DefaultTimeout(60 * time.Minute),
        },
    }
}

// For very long operations (like large VM deployments)
func resourceComplexService() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(120 * time.Minute), // 2 hours for complex deployments
            Delete: pluginsdk.DefaultTimeout(90 * time.Minute),  // Longer delete for cleanup
            Update: pluginsdk.DefaultTimeout(60 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(10 * time.Minute),  // Longer read for complex resources
        },
    }
}
```

#### Context Management

**Proper Context Usage:**
```go
package servicename

import (
    "context"
    "fmt"
    "time"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Proper context handling
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            
            // Create child context with shorter timeout for individual operations
            operationCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
            defer cancel()
            
            // Pass context to all Azure API calls
            if err := client.CreateOrUpdateThenPoll(operationCtx, id, properties); err != nil {
                return fmt.Errorf("creating %s: %+v", id, err)
            }
            
            // Check context cancellation between operations
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
                // Continue processing
            }
            
            return nil
        },
    }
}
```

#### Performance Monitoring

**Operation Timing and Metrics:**
```go
package servicename

import (
    "context"
    "fmt"
    "time"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// GOOD - Add timing for performance monitoring
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Complex Azure resources may take time to provision
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            startTime := time.Now()
            defer func() {
                duration := time.Since(startTime)
                metadata.Logger.Infof("Resource creation completed in %v", duration)

                // Note: id and properties are example placeholders representing actual resource identifiers
                // and Azure API parameters that would be defined in the actual implementation context

                // Log slow operations for investigation
                if duration > 5*time.Minute {
                    metadata.Logger.Infof("Slow operation detected: %s took %v", id, duration)
                }
            }()

            client := metadata.Client.ServiceName.ResourceClient

            // Track individual operation timings
            opStart := time.Now()
            if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
                return fmt.Errorf("creating %s: %+v", id, err)
            }
            metadata.Logger.Debugf("Azure API call completed in %v", time.Since(opStart))

            return nil
        },
    }
}
```

### Testing Standards

For comprehensive testing patterns, implementation details, and Azure-specific testing guidelines, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

#### Test Organization Standards
- Place tests in same package with `_test.go` suffix
- Use table-driven tests for multiple scenarios
- Separate unit tests from acceptance tests
- Use meaningful test names following `TestFunctionName_Scenario_ExpectedOutcome`
- Write comprehensive acceptance tests for all resources
- Include import tests for all resources (`data.ImportStep()`)

### Common Error Standards (Both Approaches)
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words. For example, write 'cannot' instead of 'can't' and 'is not' instead of 'isn't'
- Error messages must use '%+v' for verbose error output formatting
- Error messages must be clear, concise, and provide actionable guidance

#### Error Message Examples

**Field Names and Values with Backticks:**
```go
package servicename

import (
    "fmt"
)

// GOOD - Field names and values properly formatted with backticks
return fmt.Errorf("creating Storage Account %q with SKU `%s` in location `%s`: %+v", name, skuName, location, err)
return fmt.Errorf("property `account_tier` must be `Standard` or `Premium`, got `%s`", accountTier)
return fmt.Errorf("field `zones` cannot be set when `availability_set_id` is specified")

// BAD - Missing backticks around field names and values
return fmt.Errorf("creating Storage Account %q with SKU %s in location %s: %+v", name, skuName, location, err)
return fmt.Errorf("property account_tier must be Standard or Premium, got %s", accountTier)
return fmt.Errorf("field zones cannot be set when availability_set_id is specified")
```

**Lowercase, No Punctuation, Descriptive:**
```go
package servicename

import (
    "fmt"
)

// GOOD - Lowercase, no punctuation, descriptive error messages
return fmt.Errorf("creating resource group %q in location %q: %+v", name, location, err)
return fmt.Errorf("updating virtual network %q: %+v", id, err)

// BAD - Incorrect casing, punctuation, or vague messages
return fmt.Errorf("Creating Resource Group %q in Location %q: %v", name, location, err)
return fmt.Errorf("error updating virtual network: %s", err.Error())
```

**Verbose Error Formatting:**
```go
package servicename

import (
    "fmt"
)

// GOOD - Verbose error formatting provides full context
return fmt.Errorf("creating CDN Front Door Profile %q: %+v", name, err)
return fmt.Errorf("updating Network Security Group rules: %+v", err)
return fmt.Errorf("polling for completion of operation: %+v", err)

// BAD - Simple formatting loses error context
return fmt.Errorf("creating CDN Front Door Profile %q: %v", name, err)
return fmt.Errorf("updating Network Security Group rules: %s", err.Error())
return fmt.Errorf("polling for completion of operation: %w", err)
```

**Clear context and actionable information:**
```go
package servicename

import (
    "fmt"
    "log"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

// GOOD - Clear context and actionable information
return fmt.Errorf("creating Storage Account %q: name must be globally unique, try a different name: %+v", name, err)
return fmt.Errorf("VM size `%s` is not available in location `%s`, choose a different size or location", size, location)
return fmt.Errorf("property `disk_size_gb` must be between 1 and 32767, got %d", diskSize)

// BAD - Vague, unhelpful error messages
return fmt.Errorf("creating Storage Account failed: %+v", err)
return fmt.Errorf("VM size problem: %+v", err)
return fmt.Errorf("invalid disk size: %+v", err)
```

**Clear resource not found messaging**
```go
package servicename

import (
    "fmt"
    "log"

    "github.com/hashicorp/go-azure-helpers/lang/pointer"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils/response"
)

// GOOD - Clear resource not found messaging
if response.WasNotFound(resp.HttpResponse) {
    return fmt.Errorf("CDN Front Door Profile %q was not found in Resource Group %q", profileName, resourceGroupName)
}

// Typed resource approach
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// Untyped resource approach
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Storage Account %q was not found - removing from state", id.StorageAccountName)
    d.SetId("")
    return nil
}
```

**Clear parsing error context**
```go
package servicename

import (
    "fmt"

    "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicename/parse"
)

// GOOD - Clear parsing error context
id, err := parse.VirtualMachineID(d.Id())
if err != nil {
    return fmt.Errorf("parsing Virtual Machine ID %q: %+v", d.Id(), err)
}

// Typed resource approach
id, err := parse.ServiceNameID(metadata.ResourceData.Id())
if err != nil {
    return fmt.Errorf("parsing Service ID %q: %+v", metadata.ResourceData.Id(), err)
}
```
**Specific validation context**
```go
package servicename

import (
    "fmt"
    "regexp"
    "strings"
)

// GOOD - Specific validation context
func ValidateStorageAccountName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)
    
    if len(value) < 3 || len(value) > 24 {
        errors = append(errors, fmt.Errorf("property `%s` must be between 3 and 24 characters, got %d characters", k, len(value)))
    }
    
    if !regexp.MustCompile(`^[a-z0-9]+$`).MatchString(value) {
        errors = append(errors, fmt.Errorf("property `%s` can only contain lowercase letters and numbers, got `%s`", k, value))
    }
    
    // Prevent reserved names that conflict with Azure system accounts
    reservedNames := []string{"admin", "root", "system", "default"}
    for _, reserved := range reservedNames {
        if strings.EqualFold(value, reserved) {
            errors = append(errors, fmt.Errorf("property `%s` cannot use reserved name `%s`", k, reserved))
            return warnings, errors
        }
    }
    
    return warnings, errors
}

// BAD - Generic validation without context
func ValidateName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)
    
    if len(value) < 3 {
        errors = append(errors, fmt.Errorf("property `%s` is too short", k))
    }
    
    return warnings, errors
}
```
