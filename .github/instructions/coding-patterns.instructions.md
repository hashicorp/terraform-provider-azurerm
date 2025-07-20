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

            if model.Location != nil {
                state.Location = *model.Location
            }

            if model.Sku != nil {
                state.Sku = string(model.Sku.Name)
            }

            if model.Tags != nil {
                state.Tags = *model.Tags
            }

            if props := model.Properties; props != nil {
                if props.Enabled != nil {
                    state.Enabled = *props.Enabled
                }
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

            if model.Location != nil {
                state.Location = *model.Location
            }

            if model.Sku != nil {
                state.Sku = string(model.Sku.Name)
            }

            if model.Tags != nil {
                state.Tags = *model.Tags
            }

            if props := model.Properties; props != nil {
                if props.Enabled != nil {
                    state.Enabled = *props.Enabled
                }
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

### Testing Patterns

For comprehensive testing patterns, implementation details, and Azure-specific testing guidelines, see [`testing-guidelines.instructions.md`](./testing-guidelines.instructions.md).

Key testing requirements:
- Write comprehensive acceptance tests for all resources
- Use the standard acceptance test framework
- Test both success and failure scenarios
- Ensure tests are idempotent and can run in parallel
- Include import tests for all resources (`data.ImportStep()`)
- Test Azure-specific features like resource tagging and location handling
