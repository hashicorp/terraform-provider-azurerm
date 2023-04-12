# Best Practices

Since it's inception the provider has undergone various iterations and changes in convention, as a result there can be legacy by-products within the provider which are inadvertently used as references. This section contains a miscellaneous assortment of current best practices to be aware of when contributing to the provider.

## Separate Create and Update Methods

Combined create and update methods within the provider are a legacy remnant that will need to be separated going forward and prior to the next major release of the provider due to changes in behaviour in Terraform core and the providers migration from `terraform-pluginsdk` to `terraform-plugin-framework`.

Migration to `terraform-plugin-framework` will see the use of `ignore_changes` become more prevalent in Terraform configs. For `ignore_changes` to work the create and update methods must be separate and `d.HasChanges` needs to be called on properties within the resource - see [this guide on new resources](guide-new-resource.md) for an example of how this is done.

## Typed vs. Untyped Resources

To facilitate migration over to `terraform-plugin-framework` and to ease the transition towards generated resources we ask that all new resources that are added to the provider be typed instead of the former untyped resources.

What is meant by this is best explained by way of an example of each type.

Untyped resources rely on interfaces and methods within the `pluginsdk` package to build the resource - they look roughly like this:

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
			_, err := parse.SomeResourceID(id)
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

Differing behaviour in the new protocol version (v6) used in `terraform-core` mean that any changes that occur to a properties value after applying will throw an error instead of a warning. Thus any uses of `Optional` + `Computed` should be avoided and existing ones will need to be phased out and replaced with logic that allows users to add the property to `ignore_changes`.





