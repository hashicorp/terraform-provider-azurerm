---
applyTo: "internal/**/*.go"
description: This document outlines the coding standards for Go files in the Terraform AzureRM provider repository. It includes naming conventions, file organization, error handling patterns, resource implementation guidelines, and Azure SDK usage standards.
---

# Console Output Interpretation

**🚨 CRITICAL: Console Line Wrapping Detection Protocol 🚨**

**CONSOLE LINE WRAPPING WARNING**: When reviewing git diff output in terminal/console, be aware that long lines may wrap and appear malformed. Always verify actual file content for syntax validation, especially for JSON, YAML, or structured data files. Console wrapping can make valid syntax appear broken.

**VERIFICATION PROTOCOL FOR SUSPECTED ISSUES**:
## 🔍 **MANDATORY VERIFICATION STEPS:**
1. **STOP** - If text appears broken/fragmented, this is likely console wrapping
2. **VERIFY** - Use `Get-Content filename` (PowerShell) or `cat filename` (bash) to check actual file content
3. **VALIDATE** - For JSON/structured files: `Get-Content file.json | ConvertFrom-Json` or `jq "." file.json`

### 🚨 **Console Wrapping Red Flags:**
- ❌ Text breaks mid-sentence or mid-word without logical reason
- ❌ Missing closing quotes/brackets that don't make sense contextually
- ❌ Fragmented lines that appear to continue elsewhere in the diff
- ❌ Content looks syntactically invalid but conceptually correct
- ❌ Long lines in git diff output that suddenly break

#### ✅ **GOLDEN RULE**: If actual file content is valid → acknowledge console wrapping → do NOT flag as corruption

**Verification Rule**: If actual file content is valid, acknowledge console wrapping and do not flag as an issue

## Coding Standards
Given below are the coding standards for the Terraform AzureRM provider which **MUST** be followed.

### Implementation Approach Overview

This provider supports two implementation approaches for new and existing resources:

**For detailed implementation approach documentation, see the main copilot instructions file.**

#### Key Standards for Both Approaches
- Follow consistent naming conventions regardless of implementation pattern
- Maintain identical user-facing behavior and documentation
- Use appropriate error handling patterns for each approach
- Ensure comprehensive testing coverage

- Maintain identical user-facing behavior and documentation
- Use appropriate error handling patterns for each approach
- Ensure comprehensive testing coverage

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

#### Schema Field Naming Conventions

**Descriptive Field Names:**
- Use clear, descriptive names that indicate the field's purpose
- Avoid abbreviations unless they are well-established industry terms
- Use the full Azure service terminology when possible

**Field Naming Pattern Examples:**
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

**Consistency Across Resources:**
- When similar functionality exists across multiple resources, use consistent naming
- Match Azure API documentation terminology where possible
- Maintain consistency between resource and data source field names

#### Data Types for Typed Resource Models
- **Integer fields**: Use `int64` for numeric fields in typed resource models (e.g., `TimeoutSeconds int64`)
- **Pointer operations**: Use `pointer.FromInt64()` and `pointer.To(int64(...))` for int64 fields
- **Schema mapping**: Integer fields map to `pluginsdk.TypeInt` in schema definitions
- **Type conversion**: Convert from `d.Get("field").(int)` to `int64` when using with pointer operations

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
    TimeoutSeconds    int64             `tfschema:"timeout_seconds"`
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

            "location": commonschema.Location(),

            "resource_group_name": commonschema.ResourceGroupName(),
            // ... more schema definitions
        },
    }
}
```

#### ValidateFunc Patterns

If the Azure SDK package offers a `PossibleValuesForFieldName` function, use that in the `validation.StringInSlice` function instead of hardcoding the possible values manually. However, if the Azure SDK package does not offer a `PossibleValuesForFieldName` hardcoding the possible values is acceptable.

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

#### ValidateFunc Patterns

If the Azure SDK package offers a `PossibleValuesForFieldName` function, use that in the `validation.StringInSlice` function instead of hardcoding the possible values manually. However, if the Azure SDK package does not offer a `PossibleValuesForFieldName` hardcoding the possible values is acceptable.

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
    return fmt.Errorf("creating Resource `%s`: %+v", name, err)
}

// Include resource information in error messages
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Resource `%s` was not found - removing from state", id.ResourceName)
    d.SetId("")
    return nil
}

// Handle Azure-specific errors
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled"))
}
```


#### CustomizeDiff Implementation

For detailed CustomizeDiff implementation patterns, import requirements, and examples, see:
- **Implementation Details**: [`coding-patterns.instructions.md`](./coding-patterns.instructions.md) - Dual import requirements and standard patterns
- **Azure-Specific Examples**: [`provider-guidelines.instructions.md`](./provider-guidelines.instructions.md) - Azure resource validation examples

