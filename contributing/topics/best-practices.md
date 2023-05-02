# Best Practices

Since it's inception the provider has undergone various iterations and changes in convention, as a result there can be legacy by-products within the provider which are inadvertently used as references. This section contains a miscellaneous assortment of current best practices to be aware of when contributing to the provider.

## Separate Create and Update Methods

Combined create and update methods within the provider are a legacy remnant that will need to be separated going forward and prior to the next major release of the provider due to changes in behaviour in Terraform core and the providers migration from `terraform-pluginsdk` to `terraform-plugin-framework`.

Migration to `terraform-plugin-framework` will see the use of `ignore_changes` become more prevalent in Terraform configs. For `ignore_changes` to work the create and update methods must be separate and `d.HasChanges` needs to be called on properties within the resource - see [this guide on new resources](guide-new-resource.md) for an example of how this is done.

## Typed vs. Untyped Resources

At this point in time the Provider supports Data Sources and Resources built using either the Typed SDK, or `hashicorp/terraform-plugin-sdk` (which we call `Untyped`). Whilst both of these output Terraform Data Sources and Resources, we're gradually moving from using Untyped Data Sources and Resources to Typed Resources since there's a number of advantages in doing so. We currently recommend using the [internal sdk package](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/sdk#should-i-use-this-package-to-build-resources) to build Typed Resources.

An example of both Typed and Untyped Resources can be found below - however as a general rule:

* When the Resource imports `"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"` - it's using the Typed SDK.
* When the Resource doesn't import `"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"` - then it's an Untyped Resource, which is backed by `hashicorp/terraform-plugin-sdk`.

Data Sources and Resources built using the Typed SDK have a number of benefits over those using `hashicorp/terraform-plugin-sdk` directly:

* The Typed SDK requires that a number of Azure specific behaviours are present in each Data Source/Resource. For example, the `interface` defining the Typed SDK includes a `IDValidationFunc()` function, which is used during `terraform import` to ensure the Resource ID being specified matches what we're expecting. Whilst this is possible using the Untyped SDK, it's more work to do so, as such using the Typed SDK ensures that these behaviours become common across the provider.
* The Typed SDK exposes an `Encode()` and `Decode()` method, allowing the marshalling/unmarshalling of the Terraform Configuration into a Go Object - which both:
    1. Avoids logic errors when an incorrect key is used in `d.Get` and `d.Set`, since we can (TODO: https://github.com/hashicorp/terraform-provider-azurerm/blob/5652afa601d33368ebefb4a549584e214e9729cb/internal/sdk/wrapper_validate.go#L21) validate that each of the HCL keys used for the models (to get and set these from the Terraform Config) is present within the Schema via a unit test, rather than failing during the `Read` function, which takes considerably longer.
    2. Default values can be implied for fields, rather than requiring an explicit `d.Set` in the Read function for every field - this allows us to ensure that an empty value/list is set for a field, rather than being `null` and thus unreferenceable in user configs.
* Using the Typed SDK allows Data Sources and Resources to (in the future) be migrated across to using `hashicorp/terraform-plugin-framework` rather than `hashicorp/terraform-plugin-sdk` without rewriting the resource - which will unlock a number of benefits to end-users, but does involve some configuration changes (and as such will need to be done in a major release).
* Using the Typed SDK means that these Data Sources/Resources can be more easily swapped out for generated versions down the line (since the code changes will be far smaller).
  
To facilitate the migration across to Typed Resources, we ask that any new Data Source or Resource which is added to the Provider is added as a Typed Data Source/Resource. Enhancements to existing Data Sources/Resources which are Untyped Resources can remain as Untyped Resources, however these will need to be migrated across in the future.

```go
package someservice

import ...

func someResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: someResourceCreate,
		Read:   someResourceRead,
		Update: someResourceUpdate,
		Delete: someResourceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := someresource.ParseSomeResourceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// schema fields are defined here
		},
	}
}

func someResourceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	// create logic is defined here
}

func someResourceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	// update logic is defined here
}

func someResourceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// read logic is defined here
}

func someResourceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// delete logic is defined here
}

```

Typed resources are initialised using interfaces and methods from the `sdk` package within the provider and will look something like the example below:

```go
package someservice

import ...

type SomeResource struct{}

var (
	_ sdk.Resource           = SomeResource{}
	_ sdk.ResourceWithUpdate = SomeResource{}
)

type SomeResourceModel struct {
	DisplayName           string            `tfschema:"display_name"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	Sku                   string            `tfschema:"sku_name"`
	Tags                  map[string]string `tfschema:"tags"`
	TenantId              string            `tfschema:"tenant_id"`
}

func (r SomeResource) ResourceType() string {
	return "azurerm_some_resource"
}

func (r SomeResource) ModelObject() interface{} {
	return &SomeResourceModel{}
}

func (r SomeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return someService.ValidateSomeResourceID
}

func (r SomeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// settable schema fields are set here
    }
}

func (r SomeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// read-only schema fields are set here
	}
}

func (r SomeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// create logic is defined here 
		},
	}
}

func (r SomeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// update logic is defined here
		},
	}
}

func (r SomeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// read logic is defined here
		},
	}
}

func (r SomeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func:    func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// delete logic is defined here
		},
	}
}
```


## Setting Properties to Optional + Computed

Due to certain behaviours in the Azure API it had become commonplace to set properties whose values could change in the background or for new properties that returned a default value not set by users as `Computed` to prevent the provider from flagging a diff after applying a plan.

Differing behaviour in the new protocol version (v6) used in `terraform-core` mean that any changes that occur to a property's value after applying will throw an error instead of a warning. Thus any uses of `Optional` + `Computed` should be avoided and existing ones will need to be phased out and replaced with logic that allows users to add the property to `ignore_changes`.
