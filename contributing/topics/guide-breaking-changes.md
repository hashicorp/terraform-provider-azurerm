# Guide: Breaking Changes and Deprecations

To keep up with and accommodate the changing pace of Azure, the provider needs to be able to gracefully introduce and handle breaking changes. A "breaking change" within the provider is considered to be anything that requires an end user to modify previously valid terraform configuration after a provider upgrade to either deploy new resources or to maintain existing deployments. Even if a change does not affect the user's current deployment, it is still considered a breaking change if it requires the user to modify their configuration to deploy new resources. 

The `azurerm` provider attempts to be as "surface stable" as possible during minor and patch releases meaning breaking changes are typically only made during major releases, however exceptions are sometimes made for minor releases when the breaking change is deemed necessary or is unavoidable. Terraform users rely on the stability of Terraform providers as not only can configuration changes be costly to make, test, and deploy they can also affect downstream tooling such as modules. Even as part of a major release, breaking changes that are overly large or have little benefit can delay users upgrading to the next major version.

Generally we can safely introduce breaking changes into the provider for the major release using a feature flag. For the next major release that would be the `features.FivePointOh()` flag which is available in the provider today. This guide includes several topics on how to do common deprecations and breaking changes in the provider using this feature flag, as well as additional guidance on how to deal with changing default values in the Azure API. 

Types of breaking changes covered are:

- [Removing Resources or Data Sources](#removing-resources-or-data-sources)
- [Breaking Schema Changes](#breaking-schema-changes-and-deprecations)
- [Updating Default Values](#updating-default-values)
- [Post Release Breaking Change Clean Up](#post-release-breaking-change-clean-up)

## Removing Resources or Data Sources

Resources can be removed for several reasons, the service could be retiring, the API may no longer support creation of that resource or the resource has been renamed or superseded by a new version.

In all cases the resources cannot be removed from the provider in a minor release but must be deprecated and the registration of the resource made conditional using the major release feature flag.

The steps outlined below uses an example resource that is deprecated, but the same principles and steps apply for data sources as well.

1. Add the appropriate deprecation message to the resource.
    
   For Typed Resources
    ```go
    
    // For resources that have no replacement
    
    var _ sdk.ResourceWithDeprecationAndNoReplacement = ResourceWithNoReplacement{}
    
    func (r ResourceWithNoReplacement) DeprecationMessage() string {
        return "The `azurerm_resource_with_no_replacement` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider"
    }
     
    
    // For resources that have a replacement
    
    var _ sdk.ResourceWithDeprecationReplacedBy = ResourceWithReplacement{}
    
    func (r ResourceWithReplacement) DeprecatedInFavourOfResource() string {
        return "azurerm_new_resource"
    }
    
    ```

    For Untyped Resources
    ```go
    func resourceExample() *pluginsdk.Resource {
        return &pluginsdk.Resource{
            Create: resourceExampleCreate,
            Read:   resourceExampleRead,
            Update: resourceExampleUpdate,
            Delete: resourceExampleDelete,
            
            Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
            },
            
            DeprecationMessage: "The `azurerm_example` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider"
            ...
        }
    }
    ```

2. Conditionally register the resource in the `registration.go` file of the service package.
    
   For Typed Resources
    ```go
    func (r Registration) Resources() []sdk.Resource {
        resources := []sdk.Resource{
            MySqlFlexibleServerResource{},
        }
        
        if !features.FivePointOh() {
            resources = append(resources, ExampleResource{})
        }
        
        return resources
    }
    ```
    
    For Untyped Resources
    ```go
    func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
        resources := map[string]*pluginsdk.Resource{
            "azurerm_mysql_flexible_server": resourceMysqlFlexibleServer(),
    
        }
    
        if !features.FivePointOh() {
            resources["azurerm_example"] = resourceExample()
        }
    
        return resources
    }
    ```

3. Skip all tests related to the deprecated resource.

    ```go
    func TestAccExample_basic(t *testing.T) {
        if features.FivePointOh() {
            t.Skipf("Skipping since `azurerm_example` is deprecated and will be removed in 5.0")
        }
        data := acceptance.BuildTestData(t, "azurerm_example", "test")
        r := ExampleResource{}
        
        data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
            data.ImportStep()
        })
    }
    ```

4. Update the upgrade guide under `website/docs/5.0-upgrade-guide.markdown`.

   ```markdown
   ## Removed Resources
   
   ### `azurerm_example`
   
   This deprecated resources has been removed from the Azure Provider.
   ```

5. Update the resource (or data source) documentation

   ```markdown
   ~> **Note:** The `azurerm_example` resource has been deprecated because [reason here e.g. the service is retiring by 2025-10-10] and will be removed in v5.0 of the AzureRM Provider.
   ```

## Breaking Schema Changes and Deprecations

Breaking schema changes can include:
- Property renames
- When properties become Required
- When properties have Computed removed and need to be added to `ignore_changes` to prevent diffs
- Changes to the validation e.g. the validation becomes more restrictive
- Changing the default value
- Changing the type

In all cases the deprecation is handled the same way and will be illustrated by the example below.

The following example follows a fictional resource that will have the following breaking changes made:
- The property `enable_scaling` renamed to `scaling_enabled`
- The property `version` has its default changed from `1` to `2`

1. Update the Schema with the target or desired breaking schema change and patch over the breaking schema change with the current behaviour using the major release feature flag.

   ```go
   func (r ExampleResource) Arguments() map[string]*pluginsdk.Schema{
      args := map[string]*pluginsdk.Schema{
         "scaling_enabled": {
            Type:     pluginsdk.TypeBool,
            Optional: true,
            Default: false,
         },      
         "version": {
            Type:     pluginsdk.TypeString,
            Optional: true,
            Default: 2,
         },
      }
   
      // Regardless of the number of arguments changing, the whole schema definition should be updated like the following rather than inline changes for the current schema definition.
	  // This is to make cleanup easy so we can delete this block when the next major version releases.
      if !features.FivePointOh() {
         args["enable_scaling"] = &pluginsdk.Schema{
            Type:          pluginsdk.TypeBool,
            Optional:      true,
            Computed:      true,
            Default:       false,
            ConflictsWith: []string{"scaling_enabled"},
            Deprecated:    "`enable_scaling` has been deprecated in favour of `scaling_enabled` and will be removed in v5.0 of the AzureRM Provider",
         }
         // When renaming a property both properties need to have `Computed` set on them until the old property is removed in the next major release
         // We also need to remember to set ConflictsWith on both the old and the renamed property to ensure users don't set both in their config
         args["scaling_enabled"] = &pluginsdk.Schema{
            Type:          pluginsdk.TypeBool,
            Optional:      true,
            Computed:      true,
            Default:       false,
            ConflictsWith: []string{"enable_scaling"},
         }
         
         args["version"].Default = 1
      }
   
      return args
   }
   ```
   > **Note:** In the past we've accepted in-lined anonymous functions in a property's schema definition to conditionally change the default value, validation function etc. these will no longer be accepted in the provider. This is a deliberate decision to reduce the variation in how deprecations are done in the provider and also simplifies the clean-up effort of feature flagged code after the major release.

2. Update the Create/Read/Update methods if necessary.

3. Update the test configurations.

   Here are some guidelines on what good testing coverage for renamed properties looks like:
   * All test configurations that reference the old property should be updated to use the renamed property
   * One test configuration should continue using the old property to ensure that it still works as expected, but switch to using the renamed property in the major release mode. An example of what that looks like is provided below.

   ```go
   func (ExampleResource) complete(data acceptance.TestData) string {
   if !features.FivePointOh() {
        return fmt.Sprintf(`
   provider "azurerm" {
     features {}
   }
   
   resource "azurerm_resource_group" "test" {
     name     = "acctestRG-example-%[1]d"
     location = "%[2]s"
   }
   
   resource "azurerm_example" "test" {
     name           = "acctestexample%[1]d"
     enable_scaling = true
   }
   `, data.RandomInteger, data.Locations.Primary)
        }
   return fmt.Sprintf(`
   provider "azurerm" {
     features {}
   }
   
   resource "azurerm_resource_group" "test" {
     name     = "acctestRG-example-%[1]d"
     location = "%[2]s"
   }
   
   resource "azurerm_example" "test" {
     name            = "acctestexample%[1]d"
     scaling_enabled = true
   }
   `, data.RandomInteger, data.Locations.Primary)
   }
   ```
   > **Note:** Wherever possible, only update the test configuration and avoid updating the test case since changes to the test cases are more involved and higher effort to clean up.

4. Update the upgrade guide under `website/docs/5.0-upgrade-guide.markdown`
   
   Under the appropriate section of the upgrade guide, add a line for the deprecation
   ```markdown
   ## Breaking changes in Resources
   
   ### `azurerm_example_resource`
   
   * The deprecated `enable_scaling` property has been removed in favour of the `scaling_enabled` property.
   * The property `version` now defaults to `2`.
   ```
   
   The resources/data sources should be added in alphabetical order.
   
5. Update the resource documentation

   * The resource documentation should only be updated when a property is undergoing a soft deprecation. In the example above the only update to the resource documentation we need to do is to remove the property `enable_scaling` and add the property `scaling_enabled`.

   * Breaking changes such as the default value changing, or other property behaviour changing in a way that will only be active when the major release has gone out *should not* be added to the documentation since these do not apply yet. Please do not add any `**Note:** This property will do x in 5.0` notes in the documentation. 

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

This is a breaking change as Terraform should not trigger a plan between minor version upgrades. Instead, what we can do is use the major release feature flag as shown in the example below or mark the field as Required if that default value is going to continue to change in the future:

```go
func (r SparkResource) Arguments() map[string]*pluginsdk.Schema{
    args := map[string]*pluginsdk.Schema{
        "spark_version": {
            Type:     pluginsdk.TypeString,
            Optional: true,
            Default: "3.4",
            ValidateFunc: validation.StringInSlice([]string{
               "2.4",
               "3.1",
               "3.2",
               "3.3",
               "3.4",
                }, false),
            },
        }

    if !features.FivePointOh() {
        args["spark_version"].Default = "2.4"
    }
	
    return args
}
```

## Adding a new property with a default value

When adding a new property with a default value, we can introduce a similar breaking change as the one noted above, but it's even harder to pinpoint. Take for example the following property recently added to `azurerm_kusto_account`:
                             
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

## Post Release Breaking Change Clean Up

Once the next major release has happened, all blocks of code that were conditionally included for that version (e.g. `if !features.FivePointOh() { ... }`) need to be removed. Most should be fine to simply remove, however there are a few things to watch out for:

1. For typed resources, if you are removing a property, make sure you also remove it from the model(s). The fields should have a `removedInNextMajorVersion` tag. 
2. For typed resources, there may be properties that were only included once the major version was released, make sure you remove the `addedInNextMajorVersion` tag from these properties in the model(s).
3. Confirm the documentation is up-to-date with what is in code, generally this should already be the case, but it's good to double-check.