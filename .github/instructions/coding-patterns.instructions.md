---
applyTo: "internal/**/*.go"
description: This document outlines the coding patterns for Go files in the Terraform AzureRM provider repository. It includes patterns for resource implementation, client management, schema design, and Azure SDK integration.
---

## Coding Patterns for Terraform AzureRM Provider
Given below are the coding patterns for the Terraform AzureRM provider which **MUST** be followed.

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

Both approaches are covered comprehensively below with complete implementation examples.

### Typed Resource Implementation Pattern

The typed resource pattern uses the `internal/sdk` framework which provides a more structured approach with type-safe models and clear separation between arguments and attributes.

#### Typed Resource Structure
```go
type ServiceNameResourceTypeModel struct {
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

type ServiceNameResourceTypeResource struct{}

var (
    _ sdk.Resource           = ServiceNameResourceTypeResource{}
    _ sdk.ResourceWithUpdate = ServiceNameResourceTypeResource{}
)

func (r ServiceNameResourceTypeResource) ResourceType() string {
    return "azurerm_service_name_resource_type"
}

func (r ServiceNameResourceTypeResource) ModelObject() interface{} {
    return &ServiceNameResourceTypeModel{}
}

func (r ServiceNameResourceTypeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateServiceNameResourceTypeID
}

func (r ServiceNameResourceTypeResource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },

        "resource_group_name": commonschema.ResourceGroupName(),

        "location": commonschema.Location(),

        "sku_name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ValidateFunc: validation.StringInSlice([]string{
                "Standard",
                "Premium",
            }, false),
        },

        "enabled": {
            Type:     pluginsdk.TypeBool,
            Optional: true,
            Default:  true,
        },

        "timeout_seconds": {
            Type:         pluginsdk.TypeInt,
            Optional:     true,
            ValidateFunc: validation.IntBetween(1, 3600),
        },

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

        "tags": tags.Schema(),
    }
}

func (r ServiceNameResourceTypeResource) Attributes() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "id": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "endpoint": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "status": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },
    }
}
```

#### Typed CRUD Operations
```go
func (r ServiceNameResourceTypeResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceTypeClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            var model ServiceNameResourceTypeModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            id := parse.NewServiceNameResourceTypeID(subscriptionId, model.ResourceGroup, model.Name)

            metadata.Logger.Infof("Import check for %s", id)
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
            }

            if !response.WasNotFound(existing.HttpResponse) {
                return metadata.ResourceRequiresImport(r.ResourceType(), id)
            }

            metadata.Logger.Infof("Creating %s", id)

            properties := servicenametype.ResourceType{
                Location: model.Location,
                Properties: &servicenametype.ResourceTypeProperties{
                    Enabled: &model.Enabled,
                    // Map other properties
                },
                Sku: &servicenametype.Sku{
                    Name: servicenametype.SkuName(model.Sku),
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

func (r ServiceNameResourceTypeResource) Update() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceTypeClient

            id, err := parse.ServiceNameResourceTypeID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }

            metadata.Logger.Infof("Decoding state for %s", id)
            var state ServiceNameResourceTypeModel
            if err := metadata.Decode(&state); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            metadata.Logger.Infof("Updating %s", id)

            properties := servicenametype.ResourceTypeUpdate{
                Properties: &servicenametype.ResourceTypeProperties{
                    Enabled: &state.Enabled,
                },
                Sku: &servicenametype.Sku{
                    Name: servicenametype.SkuName(state.Sku),
                },
                Tags: &state.Tags,
            }

            if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
                return fmt.Errorf("updating %s: %+v", id, err)
            }

            return nil
        },
    }
}

func (r ServiceNameResourceTypeResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceTypeClient

            id, err := parse.ServiceNameResourceTypeID(metadata.ResourceData.Id())
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

            state := ServiceNameResourceTypeModel{
                Name:          id.ResourceTypeName,
                ResourceGroup: id.ResourceGroupName,
            }

            // Use pointer.From() for consistent pointer dereferencing
            state.Location = pointer.From(model.Location)
            state.Tags = pointer.From(model.Tags)

            if model.Sku != nil {
                state.Sku = string(model.Sku.Name)
            }

            if props := model.Properties; props != nil {
                state.Enabled = pointer.FromBool(props.Enabled, false)
                state.TimeoutSeconds = pointer.FromInt64(props.TimeoutSeconds, 0)
                // Map other properties
            }

            return metadata.Encode(&state)
        },
    }
}

func (r ServiceNameResourceTypeResource) Delete() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceTypeClient

            id, err := parse.ServiceNameResourceTypeID(metadata.ResourceData.Id())
            if err != nil {
                return err
            }

            metadata.Logger.Infof("Deleting %s", id)

            if err := client.DeleteThenPoll(ctx, *id); err != nil {
                return fmt.Errorf("deleting %s: %+v", id, err)
            }

            return nil
        },
    }
}
```

