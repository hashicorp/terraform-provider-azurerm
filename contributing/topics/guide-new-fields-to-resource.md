# Guide: Extending a Resource

As Azure services evolve and new features or functionalities enter public preview or become GA, their corresponding Terraform resources will need to be extended and/or modified in order to expose newly added functionalities.

Oftentimes this involves the addition of a new property or perhaps even the renaming of an existing property.

## Adding a new property

In order to incorporate a new property into a resource, modifications need to be made in multiple places. These are outlined below with pointers on what to consider and look out for as well as examples.

### Schema

Building on the example found in [adding a new resource](guide-new-resource.md) the new property will need to be added to either the user configurable list of `Arguments`, or `Attributes` if non-configurable.

Our hypothetical property `logging_enabled` will be user configurable and thus will need to be added to the `Arguments` list.

The position of the new property is determined based on the order found in [adding a new resource](guide-new-resource.md#step-3-scaffold-an-emptynew-resource) and will end up looking like the code block below. Here is an example for a typed resource:

```go
func (ResourceGroupExampleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		
		"location": commonschema.Location(),
		
		"logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		}

		"tags": commonschema.TagsDataSource(),
	}
}
```

* Remember to choose an appropriate name, see our [property naming guidelines](reference-naming.md).

* Ensure there is appropriate validation, at the very least `validation.StringIsNotEmpty` should be set for strings where a validation pattern cannot be determined.

* When adding multiple properties or blocks thought should be given on how to map these, see [schema design considerations](schema-design-considerations.md) for specific examples. 

### Create function

* The new property needs to be set in the properties struct for the resource.

```go
props := machinelearning.Workspace{
	Properties: &machinelearning.WorkspaceProperties{
		LoggingEnabled: pointer.To(model.LoggingEnabled)
	}
}
```

### Update function

* When performing selective updates check whether the property has changed.

```go
if metadata.ResourceData.HasChange("logging_enabled") {
	existing.Model.Properties.LoggingEnabled = pointer.From(model.LoggingEnabled)
}
```

### Read function

* Generally speaking all properties should have a value set into state.

* If the value returned by the API is a pointer we should account for the possibility of a nil reference to prevent panics in the provider. One way to do this is to use `pointer.From()`:

```go
state := MyResourceModel{}

if model := resp.Model; model != nil {
    state.LoggingEnabled = pointer.From(model.LoggingEnabled)
}

return metadata.Encode(&state)
```

### Tests

* It is often sufficient to add a new property to one of the existing, non-basic tests.

* If the property is `Optional` then it can be added to the `complete` test.

* Properties that are `Required` will need to be added to all the existing tests for that resource.

* In cases where a new property or block requires additional setup or pre-requisites it makes sense to create a dedicated test for it.

* Adding a new property or updating an existing property with a default should be checked properly against the current version of Terraform State and the Azure API as tests won't be able to catch a potential breaking change, see our [section on defaults and breaking changes](guide-breaking-changes.md)

### Docs

* Lastly, don't forget to update the docs with this new property!

* Property ordering within the docs follows the same conventions as in the Schema.

* `Computed` only values should be added under `Attributes Reference`

## Renaming and Deprecating a Property

Fixing typos in property names or renaming them to improve the meaning is unfortunately not just a matter of updating the name in the resource's code.

This is a breaking change and can be done by deprecating the old property, replacing it with a new one, as well as feature flagging its removal in the next major release of the provider.

A feature flag is essentially a function that returns a boolean and allows the provider to accommodate alternate behaviours that are meant for major releases. A release feature flag will always return `false` until it has been hooked up to an environment variable that allows users to toggle the behaviour and can be found in the `./internal/features` directory.

As an example, let's deprecate and replace the property `enable_compression` with `compression_enabled`. Here is an example for an untyped resource (for more information about typed and untyped resource, [see the Best Practices guide](./best-practices.md#typed-vs-untyped-resources)):

```go
Schema: map[string]*pluginsdk.Schema{
	...
	"enable_compression": {
		Type:     pluginsdk.TypeBool,
		Optional: true,
},
```

Here is an example for a typed resource:

```go
func (r ExampleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"enable_compression": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
}
```

After deprecation the schema might look like the code below. Here is an example for an untyped resource:

```go
func resource() *pluginsdk.Resource {
    resource := &pluginsdk.Resource{
        Schema: map[string]*pluginsdk.Schema{
            // The deprecated property is moved out of the schema and conditionally added back via the feature flag
            "compression_enabled": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
            },
        },
    }

    if !features.FivePointOh() {
        resource["compression_enabled"] = &pluginsdk.Schema{
            Type:          pluginsdk.TypeBool,
            Optional:      true,
            Computed:      true,
            ConflictsWith: []string{"enable_compression"}
        }
        
        resource["enable_compression"] = &pluginsdk.Schema{
            Type:	   pluginsdk.TypeBool,
            Optional:   true,
            Computed:   true,
            Deprecated: "This property has been renamed to `compression_enabled` and will be removed in v5.0 of the provider",
            ConflictsWith: []string{"compression_enabled"}
        }   
    }

    return resource
}
```

Here is an example for a typed resource:

```go
func (r ExampleResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
			// The deprecated property is moved out of the schema and conditionally added back via the feature flag
			"compression_enabled": {
				Type:	   pluginsdk.TypeBool,
				Optional:   true,
			},
		}
	}
	
	if !features.FivePointOh() {
		schema["compression_enabled"] = &pluginsdk.Schema{
			Type:	       pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"enable_compression"}
		}

		schema["enable_compression"] = &pluginsdk.Schema{
			Type:	   pluginsdk.TypeBool,
			Optional:   true,
			Computed:   true,
			Deprecated: "This property has been renamed to `compression_enabled` and will be removed in v5.0 of the provider",
			ConflictsWith: []string{"compression_enabled"}
		}   
	}
	
	return schema
}
```

Also make sure to feature flag the behaviour in the `Create()`, `Update()` and `Read()` methods. Here is an example for an untyped resource:

```go
func myResourceCreate() {
    ...
    enableCompression := false
    if !features.FivePointOh() {
        if v, ok := d.GetOkExists("enable_compression"); ok {
            enableCompression = v.(bool)
        }       
    }
    
    if v, ok := d.GetOkExists("compression_enabled"); ok {
        enableCompression = v.(bool)
    }
}

func myResourceRead() {
    ...
	d.Set("compression_enabled", props.EnableCompression)
    
    if !features.FivePointOh() {
        d.Set("enable_compression", props.EnableCompression)
    }
    ...
}
```

Here is an example for a typed resource:
```go
func (r ExampleResource) Create() sdk.ResourceFunc {
	...
	compressionEnabled := false
	if !features.FivePointOh() {
		compressionEnabled = model.EnableCompression
	}
	
	compressionEnabled = model.CompressionEnabled
	...
}

func (r ExampleResource) Read() sdk.ResourceFunc {
	...
	state.CompressionEnabled = pointer.From(props.CompressionEnabled)
	
	if !features.FivePointOh() {
		state.EnableCompression = pointer.From(props.CompressionEnabled)
	}   
	...
}
```

When deprecating a property in a Typed Resource, it is important to ensure that the Go struct representing the schema is correctly tagged to prevent the SDK decoding the removed property when the major version beta / feature flag is in use. In these cases the struct tags must be updated to include `,removedInNextMajorVersion`.  

```go
type ExampleResourceModel struct {
	Name               string `tfschema:"name"`
	EnableCompression  bool `tfschema:"enable_compression,removedInNextMajorVersion"`
	CompressionEnabled bool `tfschema:"compression_enabled"`
}
```
