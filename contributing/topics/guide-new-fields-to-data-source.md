# Guide: Extending a Data Source

It is sometimes necessary to make changes to an existing Data Source. Reasons include:

* A new property is added in the referenced resource

* A property is deprecated and/or no longer available in the referenced resource

* An API update changes the behaviour of the referenced resource

When updating an existing Data Source keep in mind the configuration of the end user that may be using it.  Mitigations must be taken, where possible, to prevent the change breaking existing user configurations.

The process is similar to [extending an existing Resource](guide-new-fields-to-resource), in that modifications in multiple places are required.

## Schema

Building on the example from [adding a new data source](guide-new-data-source.md) the new property will need to be added into the `Attributes` list which contains a list of schema fields that are Computed only.

The location of the new property within this list is determined alphabetically. Taking the hypothetical property `public_network_access_enabled` as an example this would then end up looking like this in `Attributes`.

```go
func (ResourceGroupExampleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
		},
		
		"public_network_access_enabled": {
			Type: pluginsdk.TypeBool,
			Computed: true,
        },       

		"tags": commonschema.TagsDataSource(),
	}
}
```

* For new properties that are ambiguous in their functionality or nature, follow the [property naming guidelines](reference-naming.md) when choosing a name.

## Read function

* The only thing to consider here is setting the new value into state and nil checking beforehand if required.

```go
publicNetworkAccess := true
if v := props.WorkspaceProperties.PublicNetworkAccess; v != nil {
	publicNetworkAccess = *v
}
d.Set("public_network_access_enabled", publicNetworkAccess)
```

## Tests

* New properties should be added to the basic data source test with an explicit check.

```go
func TestAccDataSourceSomeResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_some_resource", "test")
	r := AvailabilitySetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("new_property").HasValue("1"),
			),
		},
	})
}
```

## Docs

* Lastly, don't forget to update the docs where the property ordering is determined alphabetically.