**Quick Reference**: CustomizeDiff functions require dual imports:
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)
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

### Implementation Migration Guide

#### Comprehensive Step-by-Step Migration from Untyped to Typed Resources

This section provides detailed guidance for migrating existing untyped Plugin SDK resources to the modern typed resource implementation approach.

#### Pre-Migration Assessment

**1. Resource Complexity Analysis**
```go
// Assess current resource complexity
// Simple resource (recommended for first migration)
- Single-level schema (no complex nested blocks)
- Basic CRUD operations
- Minimal CustomizeDiff requirements
- Standard Azure resource patterns

// Complex resource (migrate after gaining experience)
- Multi-level nested schemas
- Complex business logic in CRUD operations
- Extensive CustomizeDiff validation
- Non-standard Azure API patterns
```

**2. Dependencies and Impact Assessment**
- **Breaking Changes**: Identify any potential breaking changes to user configurations
- **State Compatibility**: Ensure Terraform state remains compatible
- **Test Coverage**: Verify comprehensive test coverage exists
- **Documentation**: Check if documentation needs updates

#### Migration Process Phases

**Phase 1: Preparation and Planning**

```go
// 1. Create backup branch
git checkout -b migration/resource-name-to-typed

// 2. Analyze existing untyped implementation
// Study these files:
// - internal/services/servicename/resource_name_resource.go
// - internal/services/servicename/resource_name_resource_test.go
// - website/docs/r/service_name_resource.html.markdown

// 3. Document current behavior
// - All schema fields and their types
// - CRUD operation behaviors
// - Error handling patterns
// - CustomizeDiff validations
// - Import functionality
```

**Phase 2: Model Structure Creation**

```go
// Create the typed model structure
type ServiceNameResourceModel struct {
    // Required fields first (alphabetical within category)
    Name              string            `tfschema:"name"`
    ResourceGroup     string            `tfschema:"resource_group_name"`
    Location          string            `tfschema:"location"`

    // Optional fields (alphabetical)
    Enabled           bool              `tfschema:"enabled"`
    Settings          map[string]string `tfschema:"settings"`
    Tags              map[string]string `tfschema:"tags"`

    // Complex nested structures
    Configuration     []ConfigModel     `tfschema:"configuration"`

    // Computed attributes last
    Id                string            `tfschema:"id"`
    Endpoint          string            `tfschema:"endpoint"`
}

// Nested model structures
type ConfigModel struct {
    Setting1 string `tfschema:"setting1"`
    Setting2 string `tfschema:"setting2"`
}
```

**Phase 3: Resource Structure Implementation**

```go
// Implement the resource structure
type ServiceNameResource struct{}

// Required interface implementations
var (
    _ sdk.Resource           = ServiceNameResource{}
    _ sdk.ResourceWithUpdate = ServiceNameResource{} // Only if resource supports updates
)

func (r ServiceNameResource) ResourceType() string {
    return "azurerm_service_name_resource"
}

func (r ServiceNameResource) ModelObject() interface{} {
    return &ServiceNameResourceModel{}
}

func (r ServiceNameResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    // Reuse existing ID validation function
    return parse.ValidateServiceNameResourceID
}
```

**Phase 4: Schema Migration**

```go
// Migrate Arguments (user-configurable fields)
func (r ServiceNameResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        // Copy existing schema definitions from untyped resource
        // Ensure validation functions and defaults are preserved
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },

        "resource_group_name": commonschema.ResourceGroupName(),
        "location": commonschema.Location(),
        "tags": tags.Schema(),

        // Migrate complex nested schemas
        "configuration": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    // Copy nested schema from untyped implementation
                },
            },
        },
    }
}

// Migrate Attributes (computed fields)
func (r ServiceNameResource) Attributes() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "id": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "endpoint": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },
    }
}
```

**Phase 5: CRUD Operation Migration**