### Typed Data Source Pattern
```go
type ServiceNameResourceTypeDataSourceModel struct {
    Name              string            `tfschema:"name"`
    ResourceGroup     string            `tfschema:"resource_group_name"`
    Location          string            `tfschema:"location"`
    Sku               string            `tfschema:"sku_name"`
    Enabled           bool              `tfschema:"enabled"`
    TimeoutSeconds    int64             `tfschema:"timeout_seconds"`
    Configuration     []ConfigModel     `tfschema:"configuration"`
    Tags              map[string]string `tfschema:"tags"`
    Endpoint          string            `tfschema:"endpoint"`
    Status            string            `tfschema:"status"`
}

type ServiceNameResourceTypeDataSource struct{}

var _ sdk.DataSource = ServiceNameResourceTypeDataSource{}

func (r ServiceNameResourceTypeDataSource) ResourceType() string {
    return "azurerm_service_name_resource_type"
}

func (r ServiceNameResourceTypeDataSource) ModelObject() interface{} {
    return &ServiceNameResourceTypeDataSourceModel{}
}

func (r ServiceNameResourceTypeDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return parse.ValidateServiceNameResourceTypeID
}

func (r ServiceNameResourceTypeDataSource) Arguments() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },

        "resource_group_name": commonschema.ResourceGroupNameForDataSource(),
    }
}

func (r ServiceNameResourceTypeDataSource) Attributes() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "location": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "sku_name": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "enabled": {
            Type:     pluginsdk.TypeBool,
            Computed: true,
        },

        "timeout_seconds": {
            Type:     pluginsdk.TypeInt,
            Computed: true,
        },

        "configuration": {
            Type:     pluginsdk.TypeList,
            Computed: true,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "setting1": {
                        Type:     pluginsdk.TypeString,
                        Computed: true,
                    },
                    "setting2": {
                        Type:     pluginsdk.TypeString,
                        Computed: true,
                    },
                },
            },
        },

        "endpoint": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "status": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },

        "tags": tags.SchemaDataSource(),
    }
}

func (r ServiceNameResourceTypeDataSource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.ServiceName.ResourceTypeClient
            subscriptionId := metadata.Client.Account.SubscriptionId

            var state ServiceNameResourceTypeDataSourceModel
            if err := metadata.Decode(&state); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            id := parse.NewServiceNameResourceTypeID(subscriptionId, state.ResourceGroup, state.Name)

            resp, err := client.Get(ctx, id)
            if err != nil {
                if response.WasNotFound(resp.HttpResponse) {
                    return fmt.Errorf("%s was not found", id)
                }
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }

            model := resp.Model
            if model == nil {
                return fmt.Errorf("retrieving %s: model was nil", id)
            }

            state.Name = id.ResourceTypeName
            state.ResourceGroup = id.ResourceGroupName

            // Use pointer.From() for consistent pointer dereferencing
            state.Location = pointer.From(model.Location)
            state.Tags = pointer.From(model.Tags)

            if model.Sku != nil {
                state.Sku = string(model.Sku.Name)
            }

            if props := model.Properties; props != nil {
                state.Enabled = pointer.FromBool(props.Enabled, false)
                state.TimeoutSeconds = pointer.FromInt64(props.TimeoutSeconds, 0)
                // Map other properties and computed attributes
            }

            metadata.SetID(id)
            return metadata.Encode(&state)
        },
    }
}
```

### unTyped Resource Implementation Pattern (Plugin SDK)

The untyped resource pattern uses the traditional Plugin SDK approach with direct schema manipulation. This pattern is still used in many existing resources and should be maintained for backward compatibility.

#### untyped Resource Structure
```go
func resourceServiceNameResourceType() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameResourceTypeCreate,
        Read:   resourceServiceNameResourceTypeRead,
        Update: resourceServiceNameResourceTypeUpdate,
        Delete: resourceServiceNameResourceTypeDelete,

        Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
            _, err := parse.ServiceNameResourceTypeID(id)
            return err
        }),

        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },

        Schema: resourceServiceNameResourceTypeSchema(),
    }
}

func resourceServiceNameResourceTypeSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:         pluginsdk.TypeString,
            Required:     true,
            ForceNew:     true,
            ValidateFunc: validation.StringIsNotEmpty,
        },

        "location": commonschema.Location(),

        "resource_group_name": commonschema.ResourceGroupName(),

        "tags": commonschema.Tags(),
    }
}
```

#### untyped CRUD Operation Pattern
```go
func resourceServiceNameResourceTypeCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceTypeClient
    subscriptionId := meta.(*clients.Client).Account.SubscriptionId

    // Parse input parameters
    name := d.Get("name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    location := azure.NormalizeLocation(d.Get("location").(string))

    // Create resource ID
    id := parse.NewServiceNameResourceTypeID(subscriptionId, resourceGroupName, name)

    // Check for existing resource
    existing, err := client.Get(ctx, id)
    if err != nil && !response.WasNotFound(existing.HttpResponse) {
        return fmt.Errorf("checking for existing %s: %w", id, err)
    }
    if !response.WasNotFound(existing.HttpResponse) {
        return tf.ImportAsExistsError("azurerm_service_name_resource_type", id.ID())
    }

    // Build parameters
    parameters := servicenametype.ResourceType{
        Location: location,
        Properties: &servicenametype.ResourceTypeProperties{
            // Add properties here
        },
    }

    // Handle tags
    if tagsRaw := d.Get("tags"); tagsRaw != nil {
        parameters.Tags = tags.Expand(tagsRaw.(map[string]interface{}))
    }

    // Create resource
    if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
        return fmt.Errorf("creating %s: %w", id, err)
    }

    d.SetId(id.ID())
    return resourceServiceNameResourceTypeRead(ctx, d, meta)
}

func resourceServiceNameResourceTypeRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceTypeClient

    id, err := parse.ServiceNameResourceTypeID(d.Id())
    if err != nil {
        return err
    }

    resp, err := client.Get(ctx, *id)
    if err != nil {
        if response.WasNotFound(resp.HttpResponse) {
            log.Printf("[DEBUG] %s was not found - removing from state!", *id)
            d.SetId("")
            return nil
        }
        return fmt.Errorf("retrieving %s: %w", *id, err)
    }

    d.Set("name", id.ResourceTypeName)
    d.Set("resource_group_name", id.ResourceGroupName)

    if model := resp.Model; model != nil {
        d.Set("location", azure.NormalizeLocation(model.Location))

        if props := model.Properties; props != nil {
            // Set properties
        }

        if err := tags.FlattenAndSet(d, model.Tags); err != nil {
            return err
        }
    }

    return nil
}
```

