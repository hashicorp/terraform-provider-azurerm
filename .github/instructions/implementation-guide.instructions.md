---
applyTo: "internal/**/*.go"
description: Complete implementation guide for Go files in the Terraform AzureRM provider repository. Includes coding standards, patterns, style guidelines, and Azure SDK integration best practices.
---

# Terraform AzureRM Provider Implementation Guide

This comprehensive guide covers all implementation requirements for the Terraform AzureRM provider.

**Quick navigation:** [🏗️ Implementation Patterns](#🏗️-implementation-patterns) | [📏 Coding Standards](#📏-coding-standards) | [🎨 Coding Style](#🎨-coding-style) | [🔧 Azure SDK Integration](#🔧-azure-sdk-integration) | [💡 AI Coding Guidance](#💡-ai-coding-guidance)

## 🏗️ Implementation Patterns

### Implementation Approach Overview

This provider supports two implementation approaches:

#### **Typed Resource Implementation (Preferred)**
- Uses the `internal/sdk` framework with type-safe models
- Employs receiver methods on resource/data source structs
- Features structured state management with `tfschema` tags
- Provides enhanced error handling and logging through metadata
- **Recommended for all new resources and data sources**

#### **Untyped Resource Implementation (Maintenance)**
- Uses traditional Plugin SDK patterns with function-based CRUD
- Employs direct schema manipulation and `d.Set()`/`d.Get()` patterns
- Features traditional error handling and state management
- **Maintained for existing resources but not recommended for new development**

### Typed Resource Structure Pattern

```go
type ServiceNameResourceModel struct {
    Name              string            `tfschema:"name"`
    ResourceGroup     string            `tfschema:"resource_group_name"`
    Location          string            `tfschema:"location"`
    Sku               string            `tfschema:"sku_name"`
    Enabled           bool              `tfschema:"enabled"`
    TimeoutSeconds    int64             `tfschema:"timeout_seconds"`
    Configuration     []ConfigModel     `tfschema:"configuration"`
    Tags              map[string]string `tfschema:"tags"`

    // Computed attributes
    Id                string            `tfschema:"id"`
    Endpoint          string            `tfschema:"endpoint"`
    Status            string            `tfschema:"status"`
}

type ConfigModel struct {
    Setting1 string `tfschema:"setting1"`
    Setting2 string `tfschema:"setting2"`
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
    return &ServiceNameResourceModel{}
}

func (r ServiceNameResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateServiceNameID
}

func (r ServiceNameResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
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

### Typed CRUD Operations Pattern

```go
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            var model ServiceNameResourceModel
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

            state := ServiceNameResourceModel{
                Name:          id.ServiceName,
                ResourceGroup: id.ResourceGroupName,
                Location:      pointer.From(model.Location),
                Tags:          pointer.From(model.Tags),
            }

            if props := model.Properties; props != nil {
                state.Enabled = pointer.FromBool(props.Enabled, false)
                state.Endpoint = pointer.FromString(props.Endpoint, "")
            }

            return metadata.Encode(&state)
        },
    }
}
```

### Untyped Resource Structure Pattern

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
            "location": commonschema.Location(),
            "resource_group_name": commonschema.ResourceGroupName(),
            "tags": commonschema.Tags(),
        },
    }
}

func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient
    subscriptionId := meta.(*clients.Client).Account.SubscriptionId

    name := d.Get("name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    location := azure.NormalizeLocation(d.Get("location").(string))

    id := parse.NewServiceNameID(subscriptionId, resourceGroupName, name)

    existing, err := client.Get(ctx, id)
    if err != nil && !response.WasNotFound(existing.HttpResponse) {
        return fmt.Errorf("checking for existing %s: %+v", id, err)
    }
    if !response.WasNotFound(existing.HttpResponse) {
        return tf.ImportAsExistsError("azurerm_service_name", id.ID())
    }

    parameters := servicenametype.Resource{
        Location: location,
        Properties: &servicenametype.ResourceProperties{
            // Add properties here
        },
    }

    if tagsRaw := d.Get("tags"); tagsRaw != nil {
        parameters.Tags = tags.Expand(tagsRaw.(map[string]interface{}))
    }

    if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
        return fmt.Errorf("creating %s: %+v", id, err)
    }

    d.SetId(id.ID())
    return resourceServiceNameRead(ctx, d, meta)
}
```