```go
// Migrate Create operation
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute, // Copy timeout from untyped resource
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // 1. Client access (new pattern)
            client := metadata.Client.ServiceName.ResourceClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            // 2. Model decoding (new pattern)
            var model ServiceNameResourceModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // 3. Resource ID creation (same as untyped)
            id := parse.NewServiceNameResourceID(subscriptionId, model.ResourceGroup, model.Name)

            // 4. Import check (new pattern)
            metadata.Logger.Infof("Import check for %s", id)
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
            }

            if !response.WasNotFound(existing.HttpResponse) {
                return metadata.ResourceRequiresImport(r.ResourceType(), id)
            }

            // 5. API call (similar to untyped but with model data)
            metadata.Logger.Infof("Creating %s", id)

            properties := servicenametype.Resource{
                Location: model.Location,
                Properties: &servicenametype.ResourceProperties{
                    Enabled: pointer.To(model.Enabled),
                    // Map other model fields to API structure
                },
                Tags: &model.Tags,
            }

            if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
                return fmt.Errorf("creating %s: %+v", id, err)
            }

            // 6. ID setting (new pattern)
            metadata.SetID(id)
            return nil
        },
    }
}

// Migrate Read operation
func (r ServiceNameResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient

            id, err := parse.ServiceNameResourceID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }

            metadata.Logger.Infof("Reading %s", id)
            resp, err := client.Get(ctx, *id)
            if err != nil {
                if response.WasNotFound(resp.HttpResponse) {
                    // New pattern for marking resource as gone
                    return metadata.MarkAsGone(id)
                }
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }

            model := resp.Model
            if model == nil {
                return fmt.Errorf("retrieving %s: model was nil", id)
            }

            // Map API response to typed model
            state := ServiceNameResourceModel{
                Name:          id.ResourceName,
                ResourceGroup: id.ResourceGroupName,
                Location:      model.Location,
                Tags:          model.Tags,
            }

            if props := model.Properties; props != nil {
                state.Enabled = pointer.FromBool(props.Enabled, false)
                state.Endpoint = pointer.FromString(props.Endpoint, "")
                // Map other properties from API response
            }

            // New pattern for encoding state
            return metadata.Encode(&state)
        },
    }
}
```

**Phase 6: Testing Migration**

```go
// Update test patterns
func TestAccServiceNameResource_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name_resource", "test")
    r := ServiceNameResource{} // New pattern: use struct instead of function

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-%d", data.RandomInteger)),
                check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-%d", data.RandomInteger)),
            ),
        },
        data.ImportStep(), // Keep existing import tests
    })
}

// Update Exists function for typed resource
func (r ServiceNameResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := parse.ServiceNameResourceID(state.ID)
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

**Phase 7: Service Registration Update**

```go
// Update service registration to include typed resource
func (r Registration) Resources() []sdk.Resource {
    return []sdk.Resource{
        ServiceNameResource{}, // Add migrated typed resource
        // Keep other typed resources
    }
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
    return map[string]*pluginsdk.Resource{
        // Remove migrated resource from here
        // "azurerm_service_name_resource": resourceServiceNameResource(), // REMOVE

        // Keep other untyped resources
        "azurerm_other_resource": resourceOtherResource(),
    }
}
```

#### Migration Validation Checklist

**Functionality Verification:**
- [ ] All CRUD operations work correctly
- [ ] Resource import functionality preserved
- [ ] CustomizeDiff validations migrated correctly
- [ ] Error handling maintains same user experience
- [ ] Timeout configurations preserved
- [ ] All schema fields and validation rules maintained

**Testing Verification:**
- [ ] All existing acceptance tests pass
- [ ] Import tests continue to work
- [ ] Error scenario tests function correctly
- [ ] Test configurations require no changes
- [ ] Performance characteristics remain similar

**State Compatibility Verification:**
- [ ] Existing Terraform state works without manual intervention
- [ ] Resource attributes remain accessible
- [ ] Computed values continue to be populated correctly
- [ ] No unexpected ForceNew behaviors introduced

**Documentation Verification:**
- [ ] Resource documentation requires no changes
- [ ] Import syntax remains the same
- [ ] Examples continue to work
- [ ] Attribute descriptions remain accurate

#### Common Migration Pitfalls and Solutions

**1. Schema Mismatch Issues**
```go
// Problem: tfschema tag doesn't match schema key
type BadModel struct {
    Name string `tfschema:"resource_name"` // Wrong: should be "name"
}

// Solution: Ensure tfschema tags exactly match schema keys
type GoodModel struct {
    Name string `tfschema:"name"` // Correct: matches schema key
}
```

**2. Nested Structure Complexity**
```go
// Problem: Complex nested structures not properly modeled
// Solution: Break down into separate model structs
type ConfigurationModel struct {
    Setting1 string            `tfschema:"setting1"`
    Setting2 string            `tfschema:"setting2"`
    Nested   []NestedModel     `tfschema:"nested"`
}

type NestedModel struct {
    Value string `tfschema:"value"`
}
```

**3. Import Conflict Detection**
```go
// Problem: Using old import conflict pattern
if existing.StatusCode != http.StatusNotFound {
    return tf.ImportAsExistsError("azurerm_resource", id.ID())
}