#### State Management with d.GetRawConfig()

**When to Use `d.GetRawConfig()` vs `d.Get()` (untyped Resources Only):**

`d.GetRawConfig()` should be used in specific scenarios where you need to distinguish between user-configured values and computed/default values. This method is only available in untyped Plugin SDK resource implementations.

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

// 1b. Alternative pattern using AsValueMap() for multiple fields
func resourceServiceNameUpdateAlternative(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient

    id, err := parse.ServiceNameID(d.Id())
    if err != nil {
        return err
    }

    parameters := serviceapi.UpdateParameters{
        Name: d.Get("name").(string),
    }

    // Get all raw config values at once for multiple checks
    rawConfig := d.GetRawConfig().AsValueMap()

    // Check if user explicitly configured sampling_percentage
    samplingPercentage := rawConfig["sampling_percentage"]
    if !samplingPercentage.IsNull() {
        parameters.SamplingSettings = &serviceapi.SamplingSettings{
            SamplingType: serviceapi.SamplingTypeFixed,
            Percentage:   d.Get("sampling_percentage").(float64),
        }
    } else {
        parameters.SamplingSettings = nil
    }

    // Check if user explicitly configured timeout_seconds
    timeoutSeconds := rawConfig["timeout_seconds"]
    if !timeoutSeconds.IsNull() {
        timeoutValue := d.Get("timeout_seconds").(int)
        parameters.TimeoutSeconds = &timeoutValue
    }

    if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
        return fmt.Errorf("updating %s: %+v", *id, err)
    }

    return nil
}

// 2. Handling optional complex blocks that should be omitted when not configured
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient
    subscriptionId := meta.(*clients.Client).Account.SubscriptionId

    name := d.Get("name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    location := azure.NormalizeLocation(d.Get("location").(string))

    id := parse.NewServiceNameID(subscriptionId, resourceGroupName, name)

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

    if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
        return fmt.Errorf("creating %s: %+v", id, err)
    }

    d.SetId(id.ID())
    return nil
}

// 3. Preserving Azure service defaults vs Terraform defaults
func resourceServiceNameRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient

    id, err := parse.ServiceNameID(d.Id())
    if err != nil {
        return err
    }

    resp, err := client.Get(ctx, *id)
    if err != nil {
        if response.WasNotFound(resp.HttpResponse) {
            log.Printf("[DEBUG] %s was not found - removing from state", *id)
            d.SetId("")
            return nil
        }
        return fmt.Errorf("retrieving %s: %+v", *id, err)
    }

    d.Set("name", id.ServiceName)
    d.Set("resource_group_name", id.ResourceGroupName)

    if model := resp.Model; model != nil {
        d.Set("location", azure.NormalizeLocation(model.Location))

        if props := model.Properties; props != nil {
            // Only set in state if user originally configured it
            if raw := d.GetRawConfig().GetAttr("timeout_seconds"); !raw.IsNull() {
                d.Set("timeout_seconds", props.TimeoutSeconds)
            }
            // If user never set timeout_seconds, don't store Azure's default in state
        }
    }

    return nil
}
```

**When NOT to Use `d.GetRawConfig()`:**
```go
// AVOID: Using GetRawConfig for required fields
func resourceServiceNameCreate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient

    // INCORRECT - required fields should always use d.Get()
    var name string
    if raw := d.GetRawConfig().GetAttr("name"); !raw.IsNull() {
        name = d.Get("name").(string)
    }

    // CORRECT - required fields always use d.Get()
    name := d.Get("name").(string)

    // Continue with resource creation...
    return nil
}

// AVOID: Using GetRawConfig when you always need the value
func resourceServiceNameUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).ServiceName.ResourceClient

    // INCORRECT - if you always need the value, use d.Get()
    var enabled bool
    if raw := d.GetRawConfig().GetAttr("enabled"); !raw.IsNull() {
        enabled = d.Get("enabled").(bool)
    } else {
        enabled = false // This adds unnecessary complexity
    }

    // CORRECT - use d.Get() with proper default handling
    enabled := d.Get("enabled").(bool) // Schema default will be used if not set

    // Use the enabled value in API call...
    return nil
}