### Import Management Pattern

```go
import (
    // Standard library imports first
    "context"
    "fmt"
    "log"
    "time"

    // External dependencies second
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

    // Internal imports last
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)
```

### CustomizeDiff Import Requirements

**IMPORTANT**: The dual import pattern is **only** required for specific scenarios:

**When DUAL IMPORTS are Required (Typed Resources):**
```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)

// When using *schema.ResourceDiff directly in CustomizeDiff functions
CustomizeDiff: pluginsdk.All(
    pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        // Custom validation using *schema.ResourceDiff
        return nil
    }),
),
```

**When SINGLE IMPORT is Sufficient (Untyped Resources):**
```go
import (
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // Only this import needed
)

// When using *pluginsdk.ResourceDiff in CustomizeDiffShim functions
CustomizeDiff: pluginsdk.CustomDiffWithAll(
    pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
        // Custom validation using *pluginsdk.ResourceDiff (which is an alias for *schema.ResourceDiff)
        return nil
    }),
),
```

**Rule of Thumb:**
- **Typed Resources**: Usually need dual imports `schema` and `pluginsdk` when using `*schema.ResourceDiff` directly
- **Untyped Resources**: Usually only need `pluginsdk` import when using `*pluginsdk.ResourceDiff`
- **Check the function signature**: If you see `*pluginsdk.ResourceDiff` or `pluginsdk.CustomDiffWithAll`, single import is sufficient

