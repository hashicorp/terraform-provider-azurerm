# Guide: Extending a Resource

As Azure services evolve and new features or functionalities enter public preview or become GA, their corresponding Terraform resources will need to be extended and/or modified in order to expose newly added functionalities.

Oftentimes this involves the addition of a new property or perhaps even the renaming of an existing property.

# Adding a new property

In order to incorporate a new property into a resource, modifications need to be made in multiple places. These are outlined below with pointers on what to consider and look out for as well as examples.

## Schema

Building on the example found in [adding a new resource](guide-new-resource.md) the new property will need to be added to either the user configurable list of `Arguments`, or `Attributes` if non-configurable.

Our hypothetical property `public_network_access_enabled` will be user configurable and thus will need to be added to the `Arguments` list.

The position of the new property is determined alphabetically and will end up looking like the code block below.

```go
func (ResourceGroupExampleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		
		"location": commonschema.Location(),
		
		"public_network_access_enabled": {
			Type: pluginsdk.TypeBool,
			Optional: true,
        }       

		"tags": commonschema.TagsDataSource(),
	}
}
```

* Remember to choose an appropriate name, see our [property naming guidelines](reference-naming.md).

* Ensure there is appropriate validation, at the very least `validation.StringIsNotEmpty` should be set for strings where a validation pattern cannot be determined.

## Create function

* The new property needs to be set in the properties struct for the resource.

```go
props := machinelearning.Workspace{
	WorkspaceProperties: &machinelearning.WorkspaceProperties{
		PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool))
    }
}
```

## Update function

* When performing selective updates check whether the property has changed.

```go
if d.HasChange("public_network_access_enabled") {
	props.WorkspaceProperties.PublicNetworkAccess = utils.Bool(d.Get("public_network_access_enabled").(bool))
}
```

## Read function

* Generally speaking all properties should have a value set into state.

* If the value returned by the API is a pointer we should nil check this to prevent panics in the provider.

```go
publicNetworkAccess := true
if v := props.WorkspaceProperties.PublicNetworkAccess; v != nil {
	publicNetworkAccess = *v
}
d.Set("public_network_access_enabled", publicNetworkAccess)
```

## Tests

* It is often sufficient to add a new property to one of the existing, non-basic tests.

* If the property is `Optional` then it can be added to the `complete` test.

* Properties that are `Required` will need to be added to all the existing tests for that resource.

* In cases where a new property or block requires additional setup or pre-requisites it makes sense to create a dedicated test for it.

## Docs

* Lastly, don't forget to update the docs with this new property!

* Property ordering within the docs follows the same conventions as in the Schema.

* `Computed` only values should be added under `Attributes Reference`

# Renaming and Deprecating a Property

Fixing typos in property names or renaming them to improve the meaning is unfortunately not just a matter of updating the name in the resource's code.

This is a breaking change and can be done by deprecating the old property, replacing it with a new one, as well as feature flagging its removal in the next major release of the provider.

A feature flag is essentially a function that returns a boolean and allows the provider to accommodate alternate behaviours that are meant for major releases. A release feature flag will always return `false` until it has been hooked up to an environment variable that allows users to toggle the behaviour and can be found in the `./internal/features` directory.

As an example, let's deprecate and replace the property `enable_public_network_access` with `public_network_access_enabled`.

```go
Schema: map[string]*pluginsdk.Schema{
    ...
    "enable_public_network_access": {
        Type:     pluginsdk.TypeBool,
        Optional: true,
    },
}
```

After deprecation the schema might look like the example below.

TODO: how is this done for typed resources?

```go
func resource() *pluginsdk.Resource {
    resource := &pluginsdk.Resource{
        ...
        Schema: map[string]*pluginsdk.Schema{
        	...

            // The deprecated property is moved out of the schema and conditionally added back via the feature flag
        	
            "public_network_access_enabled": {
                Type:       pluginsdk.TypeBool,
                Optional:   true,
                Computed:   !features.FourPointOhBeta()  // Conditionally computed to avoid diffs when both properties are set into state
                ConflictsWith: func() []string{          // Conditionally conflict with the deprecated property which will no longer exist in the next major release
                	if !features.FourPointOhBeta() {
                		return []string{"public_network_access_enabled"}
                    }      
                    return []string{}
                }(),       
            },
        }
    }
    
    if !features.FourPointOhBeta() {
    	resource.Schema["enable_public_network_access"] = &pluginsdk.Schema{
            Type:       pluginsdk.TypeBool,
            Optional:   true,
            Computed:   true,
            Deprecated: "This property has been renamed to `public_network_access_enabled` and will be removed in v4.0 of the provider",
            ConflictsWith: []string{"public_network_access_enabled"}
        }   
    }
    
    return resource
}
```

Also make sure to feature flag the behaviour in the create, update and read methods.

```go
func create() {
	...
	publicNetworkAccess := false
	if !features.FourPointOhBeta() {
		if v, ok := d.GetOkExists("enable_public_network_access"); ok {
			publicNetworkAccess = v.(bool)
		}       
    }
    
    if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
        publicNetworkAccess = v.(bool)
    }
	...
}

func read() {
	...
	d.Set("public_network_access_enabled", props.PublicNetworkAccess)
	
	if !features.FourPointOhBeta() {
		d.Set("enable_public_network_access", props.PublicNetworkAccess)
    }   
	...
}
```