// AVOID: Using GetRawConfig in typed resource implementations
// Typed resources use metadata.Decode() patterns instead
func (r ServiceNameResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // CORRECT for typed resources - use metadata.Decode()
            var model ServiceNameModel
            if err := metadata.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            // DON'T try to use d.GetRawConfig() in typed resources
            // The metadata pattern handles this automatically
            return nil
        },
    }
}
```

**Schema Design for `d.GetRawConfig()` Patterns:**
```go
// When planning to use GetRawConfig, consider schema defaults carefully
func resourceServiceNameSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "timeout_seconds": {
            Type:     pluginsdk.TypeInt,
            Optional: true,
            // NO Default here if you want to distinguish between
            // user-set vs Azure service default
        },

        "advanced_config": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            // GetRawConfig useful here to distinguish between
            // empty block {} vs not configured at all
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "setting": {
                        Type:     pluginsdk.TypeString,
                        Required: true,
                    },
                },
            },
        },

        "enabled": {
            Type:     pluginsdk.TypeBool,
            Optional: true,
            Default:  true,
            // GetRawConfig not needed here - default handling is sufficient
        },
    }
}
```

**Best Practices for `d.GetRawConfig()`:**
- **Use sparingly**: Only when you need to distinguish between explicit configuration and defaults
- **Document intent**: Add comments explaining why GetRawConfig is necessary
- **Consider alternatives**: Often proper schema defaults eliminate the need for GetRawConfig
- **untyped only**: This pattern is only available in untyped Plugin SDK resources
- **Validate necessity**: Ask if the Azure API behavior truly requires distinguishing explicit vs default values
- **Test thoroughly**: Ensure behavior works correctly with both explicit and omitted values

**Anti-Patterns to Avoid:**
- **AVOID**: Overusing GetRawConfig when d.Get() with defaults would work
- **AVOID**: Using GetRawConfig for required fields
- **AVOID**: Using GetRawConfig when the value is always needed
- **AVOID**: Complex conditional logic that could be simplified with proper defaults
- **AVOID**: Using GetRawConfig in typed resource implementations (use metadata patterns instead)

### Import Management Pattern

#### Standard Go Import Organization
```go
package servicename