// Solution: Use metadata pattern
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}
```

**4. State Management**
```go
// Problem: Trying to use d.Set() in typed resource
d.Set("name", model.Name) // Wrong pattern

// Solution: Use metadata.Encode()
return metadata.Encode(&state) // Correct pattern
```

#### Post-Migration Best Practices

**1. Gradual Migration Strategy**
- Start with simpler resources to gain experience
- Migrate related resources together when possible
- Monitor for any user-reported issues
- Document lessons learned for future migrations

**2. Backward Compatibility Monitoring**
- Test with existing Terraform configurations
- Verify state file compatibility
- Check for any performance regressions
- Monitor for unexpected breaking changes

**3. Code Quality Maintenance**
- Follow typed resource best practices
- Maintain consistent error message formatting
- Use structured logging patterns
- Keep comprehensive test coverage

This comprehensive migration guide provides the foundation for systematically converting untyped resources to the modern typed implementation approach while maintaining backward compatibility and code quality standards.

#### Version Compatibility Matrix

Understanding version compatibility is crucial for planning migrations and ensuring the migration techniques are available in your target provider version.

**Terraform Plugin SDK Compatibility:**

| Feature | Plugin SDK v2.0+ | Plugin SDK v2.10+ | Plugin SDK v2.20+ | Notes |
|---------|------------------|-------------------|-------------------|--------|
| Basic Typed Resources | ❌ | ✅ | ✅ | Minimum version for typed resource framework |
| `metadata.Decode()` | ❌ | ✅ | ✅ | State decoding for typed resources |
| `metadata.Encode()` | ❌ | ✅ | ✅ | State encoding for typed resources |
| `metadata.ResourceRequiresImport()` | ❌ | ✅ | ✅ | Import conflict detection |
| `metadata.MarkAsGone()` | ❌ | ✅ | ✅ | Resource deletion handling |
| `metadata.Logger` | ❌ | ✅ | ✅ | Structured logging |
| `sdk.ResourceWithUpdate` | ❌ | ✅ | ✅ | Update operation interface |
| Enhanced `CustomizeDiff` | ✅ | ✅ | ✅ | Available in all v2.x versions |

**AzureRM Provider Framework Evolution:**

| AzureRM Version | Typed Resources | Migration Support | Dual Registration | Recommendation |
|-----------------|----------------|-------------------|-------------------|----------------|
| v3.0 - v3.50 | ❌ | ❌ | ❌ | Use untyped resources only |
| v3.51 - v3.80 | ⚠️ | ⚠️ | ❌ | Early typed resource support (experimental) |
| v3.81+ | ✅ | ✅ | ✅ | Full typed resource support with migration capabilities |
| v4.0+ (planned) | ✅ | ✅ | ✅ | Preferred typed resource implementation |

**Azure SDK for Go Compatibility:**

| Azure SDK Version | Pointer Helpers | Response Helpers | Polling Support | Migration Impact |
|-------------------|----------------|------------------|-----------------|------------------|
| HashiCorp Go Azure SDK v0.20+ | ✅ | ✅ | ✅ | Full migration support |
| Azure SDK for Go v68+ | ⚠️ | ✅ | ✅ | Limited pointer helper support |
| Legacy Azure SDK | ❌ | ⚠️ | ⚠️ | Migration not recommended |

**Migration Timeline Recommendations:**

```go
// Version-specific migration approach
func planMigration(providerVersion string) MigrationStrategy {
    switch {
    case providerVersion >= "v3.81":
        return MigrationStrategy{
            TypedResources:    true,
            DualRegistration: true,
            Recommendation:   "Full migration support available",
        }
    case providerVersion >= "v3.51":
        return MigrationStrategy{
            TypedResources:    true,
            DualRegistration: false,
            Recommendation:   "Experimental - wait for v3.81+ for production use",
        }
    default:
        return MigrationStrategy{
            TypedResources:    false,
            DualRegistration: false,
            Recommendation:   "Upgrade provider version before migration",
        }
    }
}
```

**Feature Compatibility Checks:**

Before starting migration, verify your environment supports the required features:

```go
// Version compatibility validation
func validateMigrationCompatibility() error {
    // Check Plugin SDK version
    if !hasPluginSDKVersion("v2.10.0") {
        return fmt.Errorf("migration requires Terraform Plugin SDK v2.10.0 or later")
    }

    // Check AzureRM provider framework
    if !hasTypedResourceSupport() {
        return fmt.Errorf("migration requires AzureRM provider v3.81 or later")
    }

    // Check Azure SDK compatibility
    if !hasPointerHelpers() {
        return fmt.Errorf("migration requires HashiCorp Go Azure SDK v0.20 or later")
    }

    return nil
}
```

**Breaking Changes by Version:**

| Version | Breaking Changes | Migration Impact | Mitigation |
|---------|-----------------|------------------|------------|
| v3.81 | Introduction of typed resources | None (backward compatible) | Gradual migration approach |
| v4.0 (planned) | Deprecation warnings for untyped patterns | Development warnings only | Plan migration timeline |
| v5.0 (planned) | Removal of legacy untyped support | Full migration required | Complete migration before upgrade |

**Best Practices by Version:**

- **v3.81 - v3.99**: Start with simple resources, maintain dual registration
- **v4.0+**: Migrate actively developed resources first, plan legacy resource migration
- **v5.0+**: Complete all migrations before upgrading

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
intPtr := pointer.To(int64(42))
boolPtr := pointer.To(true)

// Convert pointers to values with defaults
stringValue := pointer.From(stringPtr)
stringValueWithDefault := pointer.FromString(stringPtr, "default")
intValue := pointer.FromInt64(intPtr, 0)
boolValue := pointer.FromBool(boolPtr, false)

// Check if pointer is nil or has value
if pointer.IsNil(stringPtr) {
    // Handle nil pointer
}
```