---
[⬆️ Back to top](#terraform-azurerm-provider-implementation-guide)

## 📏 Coding Standards

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

### Error Handling Standards

#### Typed Resource Error Patterns
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
```

#### UnTyped Error Patterns
```go
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
```

#### Common Error Standards (Both Approaches)
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words
- Error messages must use '%+v' for verbose error output formatting
- Error messages must be clear, concise, and provide actionable guidance

### File Organization

#### Directory Structure
- **Resource files**: `internal/services/[service]/[resource_type]_resource.go`
- **Resource Test files**: Same directory and name as source with `_test.go` suffix
- **Data source files**: `internal/services/[service]/[resource_type]_data_source.go`
- **Utility files**: Group related functions (e.g., `validate.go`, `parse.go`, `flatten.go`, `expand.go`)
- **Registration**: Each service has a `registration.go` file

#### File Naming
- Use snake_case for file names
- Keep files focused on single responsibility
- Aim for files under 1000 lines when possible
- Separate complex logic into utility functions

---
[⬆️ Back to top](#terraform-azurerm-provider-implementation-guide)

## 🎨 Coding Style

### Copyright Header (Required)
All Go files must include this exact copyright header at the top:
```go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
```

### Code Formatting (gofmt/gofumpt Enforced)
- **gofmt**: All Go code must be formatted with `gofmt` (automatically handled by most editors)
- **gofumpt**: Use `gofumpt` for stricter formatting that enforces additional style rules beyond gofmt
- **goimports**: Use `goimports` to automatically organize import statements
- **Indentation**: Use tabs (Go standard, handled by gofmt/gofumpt)
- **Line Length**: Aim for 120 characters max, break longer lines sensibly

### Basic Go Naming Conventions

#### Basic Rules
- **Exported identifiers**: Use PascalCase (e.g., `CreateResource`, `ValidateInput`)
- **Unexported identifiers**: Use camelCase (e.g., `parseResourceID`, `buildParameters`)
- **Acronyms**: Keep as uppercase (e.g., `resourceGroupID`, `vmSSH`, `apiURL`)
- **Interface names**: Often end with 'er' (e.g., `ResourceProvider`, `Validator`)

### Variable Assignment Standards

#### Simplified Variable Assignment Pattern

**PREFERRED - Direct Assignment:**
```go
// Simple, clear assignment
name := d.Get("name").(string)
enabled := d.Get("enabled").(bool)
```

**FORBIDDEN - Unnecessary Variable Assignment:**
```go
// Don't create intermediate variables for simple operations
nameFromConfig := d.Get("name").(string)
name := nameFromConfig
```

### Comment Guidelines

#### Minimal Comment Standards

**⚠️ CRITICAL: Follow ZERO TOLERANCE FOR UNNECESSARY COMMENTS POLICY**

📋 **For complete policy details, enforcement guidelines, decision trees, and comprehensive examples, see:** [Code Clarity Enforcement Guidelines](./code-clarity-enforcement.instructions.md#🚫-zero-tolerance-for-unnecessary-comments-policy)

**Quick Reference - Comments ONLY for:**
- Azure API-specific quirks not obvious from code
- Complex business logic that cannot be simplified
- Azure SDK workarounds for limitations/bugs
- Non-obvious state patterns (PATCH operations, residual state)

**All other comment scenarios are FORBIDDEN - refactor code instead.**

**🔍 MANDATORY JUSTIFICATION:** Every comment requires explicit justification documented in review response explaining which exception case applies and why code cannot be self-explanatory through refactoring.

---
[⬆️ Back to top](#terraform-azurerm-provider-implementation-guide)

## 🔧 Azure SDK Integration

### Pointer Package Usage

**Use the `pointer` package instead of the `utils` package for pointer operations where applicable:**

```go
import (
    "github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// PREFERRED - Use pointer package for creating pointers
stringPtr := pointer.To("example")
intPtr := pointer.To(int64(42))
boolPtr := pointer.To(true)
slicePtr := pointer.To([]string{"item1", "item2"})

// Convert pointers to values with defaults
stringValue := pointer.From(stringPtr)
stringValueWithDefault := pointer.FromString(stringPtr, "default")
intValue := pointer.FromInt64(intPtr, 0)
boolValue := pointer.FromBool(boolPtr, false)

// Azure API Parameter Patterns
parameters := azuretype.CreateParameters{
    Name:     pointer.To("resource-name"),
    Location: pointer.To("eastus"),
    Enabled:  pointer.To(true),
    Tags:     pointer.To(map[string]string{"env": "prod"}),
}

// Nested Structure Patterns
properties := &azuretype.Properties{
    Config: &azuretype.Config{
        Timeout:  pointer.To(int64(300)),
        Retries:  pointer.To(int32(3)),
        Advanced: pointer.To(false),
    },
}

// FORBIDDEN - Manual pointer creation (inconsistent pattern)
name := "resource-name"
namePtr := &name  // Use pointer.To() instead

// FORBIDDEN - Manual nil checks with dereferencing
if props.Name != nil {
    state.Name = *props.Name  // Use pointer.From() instead
}
```

### Pointer Dereferencing Best Practices

**PREFERRED - Use `pointer.From()` for consistent dereferencing:**
```go
// GOOD - Use pointer.From() for safe dereferencing
state.DisplayName = pointer.From(props.DisplayName)
state.Tags = pointer.From(model.Tags)

if props.Api != nil {
    state.ManagedApiId = pointer.From(props.Api.Id)
}
```

**FORBIDDEN - Manual nil checks with dereferencing:**
```go
// FORBIDDEN - Manual nil checks and dereferencing (inconsistent pattern)
if props.DisplayName != nil {
    state.DisplayName = *props.DisplayName
}
```

### Client Management Pattern

#### Typed Resource Client Usage
```go
// Use metadata.Client for accessing clients
client := metadata.Client.ServiceName.ResourceClient
subscriptionId := metadata.Client.Account.SubscriptionId

// Use pointer package for pointer operations
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

#### Untyped Client Usage
```go
// Standard client initialization
client := meta.(*clients.Client).ServiceName.ResourceClient

// Use pointer package for pointer operations
enabled := pointer.To(d.Get("enabled").(bool))
timeout := pointer.To(int64(d.Get("timeout_seconds").(int)))

// Use resource ID parsing for type safety
id := parse.NewResourceID(subscriptionId, resourceGroupName, resourceName)

// Long-running operations
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating Resource `%s`: %+v", id.ResourceName, err)
}
```

### Schema Design Patterns

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
    // Omitting ForceNew defaults to false, allowing in-place updates
},
```

#### ValidateFunc Patterns

If the Azure SDK package offers a `PossibleValuesForFieldName` function, use that in the `validation.StringInSlice` function instead of hardcoding the possible values manually.

```go
// PREFERRED - Use SDK-provided possible values function
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
        false,
    ),
},

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
```

### Expand/Flatten Function Patterns

#### HashiCorp Standard Expand Function Pattern

```go
func expandServiceConfiguration(input []interface{}) *serviceapi.Configuration {
    if len(input) == 0 || input[0] == nil {
        return nil
    }

    raw := input[0].(map[string]interface{})

    return &serviceapi.Configuration{
        Setting1: pointer.To(raw["setting1"].(string)),
        Setting2: pointer.To(raw["setting2"].(bool)),
        Setting3: pointer.To(raw["setting3"].(int)),
    }
}
```

#### HashiCorp Standard Flatten Function Pattern

```go
func flattenServiceConfiguration(input *serviceapi.Configuration) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    return []interface{}{
        map[string]interface{}{
            "setting1": pointer.From(input.Setting1),
            "setting2": pointer.From(input.Setting2),
            "setting3": pointer.From(input.Setting3),
        },
    }
}
```

### Resource ID Management

```go
// Parse resource IDs consistently
id, err := parse.ResourceID(d.Id())
if err != nil {
    return fmt.Errorf("parsing Resource ID `%s`: %+v", d.Id(), err)
}

