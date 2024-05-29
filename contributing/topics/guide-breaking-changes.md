# Guide: Breaking Changes

Over time, new and existing properties will need either a Default Value added or updated. This can be seen as a quick fix but could expose a breaking change that won't be caught by tests so doing a more in depth look into how Default values can impact Terraform is key.

## Updating Default Values

There are some cases where Azure updates the default value for an attribute when creating a new resource, and we would want to do the same for the provider but this is an easy breaking change to miss.

We have a property like the following and Azure added a new spark version `3.4` and said that all new resources being created will be created with `3.4` as the default. 

In Terraform, we start with:

```hcl
    "spark_version": {
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  "2.4",
		ValidateFunc: validation.StringInSlice([]string{
			"2.4",
			"3.1",
			"3.2",
			"3.3",
		}, false),
	},
```

Then we would want to update `ValidateFunc` to include the new accepted value and update `Default` to `3.4` to keep it in line with Azure like so:

```hcl
    "spark_version": {
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  "3.4",
		ValidateFunc: validation.StringInSlice([]string{
			"2.4",
			"3.1",
			"3.2",
			"3.3",
			"3.4",
		}, false),
	},
```

But if we do that, people who have created that resource without the attribute specified will see a plan diff when upgrading to this version of the provider like so:

This config does not specify `spark_version` because we know we can rely on the default to fill it for us

```hcl
resource "azurerm_synapse_spark_pool" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
}
```

Running `terraform show` we can see `spark_version` has been filled in with the default of `2.4`

```hcl
# azurerm_synapse_spark_pool.example:
resource "azurerm_synapse_spark_pool" "example" {
    name                                = "example"
.
.
.
    spark_version                       = "2.4"
}
```

When running the version of the provider where the default has changed from `2.4` to `3.4`, we'll see the following plan:

```hcl
Terraform will perform the following actions:

  # azurerm_synapse_spark_pool.example will be updated in-place
  ~ resource "azurerm_synapse_spark_pool" "test" {
        id                                  = "exampleid"
        name                                = "example"
      ~ spark_version                       = "2.4" -> "3.4"
        tags                                = {}
        # (12 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

This is a breaking change as Terraform should not trigger a plan between minor version upgrades. Instead, what we can do is add a TODO next to the `Default` tag to update the default value in the next major version of the provider or mark the field as Required if that default value is going to continue to change in the future:

```hcl
    "spark_version": {
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default: func() string {
					if !features.FourPointOh() {
						return "2.4"
					}
					return "3.4"
				}(),
		ValidateFunc: validation.StringInSlice([]string{
			"2.4",
			"3.1",
			"3.2",
			"3.3",
			"3.4", 
		}, false),
	},
```

## Adding a new property with a default value

When adding a new property with a default value, we can introduce a similar breaking change as the one noted above but it's even harder to pinpoint. Take for example the following property recently added to `azurerm_kusto_account`:
                             
It originally came in like this:

```hcl
"auto_stop_enabled": {
	Type:     pluginsdk.TypeBool,
	Optional: true,
},
```

Our tests were failing because the Azure API was returning this value as true while Terraform does not expect this value to be set because it isn't specified in the config file. To fix this breaking change, we need to add a Default like so:

```hcl
"auto_stop_enabled": {
	Type:     pluginsdk.TypeBool,
	Optional: true,
	Default:  true,
},
```

There are many ways to accidentally add a breaking change when looking at properties with a Default or lack thereof so extra work needs to be done to confirm what Terraform and the Azure API are returning before deciding how best to incorporate the Default tag.