#### Pointer Dereferencing Best Practices

**PREFERRED - Use `pointer.From()` for consistent dereferencing:**
```go
package servicename

import (
    "github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// GOOD - Use pointer.From() for safe dereferencing
state.DisplayName = pointer.From(props.DisplayName)
state.Tags = pointer.From(model.Tags)

if props.Api != nil {
    state.ManagedApiId = pointer.From(props.Api.Id)
}
```

**AVOID - Manual nil checks with dereferencing:**
```go
package servicename

import (
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)

// AVOID - Manual nil checks and dereferencing (inconsistent pattern)
if props.DisplayName != nil {
    state.DisplayName = *props.DisplayName
}

if model.Tags != nil {
    state.Tags = *model.Tags
}

// AVOID - Complex nested nil checks
if props.Api != nil && props.Api.Id != nil {
    state.ManagedApiId = *props.Api.Id
}

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
- **Code Review Focus**: Replace manual nil checks with `pointer.From()` for consistent dereferencing patterns

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
    return fmt.Errorf("creating Resource `%s`: %+v", id.ResourceName, err)
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
    return fmt.Errorf("parsing Resource ID `%s`: %+v", d.Id(), err)
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
return fmt.Errorf("creating Storage Account `%s` with SKU `%s` in location `%s`: %+v", name, skuName, location, err)
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
return fmt.Errorf("creating resource group `%s` in location `%s`: %+v", name, location, err)
return fmt.Errorf("updating virtual network `%s`: %+v", id, err)

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
return fmt.Errorf("creating CDN Front Door Profile `%s`: %+v", name, err)
return fmt.Errorf("updating Network Security Group rules: %+v", err)
return fmt.Errorf("polling for completion of operation: %+v", err)

// BAD - Simple formatting loses error context
return fmt.Errorf("creating CDN Front Door Profile `%s`: %v", name, err)
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
return fmt.Errorf("creating Storage Account `%s`: name must be globally unique, try a different name: %+v", name, err)
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
    return fmt.Errorf("CDN Front Door Profile `%s` was not found in Resource Group `%s`", profileName, resourceGroupName)
}

// Typed resource approach
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// Untyped resource approach
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Storage Account `%s` was not found - removing from state", id.StorageAccountName)
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
    return fmt.Errorf("parsing Virtual Machine ID `%s`: %+v", d.Id(), err)
}

// Typed resource approach
id, err := parse.ServiceNameID(metadata.ResourceData.Id())
if err != nil {
    return fmt.Errorf("parsing Service ID `%s`: %+v", metadata.ResourceData.Id(), err)
}
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
    return fmt.Errorf("CDN Front Door Profile `%s` was not found in Resource Group `%s`", profileName, resourceGroupName)
}

// Typed resource approach
if response.WasNotFound(resp.HttpResponse) {
    return metadata.MarkAsGone(id)
}

// Untyped resource approach
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] Storage Account `%s` was not found - removing from state", id.StorageAccountName)
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
    return fmt.Errorf("parsing Virtual Machine ID `%s`: %+v", d.Id(), err)
}

// Typed resource approach
id, err := parse.ServiceNameID(metadata.ResourceData.Id())
if err != nil {
    return fmt.Errorf("parsing Service ID `%s`: %+v", metadata.ResourceData.Id(), err)
}
```