// Set resource ID after creation
d.SetId(id.ID())
```

---
[⬆️ Back to top](#terraform-azurerm-provider-implementation-guide)

## 💡 AI Coding Guidance

### Smart Code Generation Patterns

#### Resource Implementation Decision Tree
```text
New Resource Request
├─ Implementation Approach
│  ├─ NEW resource/data source → Use Typed Resource Implementation
│  ├─ EXISTING resource maintenance → Continue Untyped Resource Implementation
│  └─ MAJOR refactor → Consider migration to Typed Resource Implementation
│
├─ File Structure Creation
│  ├─ Resource: internal/services/[service]/[resource]_resource.go
│  ├─ Data Source: internal/services/[service]/[resource]_data_source.go
│  ├─ Tests: same directory with _test.go suffix
│  └─ Utilities: parse.go, validate.go, expand.go, flatten.go
│
└─ Implementation Order
   ├─ 1. Define model structs (Typed) or schema (Untyped)
   ├─ 2. Implement CRUD operations
   ├─ 3. Add validation and error handling
   ├─ 4. Create acceptance tests
   └─ 5. Write documentation
```

#### Cross-Implementation Consistency Validation
When working with related Azure resources (like Linux and Windows variants), always verify:
```text
Consistency Checklist
├─ VALIDATION LOGIC
│  ├─ CustomizeDiff functions must be identical across variants
│  ├─ Field requirements must match (if Windows requires X, Linux must too)
│  ├─ Error messages must use identical patterns
│  └─ Default value handling must be consistent
│
├─ DOCUMENTATION
│  ├─ Field descriptions must be identical for shared fields
│  ├─ Note blocks must apply same conditional logic
│  ├─ Examples must demonstrate equivalent patterns
│  └─ Validation rules must be documented consistently
│
└─ TESTING
   ├─ Test coverage must be equivalent between implementations
   ├─ Test naming must follow parallel patterns
   ├─ Helper function naming must use consistent camelCase
   └─ Configuration templates must demonstrate same behaviors
