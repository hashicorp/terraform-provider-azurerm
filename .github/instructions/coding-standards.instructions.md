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
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
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
            // Create logic here

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

            // Map response to model and encode
            return metadata.Encode(&model)
        },
    }
}
```

#### UnTyped Resource Implementation

**Standard CRUD Functions:**
```go
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
func resourceServiceNameRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error  
func resourceServiceNameUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
func resourceServiceNameDelete(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error
```

**Resource Schema Patterns:**
```go
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

### Error Handling

#### typed resource Error Patterns
```go
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
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
```

**Standard CustomizeDiff Resource Pattern:**
```go
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
//  INCORRECT - Will cause compilation error
func(ctx context.Context, diff *pluginsdk.ResourceDiff, meta interface{}) error {
    // pluginsdk.ResourceDiff doesn't exist - this will fail to compile
}

//  CORRECT - Must use external schema package
func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
    // Custom validation logic using the correct type
}
```

**Missing Import Error**
If you see errors like `undefined: schema.ResourceDiff`, ensure you have both imports:
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"    // Required for *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"  // Required for helpers
)
```

**Helper Function Usage**
```go
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

### Migration Guidelines

#### Migrating from Untyped to Typed Implementation

**Migration Decision Criteria:**
- **New Resources**: Always use Typed Resource Implementation
- **Major Feature Additions**: Consider migration opportunity when adding significant functionality
- **Bug Fixes**: Maintain existing untyped implementation for simple fixes
- **End-of-Life Planning**: Plan migration for resources approaching major version changes

**Migration Process:**
1. **Assessment Phase**
   - Analyze existing resource complexity and usage patterns
   - Identify breaking changes that may be required
   - Plan migration timeline and testing strategy

2. **Implementation Phase**
   ```go
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
// Ensure consistent resource ID format
func (r ExistingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateExistingResourceID  // Use same parser as untyped version
}

// Maintain backward compatibility in Read operation
func (r ExistingResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Handle both old and new state formats if necessary
            id, err := parse.ExistingResourceID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }
            // Continue with standard typed implementation
            return nil
        },
    }
}
```

### Azure SDK Integration

#### Pointer Package Usage

**Use the `pointer` package instead of the `utils` package for pointer operations where applicable:**

```go
// PREFERRED - Use pointer package
import "github.com/hashicorp/go-azure-helpers/lang/pointer"

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
// AVOID - Legacy utils package patterns (where pointer package can be used)
import "github.com/hashicorp/terraform-provider-azurerm/utils"

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
- **utils package**: For Azure-specific utilities, complex transformations, and legacy compatibility where pointer package doesn't provide equivalent functionality

**Migration Guidelines:**
- **New Code**: Always use `pointer` package for pointer operations
- **Existing Code**: Gradually migrate to `pointer` package during refactoring
- **Legacy Compatibility**: Maintain `utils` package usage only where `pointer` package doesn't provide equivalent functionality

#### typed resource Client Usage
```go
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

### State Management Standards

For detailed state management patterns including when to use `d.GetRawConfig()` vs `d.Get()` in untyped Plugin SDK resources, see the State Management section in [`coding-patterns.instructions.md`](./coding-patterns.instructions.md).

### Testing Standards

For comprehensive testing patterns, implementation details, and Azure-specific testing guidelines, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

#### Test Organization Standards
- Place tests in same package with `_test.go` suffix
- Use table-driven tests for multiple scenarios
- Separate unit tests from acceptance tests
- Use meaningful test names following `TestFunctionName_Scenario_ExpectedOutcome`
- Write comprehensive acceptance tests for all resources
- Include import tests for all resources (`data.ImportStep()`)