import (
    // Standard library imports first
    "context"
    "fmt"
    "log"
    "regexp"
    "time"

    // External dependencies second
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

    // Internal imports last
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/common"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)
```

#### CustomizeDiff Import Requirements

**Important**: When using CustomizeDiff functions, you must import both packages:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)
```

**Why Both Imports are Required:**
- pluginsdk provides the main resource types (pluginsdk.Resource, pluginsdk.Schema, etc.)
- schema provides CustomizeDiff function types that are **not** aliased in the internal pluginsdk package
- CustomizeDiff functions must use *schema.ResourceDiff parameter type from the external package

**Correct CustomizeDiff Pattern:**
```go
func resourceServiceName() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceServiceNameCreate,
        Read:   resourceServiceNameRead,
        Update: resourceServiceNameUpdate,
        Delete: resourceServiceNameDelete,

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
            // Schema definition using pluginsdk types
        },

        CustomizeDiff: pluginsdk.CustomDiffWithAll(
            // CustomizeDiff functions must use *schema.ResourceDiff with CustomizeDiffShim wrapper
            pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
                // Validation logic here
                return nil
            }),
        ),
    }
}
```

**What NOT to do:**
```go
//  INCORRECT - Cannot use pluginsdk.ResourceDiff (doesn't exist)
func(ctx context.Context, diff *pluginsdk.ResourceDiff, meta interface{}) error {
    // This will cause compilation errors
}
```

#### Azure-Specific CustomizeDiff Validation Example

**Key Points for Azure CustomizeDiff Validation:**
- **Azure API Constraints**: Validate field combinations that are specific to Azure service requirements
- **Clear Error Messages**: Include field paths and specific constraint violations in error messages
- **Field Name Formatting**: Wrap field names in backticks for clarity (`selector`, `match_variable`)
- **Conditional Logic**: Handle complex Azure resource configuration dependencies

### Azure PATCH Operation Pattern

#### Critical PATCH Behavior Understanding

**Azure Resource Manager PATCH Operations:**
Many Azure services use PATCH operations for resource updates, which have fundamentally different behavior from PUT operations:

- **PATCH preserves existing values** when fields are omitted from the request
- **PUT replaces the entire resource** with the provided configuration
- **Azure SDK nil filtering** removes nil values before sending requests to Azure
- **Residual state persistence** means previously enabled features remain active unless explicitly disabled

**PATCH Operation Challenges:**
```go
// PROBLEM: This approach fails with PATCH operations
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    if len(input) == 0 {
        return nil // SDK filters this out, Azure never gets disable command
    }

    // Only sends enabled policies, disabled policies remain unchanged
    result := &azuretype.Policy{}
    // Configure only enabled features...
    return result
}
```

**SOLUTION: Explicit Disable Pattern for PATCH Operations:**
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

**Key PATCH Operation Principles:**
- **Always return complete structures** - Never return nil for PATCH operations
- **Explicit disable commands** - Use `enabled=false` instead of omitting fields
- **Required field compliance** - Provide all required fields even when features are disabled
- **Residual state management** - Ensure previously enabled features can be properly disabled
- **"None" Pattern Integration** - Combine with flatten functions that exclude disabled features from state

**Documentation Requirements for PATCH Operations:**
```go
// PATCH Behavior Note: Azure VMSS uses PATCH operations which preserve existing values
// when fields are omitted. This means previously enabled policies will remain active
// unless explicitly disabled with enabled=false. Sending nil values gets filtered out
// by the Azure SDK, so Azure never receives disable commands. We must explicitly
// send enabled=false for all policies that should be disabled.
```

#### Real-World PATCH + "None" Pattern Integration

**Resiliency Policy Example - Key Learnings:**
The Azure VMSS resiliency policy implementation demonstrates the critical integration of PATCH operations with the "None" pattern:

```go
func ExpandVirtualMachineScaleSetResiliency(input []interface{}, resilientVMCreationEnabled, resilientVMDeletionEnabled bool) *virtualmachinescalesets.ResiliencyPolicy {
    // PATCH Behavior Note: Azure VMSS uses PATCH operations which preserve existing values
    // when fields are omitted. This means previously enabled policies will remain active
    // unless explicitly disabled with enabled=false. Sending nil values gets filtered out
    // by the Azure SDK, so Azure never receives disable commands. We must explicitly
    // send enabled=false for all policies that should be disabled.

    // Define result with all policies disabled by default
    // Azure requires all fields to be present even when disabled
    result := &virtualmachinescalesets.ResiliencyPolicy{
        AutomaticZoneRebalancingPolicy: &virtualmachinescalesets.AutomaticZoneRebalancingPolicy{
            Enabled:           pointer.To(false),                                                       // Disabled by default
            RebalanceBehavior: pointer.To(virtualmachinescalesets.RebalanceBehaviorCreateBeforeDelete), // Required even when disabled
            RebalanceStrategy: pointer.To(virtualmachinescalesets.RebalanceStrategyRecreate),           // Required even when disabled
        },
        ResilientVMCreationPolicy: &virtualmachinescalesets.ResilientVMCreationPolicy{
            Enabled: pointer.To(resilientVMCreationEnabled),
        },
        ResilientVMDeletionPolicy: &virtualmachinescalesets.ResilientVMDeletionPolicy{
            Enabled: pointer.To(resilientVMDeletionEnabled),
        },
    }

    // Handle automatic zone rebalancing configuration - "None" pattern integration
    if len(input) > 0 && input[0] != nil {
        raw := input[0].(map[string]interface{})

        // Enable automatic zone rebalancing if explicitly configured
        if automaticZoneRebalancingRaw, exists := raw["automatic_zone_rebalancing"]; exists {
            automaticZoneRebalancingList := automaticZoneRebalancingRaw.([]interface{})
            if len(automaticZoneRebalancingList) > 0 && automaticZoneRebalancingList[0] != nil {
                // User explicitly configured this feature - enable it
                result.AutomaticZoneRebalancingPolicy.Enabled = pointer.To(true)
                // Apply user configuration...
            }
            // If exists but empty block, feature remains disabled (following "None" pattern)
        }
        // If not configured at all, feature remains disabled (following "None" pattern)
    }

    return result
}
```

**Critical Debugging Insights:**
- **State Management Issue**: Originally returned `nil` when no configuration provided, causing residual state persistence
- **PATCH Operation Interaction**: Azure API preserved previously enabled policies because disable commands were never sent
- **"None" Pattern Compliance**: Users omit optional fields instead of explicitly setting them to default values
- **Required Field Handling**: Even disabled features need all required fields populated with appropriate defaults

**Testing Implications:**
- **Test Scenarios**: Must test enabling → disabling → re-enabling features to catch residual state bugs
- **State Verification**: Verify that removing configuration truly disables Azure features (not just Terraform state)
- **PATCH Operation Testing**: Specifically test update scenarios where features are toggled on/off

### Client Management Pattern

#### Client Registration
```go
type Client struct {
    ResourceTypeClient *servicenametype.ResourceTypeClient
}

func NewClient(o *common.ClientOptions) *Client {
    resourceTypeClient := servicenametype.NewResourceTypeClientWithBaseURI(o.ResourceManagerEndpoint)
    o.ConfigureClient(&resourceTypeClient.Client, o.ResourceManagerAuthorizer)

    return &Client{
        ResourceTypeClient: &resourceTypeClient,
    }
}
```

#### Azure Client Usage Patterns

**Typed Resource Client Usage:**
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

**untyped Client Usage:**
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
timeout := pointer.To(int64(d.Get("timeout_seconds").(int)))

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

### untyped Data Source Pattern

#### Standard Data Source Structure
```go
func dataSourceServiceNameResourceType() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Read: dataSourceServiceNameResourceTypeRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

            "disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_settings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"disk_encryption_key": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"secret_url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

            "tags": commonschema.TagsDataSource(),
        },
    }
}

func dataSourceServiceNameResourceTypeRead(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Compute.DisksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServiceNameResourceTypeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		storageAccountType := ""
		if sku := model.Sku; sku != nil {
			storageAccountType = string(*sku.Name)
		}
		d.Set("storage_account_type", storageAccountType)

		if props := model.Properties; props != nil {
			creationData := props.CreationData

			diskEncryptionSetId := ""
			if props.Encryption != nil && props.Encryption.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *props.Encryption.DiskEncryptionSetId
			}
			d.Set("disk_encryption_set_id", diskEncryptionSetId)

			if err := d.Set("encryption_settings", flattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
```

### Schema Design Patterns

#### Complex Schema with Validation
```go
func resourceComplexSchema() map[string]*pluginsdk.Schema {
    return map[string]*pluginsdk.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
            ValidateFunc: validation.All(
                validation.StringIsNotEmpty,
                validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "name can only contain alphanumeric characters and hyphens"),
                validation.StringLenBetween(1, 64),
            ),
        },

        "configuration": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "enabled": {
                        Type:     pluginsdk.TypeBool,
                        Optional: true,
                        Default:  true,
                    },
                    "settings": {
                        Type:     pluginsdk.TypeMap,
                        Optional: true,
                        Elem: &pluginsdk.Schema{
                            Type: pluginsdk.TypeString,
                        },
                    },
                },
            },
        },

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
    }
}
```

### Error Handling Patterns

#### Standard Error Handling
```go
// Check for existing resource
existing, err := client.Get(ctx, id)
if err != nil {
    if !response.WasNotFound(existing.HttpResponse) {
        return fmt.Errorf("checking for existing %s: %w", id, err)
    }
}

// Handle resource not found in Read operation
if response.WasNotFound(resp.HttpResponse) {
    log.Printf("[DEBUG] %s was not found - removing from state", id)
    d.SetId("")
    return nil
}

// Handle throttling
if response.WasThrottled(resp.HttpResponse) {
    return resource.RetryableError(fmt.Errorf("request was throttled, retrying"))
}

// Handle long-running operations
if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
    return fmt.Errorf("creating %s: %w", id, err)
}
```

### Resource ID Parsing Pattern

#### ID Parser Implementation
```go
type ServiceNameResourceTypeId struct {
    SubscriptionId      string
    ResourceGroupName   string
    ResourceTypeName    string
}

func NewServiceNameResourceTypeID(subscriptionId, resourceGroupName, resourceTypeName string) ServiceNameResourceTypeId {
    return ServiceNameResourceTypeId{
        SubscriptionId:    subscriptionId,
        ResourceGroupName: resourceGroupName,
        ResourceTypeName:  resourceTypeName,
    }
}

func (id ServiceNameResourceTypeId) ID() string {
    fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceName/resourceTypes/%s"
    return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceTypeName)
}

func ServiceNameResourceTypeID(input string) (*ServiceNameResourceTypeId, error) {
    parser := resourceids.NewParserFromResourceIdType(ServiceNameResourceTypeId{})
    parsed, err := parser.Parse(input, false)
    if err != nil {
        return nil, fmt.Errorf("parsing %q: %w", input, err)
    }

    var ok bool
    id := ServiceNameResourceTypeId{}

    if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
    }

    if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
    }

    if id.ResourceTypeName, ok = parsed.Parsed["resourceTypeName"]; !ok {
        return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceTypeName", *parsed)
    }

    return &id, nil
}
```

### Progressive Code Simplification Patterns

#### When Complex Logic Needs Simplification

When implementing Azure resource expand/flatten functions, especially for PATCH operations, complex conditional logic can often be simplified through strategic refactoring.

**Anti-Pattern: Complex Conditional Logic**
```go
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    result := &azuretype.Policy{}

    if len(input) == 0 || input[0] == nil {
        // Create disabled policies with complex duplication
        result.FeatureA = &azuretype.FeatureA{
            Enabled: pointer.To(false),
            RequiredField: pointer.To(azuretype.DefaultValue),
        }
        result.FeatureB = &azuretype.FeatureB{
            Enabled: pointer.To(false),
        }
        return result
    }

    raw := input[0].(map[string]interface{})

    // More complex conditional logic with duplication...
    if featureARaw, exists := raw["feature_a"]; exists {
        // Duplicate policy creation logic
        result.FeatureA = &azuretype.FeatureA{
            Enabled: pointer.To(true),
            RequiredField: pointer.To(azuretype.DefaultValue),
        }
        // Configure...
    } else {
        // Duplicate disabled policy creation
        result.FeatureA = &azuretype.FeatureA{
            Enabled: pointer.To(false),
            RequiredField: pointer.To(azuretype.DefaultValue),
        }
    }

    return result
}
```

**Progressive Simplification Strategy:**

**Step 1: Extract Common Patterns**
```go
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    // Extract disabled policy as variable to reduce duplication
    disabledFeatureA := &azuretype.FeatureA{
        Enabled: pointer.To(false),
        RequiredField: pointer.To(azuretype.DefaultValue),
    }

    result := &azuretype.Policy{}

    if len(input) == 0 || input[0] == nil {
        result.FeatureA = disabledFeatureA
        result.FeatureB = &azuretype.FeatureB{Enabled: pointer.To(false)}
        return result
    }

    raw := input[0].(map[string]interface{})

    // Start with disabled, enable if configured
    result.FeatureA = disabledFeatureA
    if featureARaw, exists := raw["feature_a"]; exists {
        // Only create enabled version when needed
        result.FeatureA = &azuretype.FeatureA{
            Enabled: pointer.To(true),
            RequiredField: pointer.To(azuretype.DefaultValue),
        }
        // Configure...
    }

    return result
}
```

**Step 2: Define Complete Disabled Structure (Recommended)**
```go
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    // Define complete result with all features disabled by default
    result := &azuretype.Policy{
        FeatureA: &azuretype.FeatureA{
            Enabled: pointer.To(false),                    // Disabled by default
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

**Benefits of Progressive Simplification:**
- **Readability**: Clear intent - start disabled, enable what's configured
- **Maintainability**: Single source of truth for structure definition
- **Testability**: Easier to test individual feature enabling/disabling
- **PATCH Compliance**: Naturally handles Azure PATCH operation requirements
- **Reduced Duplication**: No repeated policy creation logic

**Key Simplification Principles:**
1. **Define end state first** - Create complete structure with desired defaults
2. **Use simple field flipping** - Change only what needs to change based on input
3. **Eliminate conditional returns** - Single return path reduces complexity
4. **Extract common patterns** - Use variables for repeated structures
5. **Start with working code** - Simplify incrementally, don't rewrite from scratch

### Expand/Flatten Function Patterns

#### HashiCorp Standard Expand Function Pattern

**Simple Configuration Blocks:**
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

**Array Configuration Blocks:**
```go
func expandServiceRules(input []interface{}) *[]serviceapi.Rule {
    if len(input) == 0 {
        return nil
    }

    rules := make([]serviceapi.Rule, 0)
    for _, item := range input {
        if item == nil {
            continue
        }

        raw := item.(map[string]interface{})
        rule := serviceapi.Rule{
            Name:    pointer.To(raw["name"].(string)),
            Enabled: pointer.To(raw["enabled"].(bool)),
        }

        rules = append(rules, rule)
    }

    return &rules
}
```

#### HashiCorp Standard Flatten Function Pattern

**Simple Configuration Blocks:**
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

**Array Configuration Blocks:**
```go
func flattenServiceRules(input *[]serviceapi.Rule) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    rules := make([]interface{}, 0)
    for _, rule := range *input {
        rules = append(rules, map[string]interface{}{
            "name":    pointer.From(rule.Name),
            "enabled": pointer.From(rule.Enabled),
        })
    }

    return rules
}
```

#### Expand/Flatten Best Practices

**DO:**
- Use `pointer.To()` and `pointer.From()` for consistent pointer handling
- Return `nil` for empty inputs in expand functions
- Return `make([]interface{}, 0)` for empty inputs in flatten functions
- Use clear, descriptive function names matching the schema field name
- Handle nil pointers safely in flatten functions

**DON'T:**
- Create multiple expand/flatten functions for the same logical block
- Use excessive variable assignments - directly construct structs when possible
- Add unnecessary comments explaining obvious operations
- Create complex conditional logic when simple struct initialization works

### Azure PATCH Operation Pattern

#### Critical PATCH Behavior Understanding

**Azure Resource Manager PATCH Operations:**
Many Azure services use PATCH operations for resource updates, which have fundamentally different behavior from PUT operations:

- **PATCH preserves existing values** when fields are omitted from the request
- **PUT replaces the entire resource** with the provided configuration
- **Azure SDK nil filtering** removes nil values before sending requests to Azure
- **Residual state persistence** means previously enabled features remain active unless explicitly disabled

**PATCH Operation Challenges:**
```go
// PROBLEM: This approach fails with PATCH operations
func ExpandPolicy(input []interface{}) *azuretype.Policy {
    if len(input) == 0 {
        return nil // SDK filters this out, Azure never gets disable command
    }

    // Only sends enabled policies, disabled policies remain unchanged
    result := &azuretype.Policy{}
    // Configure only enabled features...
    return result
}
```

**SOLUTION: Explicit Disable Pattern for PATCH Operations:**
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

**Key PATCH Operation Principles:**
- **Always return complete structures** - Never return nil for PATCH operations
- **Explicit disable commands** - Use `enabled=false` instead of omitting fields
- **Required field compliance** - Provide all required fields even when features are disabled
- **Residual state management** - Ensure previously enabled features can be properly disabled
- **"None" Pattern Integration** - Combine with flatten functions that exclude disabled features from state

**Documentation Requirements for PATCH Operations:**
```go
// PATCH Behavior Note: Azure VMSS uses PATCH operations which preserve existing values
// when fields are omitted. This means previously enabled policies will remain active
// unless explicitly disabled with enabled=false. Sending nil values gets filtered out
// by the Azure SDK, so Azure never receives disable commands. We must explicitly
// send enabled=false for all policies that should be disabled.
```

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

// Don't over-assign simple boolean values
enabledFromConfig := d.Get("enabled").(bool)
enabled := enabledFromConfig
```