```

#### Template Selection Guide
```go
// TYPED RESOURCE TEMPLATE - Use for NEW resources
type ServiceNameResource struct{}
var _ sdk.Resource = ServiceNameResource{}

func (r ServiceNameResource) ResourceType() string {
    return "azurerm_service_name"
}

func (r ServiceNameResource) ModelObject() interface{} {
    return &ServiceNameResourceModel{}
}

// UNTYPED RESOURCE TEMPLATE - Use for EXISTING resource maintenance
func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,
        Schema: map[string]*pluginsdk.Schema{
            // Schema definitions
        },
    }
}
```

### Efficient Development Workflow

#### Step-by-Step Implementation Checklist
```text
□ 1. ANALYZE REQUEST
  □ Identify Azure service and resource type
  □ Check if resource already exists (grep search)
  □ Determine implementation approach (Typed vs Untyped)

□ 2. SETUP STRUCTURE
  □ Create/locate service directory: internal/services/[service]/
  □ Identify required files: resource, tests, utilities
  □ Check client registration in internal/clients/

□ 3. IMPLEMENT CORE LOGIC
  □ Define model/schema with ALL required Azure properties
  □ Implement Create() with proper validation and error handling
  □ Implement Read() with nil checks and state management
  □ Implement Update() if supported (check Azure API capabilities)
  □ Implement Delete() with proper cleanup

□ 4. ADD VALIDATION & ERROR HANDLING
  □ Implement IDValidationFunc() for resource ID parsing
  □ Add CustomizeDiff for complex Azure API constraints
  □ Use proper error formatting with field names in backticks
  □ Add timeout configurations appropriate for Azure operations

□ 5. CREATE TESTS
  □ Basic test with minimal configuration
  □ RequiresImport test for import conflict detection
  □ Update test if resource supports updates
  □ CustomizeDiff validation tests if applicable

□ 6. WRITE DOCUMENTATION
  □ Resource documentation with examples
  □ Data source documentation if applicable
  □ Import documentation with example resource ID
```

### Common Implementation Patterns

#### Quick Pattern Reference
```go
// AZURE RESOURCE ID PARSING
id, err := parse.ServiceNameID(metadata.ResourceData.Id())
if err != nil {
    return err
}

// AZURE API CLIENT ACCESS (Typed)
client := metadata.Client.ServiceName.ResourceClient

// AZURE API CLIENT ACCESS (Untyped)
client := meta.(*clients.Client).ServiceName.ResourceClient

// ERROR HANDLING WITH CONTEXT
if err != nil {
    return fmt.Errorf("creating %s: %+v", id, err)
}

// AZURE RESOURCE EXISTENCE CHECK
if !response.WasNotFound(existing.HttpResponse) {
    return metadata.ResourceRequiresImport(r.ResourceType(), id)
}

// POINTER OPERATIONS
enabled := pointer.To(true)
value := pointer.From(response.Enabled)
valueWithDefault := pointer.FromString(response.Name, "default")

// AZURE RESOURCE STATE MANAGEMENT (Typed)
metadata.SetID(id)
return metadata.Encode(&model)

