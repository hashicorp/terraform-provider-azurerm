---
applyTo: "internal/**/*.go"
description: Complete implementation guide for Go files in the Terraform AzureRM provider repository. Includes coding standards, patterns, style guidelines, and Azure SDK integration best practices.
---

# Terraform AzureRM Provider Implementation Guide

This comprehensive guide covers all implementation requirements for the Terraform AzureRM provider. Quick navigation: [🏗️ Implementation Patterns](#-implementation-patterns) | [📏 Standards](#-coding-standards) | [🎨 Style](#-coding-style) | [🔧 Azure Integration](#-azure-sdk-integration)

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

#### CustomizeDiff Import Requirements

When using CustomizeDiff functions, you must import both packages:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"            // For *schema.ResourceDiff
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk" // For helpers
)
```

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

**AVOID - Unnecessary Variable Assignment:**
```go
// Don't create intermediate variables for simple operations
nameFromConfig := d.Get("name").(string)
name := nameFromConfig
```

### Comment Guidelines

#### Minimal Comment Standards

**⚠️ CRITICAL: Avoid Unnecessary Comments**

The AI should **NEVER** add comments to code unless absolutely necessary. Code should be self-documenting through clear variable names, function names, and structure.

**Comments should ONLY be used for:**
- **Azure API-specific quirks or behaviors** that are not obvious from the code
- **Complex business logic** that cannot be made clear through code structure alone
- **Workarounds for Azure SDK limitations** or API bugs
- **Non-obvious state management patterns** (like PATCH operations or residual state handling)
- **Azure service constraints** that require explanation (timeout ranges, SKU limitations, etc.)

**DO NOT comment:**
- Obvious operations like variable assignments
- Standard expand/flatten operations
- Simple struct initialization
- Basic conditional logic
- Self-explanatory function calls
- Routine Azure API calls

## 🔧 Azure SDK Integration

### Pointer Package Usage

**Use the `pointer` package instead of the `utils` package for pointer operations where applicable:**

```go
import (
    "github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// PREFERRED - Use pointer package
stringPtr := pointer.To("example")
intPtr := pointer.To(int64(42))
boolPtr := pointer.To(true)

// Convert pointers to values with defaults
stringValue := pointer.From(stringPtr)
stringValueWithDefault := pointer.FromString(stringPtr, "default")
intValue := pointer.FromInt64(intPtr, 0)
boolValue := pointer.FromBool(boolPtr, false)
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

**AVOID - Manual nil checks with dereferencing:**
```go
// AVOID - Manual nil checks and dereferencing (inconsistent pattern)
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
    // ForceNew: false allows in-place updates
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

## Quick Reference Links

- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- ⚡ **Quick References**: [reference/](./reference/)