#### Pointer Helper Usage Standards

**PREFERRED - Use Pointer Helpers:**
```go
properties := &serviceapi.Properties{
    Name:    pointer.To(name),
    Enabled: pointer.To(enabled),
    Count:   pointer.To(int64(count)),
}
```

**AVOID - Manual Pointer Creation:**
```go
namePtr := &name
enabledPtr := &enabled
properties := &serviceapi.Properties{
    Name:    namePtr,
    Enabled: enabledPtr,
}
```

### Comment Guidelines

#### Minimal Comment Standards

**Comments should ONLY be used for:**
- Complex business logic that isn't obvious from the code
- Azure API-specific requirements or constraints
- Workarounds for Azure SDK or API limitations
- Non-obvious state management patterns (like PATCH operations)

**DO NOT comment:**
- Obvious operations like variable assignments
- Standard expand/flatten operations
- Simple struct initialization
- Basic conditional logic

**GOOD Examples:**
```go
// Azure Front Door requires 16-240 second timeout range
if timeout < 16 || timeout > 240 {
    return fmt.Errorf("timeout must be between 16 and 240 seconds")
}

// PATCH operations require explicit disable commands to clear residual state
result := &azuretype.Policy{
    AutoFeature: &azuretype.AutoFeature{
        Enabled: pointer.To(false), // Explicit disable for PATCH
    },
}
```