// AZURE RESOURCE CLEANUP (Untyped)
d.SetId("")
return nil
```

### Azure-Specific Coding Patterns

#### PATCH Operations Handling
```go
// Azure PATCH operations preserve existing values when fields are omitted
// Always return complete structure with explicit enabled=false for disabled features
func expandPolicy(input []interface{}) *azuretype.Policy {
    result := &azuretype.Policy{
        Feature1: &azuretype.Feature1{
            Enabled: pointer.To(false), // Explicit disable for PATCH
        },
        Feature2: &azuretype.Feature2{
            Enabled: pointer.To(false), // Explicit disable for PATCH
        },
    }

    if len(input) == 0 || input[0] == nil {
        return result // Returns everything disabled
    }

    // Enable only configured features
    raw := input[0].(map[string]interface{})
    if feature1Raw, exists := raw["feature1"]; exists {
        result.Feature1.Enabled = pointer.To(true)
        // Apply configuration...
    }

    return result
}
```

#### CustomizeDiff Validation Patterns

**Typed Resource CustomizeDiff Pattern:**
```go
// NOTE: Typed resources typically use dual imports when using *schema.ResourceDiff directly
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)

// Typed resource CustomizeDiff implementation
func (r ServiceNameResource) CustomizeDiff() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // Azure SKU validation for typed resources
            if model.SkuName == "Premium" && !model.ZoneRedundant {
                return fmt.Errorf("`zone_redundant` must be true for Premium SKU")
            }

            return nil
        },
    }
}
```

**Untyped Resource CustomizeDiff Pattern:**
```go
// NOTE: Untyped resources often use single import with *pluginsdk.ResourceDiff
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // Only this import needed
)

// Untyped resource CustomizeDiff implementation
func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,

        Schema: map[string]*pluginsdk.Schema{
            // Schema definitions
        },

        CustomizeDiff: pluginsdk.CustomDiffWithAll(
            pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Azure SKU validation
                if diff.Get("sku_name").(string) == "Premium" && !diff.Get("zone_redundant").(bool) {
                    return fmt.Errorf("`zone_redundant` must be true for Premium SKU")
                }

                // Azure region constraints
                location := diff.Get("location").(string)
                if location == "West US" && diff.Get("advanced_features").(bool) {
                    return fmt.Errorf("advanced features not available in West US region")
                }

                return nil
            }),
            // Force recreation for immutable Azure properties
            pluginsdk.ForceNewIfChange("location", func(ctx context.Context, old, new, meta interface{}) bool {
                return old.(string) != new.(string)
            }),

            // Programmatic ForceNew for complex state transitions
            pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
                oldSkuProfile, newSkuProfile := diff.GetChange("sku_profile")
                oldSkuProfileList := oldSkuProfile.([]interface{})
                newSkuProfileList := newSkuProfile.([]interface{})

                // Detect complex state transition requiring recreation
                skuProfileBeingRemoved := len(oldSkuProfileList) > 0 && len(newSkuProfileList) == 0
                if skuProfileBeingRemoved {
                    oldSkuName, newSkuName := diff.GetChange("sku_name")

                    // Force recreation for Azure API constraint
                    if oldSkuName.(string) == "Mix" && newSkuName.(string) != "Mix" {
                        if err := diff.ForceNew("sku_profile"); err != nil {
                            return fmt.Errorf("forcing new resource when removing `sku_profile` with `sku_name` change from `Mix`: %+v", err)
                        }
                    }
                }
                return nil
            }),
        ),
    }
}
```

**Key Differences:**
- **Typed Resources**: Use receiver methods and `sdk.ResourceFunc` patterns, validate against model structs
- **Untyped Resources**: Use function-based patterns and `*schema.ResourceDiff` for field access
- **Import Requirements**: Typed typically need dual imports, untyped often use single import
- **Validation Style**: Typed validate against decoded models, untyped use `diff.Get()` patterns

**For Azure-specific CustomizeDiff validation techniques including zero value handling patterns, see:** [Azure Patterns - Zero Value Validation](./azure-patterns.instructions.md#zero-value-validation-pattern)

**Programmatic ForceNew Pattern Explanation:**
Use `diff.ForceNew()` within CustomizeDiffShim when:
1. Complex conditional logic determines if recreation is needed
2. Multiple field changes combine to require ForceNew
3. Azure API constraints require recreation for specific state transitions
4. Static ForceNew: true or ForceNewIfChange cannot express the logic

#### Schema Design Patterns
```go
// Use Azure SDK constants for validation
"match_variable": {
    Type:     pluginsdk.TypeString,
    Required: true,
    ValidateFunc: validation.StringInSlice(
        profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(), // SDK function
        false,
    ),
},