**BAD Examples:**
```go
// Get the name from configuration
name := d.Get("name").(string)

// Create the configuration object
config := expandConfiguration(d.Get("configuration").([]interface{}))

// Set the properties
properties.Name = pointer.To(name)
```

### Testing Patterns

For comprehensive testing patterns, implementation details, and Azure-specific testing guidelines, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

Key testing requirements:
- Write comprehensive acceptance tests for all resources
- Use the standard acceptance test framework
- Test both success and failure scenarios
- Ensure tests are idempotent and can run in parallel
- Include import tests for all resources (`data.ImportStep()`)
- Test Azure-specific features like resource tagging and location handling

#### Test Template Organization Patterns

**Semantic Template Functions:**
```go
// Create semantic template functions that clearly indicate their behavior
func (r ServiceNameResource) templateWithForceDelete(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {
    virtual_machine_scale_set {
      force_delete = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-service-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

// Use semantic templates in test configurations
func (r ServiceNameResource) resiliencyTestWithCleanup(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

resource "azurerm_service_resource" "test" {
  # Resource configuration with protective features
  resiliency {
    enabled = true
  }
}
`, r.templateWithForceDelete(data))
}
```

**Template Naming Conventions:**
- `templateWithForceDelete`: Provider with aggressive cleanup flags
- `templateWithLocation`: Base resources in specific region
- `templateWithBackupRetention`: Resources with backup/retention policies
- `templatePublicKey`: Base template with SSH key generation

### Troubleshooting Azure Resource State Management

#### PATCH Operation + "None" Pattern Debugging

When implementing Azure resources that use PATCH operations combined with the "None" pattern, several common issues can arise. This section documents debugging approaches learned from real-world troubleshooting sessions.

**Common Symptoms:**
- Resource state shows fields as disabled, but Azure portal shows them as enabled
- Tests pass on creation but fail when testing disable  re-enable scenarios
- Azure API calls return success, but resource configuration doesn't change
- Residual state persists after removing Terraform configuration blocks

**Root Cause Analysis Framework:**

1. **Identify the HTTP Method**: Check if the Azure service uses PATCH vs PUT operations
   ```bash
   # Look for PatchThenPoll vs CreateOrUpdateThenPoll in Azure SDK calls
   grep -r "PatchThenPoll\|CreateOrUpdateThenPoll" internal/services/servicename/
   ```

2. **Trace Azure SDK Filtering**: Verify if nil values are being filtered out
   ```go
   // Look for patterns like this that cause issues:
   if len(input) == 0 {
       return nil // SDK filters this out, Azure never gets disable command
   }
   ```

3. **Check "None" Pattern Implementation**: Ensure disabled features are explicit
   ```go
   // WRONG - Causes residual state
   func ExpandFeature(input []interface{}) *azuretype.Feature {
       if len(input) == 0 {
           return nil
       }
       // Configure only enabled features
   }

   // RIGHT - Prevents residual state
   func ExpandFeature(input []interface{}) *azuretype.Feature {
       result := &azuretype.Feature{
           Enabled: pointer.To(false), // Explicit disable
       }
       if len(input) > 0 {
           result.Enabled = pointer.To(true)
       }
       return result
   }
   ```

**Testing Strategy for State Management Issues:**

```go
func TestAccServiceName_stateManagement(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_service_name", "test")
    r := ServiceNameResource{}

    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            // Step 1: Enable feature
            Config: r.featureEnabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).Key("feature.#").HasValue("1"),
                check.That(data.ResourceName).Key("feature.0.enabled").HasValue("true"),
            ),
        },
        {
            // Step 2: Disable feature (critical test for residual state)
            Config: r.featureDisabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                // Verify feature is removed from Terraform state ("None" pattern)
                check.That(data.ResourceName).Key("feature.#").HasValue("0"),
                // Custom check to verify Azure resource is actually disabled
                check.That(data.ResourceName).Key("id").Exists(),
            ),
        },
        {
            // Step 3: Re-enable feature (tests that disable actually worked)
            Config: r.featureEnabled(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).Key("feature.#").HasValue("1"),
                check.That(data.ResourceName).Key("feature.0.enabled").HasValue("true"),
            ),
        },
    })
}
```

**Documentation Patterns for Complex State Management:**

```go
func ExpandComplexFeature(input []interface{}) *azuretype.ComplexFeature {
    // State Management Note: This Azure service uses PATCH operations combined with the "None" pattern.
    // Previously enabled features will remain active unless explicitly disabled with enabled=false.
    // When users omit the configuration block, we must send disabled=true to clear residual state.

    // Always return complete structure with explicit disable commands
    result := &azuretype.ComplexFeature{
        AutoFeature: &azuretype.AutoFeature{
            Enabled: pointer.To(false),                           // Explicit disable
            RequiredField: pointer.To(azuretype.DefaultValue),    // Required even when disabled
        },
    }

    // Enable only explicitly configured features
    if len(input) > 0 && input[0] != nil {
        // User provided configuration - enable and apply settings
    }

    return result
}
```

**Key Debugging Insights:**
- **PATCH + nil filtering**: Azure SDK removes nil values, so disable commands never reach Azure
- **Residual state persistence**: Previously enabled features stay active without explicit disable
- **"None" pattern compliance**: Users expect omitted fields to disable features, not preserve them
- **Test coverage gaps**: Must test enable  disable  re-enable to catch state management bugs
- **Documentation clarity**: Comment the interaction between PATCH operations and "None" patterns

This debugging framework helps identify and resolve state management issues in Azure resources that combine PATCH operations with the "None" pattern for user-friendly configuration management.