// Azure resource naming pattern
"name": {
    Type:         pluginsdk.TypeString,
    Required:     true,
    ForceNew:     true,
    ValidateFunc: validate.AzureResourceName, // Use validation utilities
},

// Azure common schema helpers
"location": commonschema.Location(),
"resource_group_name": commonschema.ResourceGroupName(),
"tags": commonschema.Tags(),
```

### Testing Pattern Templates

#### Essential Test Structure
```go
// Basic functionality test
func TestAccServiceName_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

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

// Import conflict test
func TestAccServiceName_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

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

// Update test (only if resource supports updates)
func TestAccServiceName_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

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

// Complete test (comprehensive feature demonstration)
func TestAccServiceName_complete(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.complete(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

// TypeSet deduplication test (valid use of Key() validation)
func TestAccServiceName_typeSetDeduplication(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.withDuplicateSetItems(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
                // VALID: Testing TypeSet deduplication behavior - not simple field validation
                check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
            ),
        },
        data.ImportStep(), // Still validates all other field values
    })
}

// CustomizeDiff validation test (if CustomizeDiff is implemented)
func TestAccServiceName_customizeDiffValidation(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config:      r.invalidConfiguration(data),
            ExpectError: regexp.MustCompile("`configuration` is required when `enabled` is true"),
        },
    })
}
```

### Error Handling Best Practices

#### Structured Error Messages
```go
// Field names and values in backticks
return fmt.Errorf("the `fieldName` field cannot be set when `otherField` is `%s`", value)

// Azure resource context
return fmt.Errorf("creating %s: %+v", id, err)

// Import conflicts (Typed)
return metadata.ResourceRequiresImport(r.ResourceType(), id)

// Import conflicts (Untyped)
return tf.ImportAsExistsError("azurerm_service_name", id.ID())

// Resource not found (Typed)
return metadata.MarkAsGone(id)

// Resource not found (Untyped)
log.Printf("[DEBUG] %s was not found - removing from state", id)
d.SetId("")
return nil
```

### Performance Optimization Patterns

#### Efficient Azure API Usage
```go
// Use appropriate timeouts for Azure operations
Timeouts: &pluginsdk.ResourceTimeout{
    Create: pluginsdk.DefaultTimeout(30 * time.Minute), // Long-running operations
    Read:   pluginsdk.DefaultTimeout(5 * time.Minute),  // Quick reads
    Update: pluginsdk.DefaultTimeout(30 * time.Minute), // Potentially long updates
    Delete: pluginsdk.DefaultTimeout(30 * time.Minute), // Cleanup operations
},

// Batch operations when possible
var operations []azuretype.Operation
for _, item := range items {
    operations = append(operations, expandOperation(item))
}
return client.BatchUpdate(ctx, id, operations)

// Use context with timeouts
ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
defer cancel()
```

### Code Quality Enforcement

#### AI Self-Check Patterns
```text
Before submitting code, verify:
□ All errors use %+v formatting
□ Field names in error messages use backticks
□ No hardcoded values in tests
□ Proper timeout configurations
□ Azure SDK constants used for validation
□ Import functionality tested
□ Documentation follows templates
□ CustomizeDiff tested if used
□ Proper pointer usage with pointer package
□ Resource ID parsing implemented correctly
□ Tests use ONLY ExistsInAzure() check with ImportStep() - NO redundant field validation
```

---
[⬆️ Back to top](#terraform-azurerm-provider-implementation-guide)

---

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- 📋 **Code Clarity Enforcement**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
