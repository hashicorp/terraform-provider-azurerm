# Best Practices

Since its inception, the provider has undergone various iterations and changes in convention, as a result there can be legacy by-products within the provider which are inadvertently used as references. This section contains a miscellaneous assortment of current best practices to be aware of when contributing to the provider.

## Separate Create and Update Methods

Historically the Provider has opted to combine the Create and Update methods due to the behaviour of the Azure API, where the same API is used for both Create and Update, meaning that the same payload has to be sent during both the Creation and Update of the resource.

In order to properly support Terraform's `ignore_changes` feature, rather than using a combined method for Create and Update, we're now requiring that these be separate, and that in the Update partial/delta differences are performed, to only update the value for a field if it's marked as changed.

For example, whilst a Create method may look similar to below:

```go
payload := resources.Group{
    Location: location.Normalize(d.Get("location").(string)),
    Tags: tags.Expand(d.Get("tags").(map[string]interface{}),
}

if err := client.CreateThenPoll(ctx, id, payload); err != nil {
  return fmt.Errorf("creating %s: %+v", id, err)
}
```

The update method should be checking if the updatable fields (in this example, only tags) - have changes (using `d.HasChanges` - which will flag updated values in the config if they're not ignored via `ignore_changes`).

Depending on the API there are two types of Updates, a patch/delta update (where only the fields containing changes are sent) - and a full update (which requires sending the full payload) - these are differentiable via the method name in the SDK, patch/delta updates are generally called `Update`, with a full update being called `CreateOrUpdate`.

A patch/delta update would look similar to below:

```go
payload := resources.GroupUpdate{}
if d.HasChanges("tags") {
  // this uses `pointer.To` since all fields are optional in a patch/delta update, so they'll only be updated if specified
  payload.Tags = pointer.To(tags.Expand(d.Get("tags").(map[string]interface{}))
}

if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
  return fmt.Errorf("updating %s: %+v", id, err)
}
```

A full update would retrieve the existing object from the API and then patch it, for example:

```go
resp, err := client.Get(ctx, id)
if err != nil {
  return fmt.Errorf("retrieving %s: %+v", id, err)
}

if resp.Model == nil {
  return fmt.Errorf("retrieving %s: model was nil", id)
}

payload := *resp.Model
if d.HasChanges("tags") {
  payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{})
}

if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
  return fmt.Errorf("updating %s: %+v", id, err)
}
```

## Typed vs. Untyped Resources

At this point in time the Provider supports Data Sources and Resources built using either the Typed SDK, or `hashicorp/terraform-plugin-sdk` (which we call `Untyped`). Whilst both of these output Terraform Data Sources and Resources, we're gradually moving from using Untyped Data Sources and Resources to Typed Resources since there's a number of advantages in doing so. We currently recommend using the [internal sdk package](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/sdk#should-i-use-this-package-to-build-resources) to build Typed Resources.

An example of both Typed and Untyped Resources can be found below - however as a general rule:

* When the Resource imports `"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"` - it's using the Typed SDK.
* When the Resource doesn't import `"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"` - then it's an Untyped Resource, which is backed by `hashicorp/terraform-plugin-sdk`.

Data Sources and Resources built using the Typed SDK have a number of benefits over those using `hashicorp/terraform-plugin-sdk` directly:

* The Typed SDK requires that a number of Azure specific behaviours are present in each Data Source/Resource. For example, the `interface` defining the Typed SDK includes an `IDValidationFunc()` function, which is used during `terraform import` to ensure the Resource ID being specified matches what we're expecting. Whilst this is possible using the Untyped SDK, it's more work to do so, as such using the Typed SDK ensures that these behaviours become common across the provider.
* The Typed SDK exposes an `Encode()` and `Decode()` method, allowing the marshalling/unmarshalling of the Terraform Configuration into a Go Object - which both:
    1. Avoids logic errors when an incorrect key is used in `d.Get` and `d.Set`, since we can validate that each of the HCL keys used for the models (to get and set these from the Terraform Config) is present within the Schema via a unit test, rather than failing during the `Read` function, which takes considerably longer.
    2. Default values can be implied for fields, rather than requiring an explicit `d.Set` in the Read function for every field - this allows us to ensure that an empty value/list is set for a field, rather than being `null` and thus not able to be referenced in user configs.
* Using the Typed SDK allows Data Sources and Resources to (in the future) be migrated across to using `hashicorp/terraform-plugin-framework` rather than `hashicorp/terraform-plugin-sdk` without rewriting the resource - which will unlock a number of benefits to end-users, but does involve some configuration changes (and as such will need to be done in a major release).
* Using the Typed SDK means that these Data Sources/Resources can be more easily swapped out for generated versions down the line (since the code changes will be far smaller).
  
To facilitate the migration across to Typed Resources, we ask that any new Data Source or Resource which is added to the Provider is added as a Typed Data Source/Resource. Enhancements to existing Data Sources/Resources which are Untyped Resources can remain as Untyped Resources, however these will need to be migrated across in the future.

Here is an example of an Untyped Resource:

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

There are many APIs within Azure that will specify a default value for a field if one isn't specified, for example the `createMode` field is typically defaulted (server-side) to `Default`.

The Azure Provider currently makes use of `hashicorp/terraform-plugin-sdk@v2` to define Data Sources and Resources, which under the hood uses v5 of the Terraform Protocol to interact with Terraform Core.

In version 5 of the Terraform Protocol, if a field is created with one value at Create time and returns a different value immediately after creation, then an internal warning is logged (but no error is raised) - meaning that the only way this change is visible is through a diff when `terraform plan` is run. The next version of the Terraform Protocol (v6 - used by `hashicorp/terraform-plugin-framework`) changes this from a logged warning to an error at runtime - meaning that these diffs will become more visible to users (and need to be accounted for in the provider).

To work around situations where we need to expose the default value from the Azure API - we've historically marked fields as both `Optional` and `Computed` - meaning that a value will be returned from the API when it's not defined.

Whilst this works, a side effect is that it's hard for users to reset a field to its default value when this is done - as such some fields today (such as the subnets block within the azurerm_virtual_network resource) require that an explicit empty list is specified (for example `subnets = []`) to remove this value, where this field is `Optional` and `Computed`.

In order to solve this, (new) fields should no longer be marked as `Optional` and `Computed` - instead where a split Create and Update method is used (see above) users can lean on `ignore_changes` to ignore values from a field with a default value, should they wish to continue using the default value.

This approach means that we can support users who want to use the default value by specifying `ignore_changes = ["some_field"]`, users who want to explicitly define this value e.g. some_field = "bar" and users who need to remove this value by either omitting the field or defining it as null, so that it gets removed.

Over time, the existing resources will be migrated from `Optional` and `Computed` to just `Optional` (allowing users to rely on `ignore_changes`) so that this becomes more behaviourally consistent - however new fields should be defined as `Optional` alone, rather than `Optional` and `Computed`.

If you encounter a field that must be `Optional` and `Computed`, make sure it follows the following conventions:
* The properties are in this sequence: Optional, Explanatory Comment, Computed
* The comment should start with `// NOTE: O+C `, and then explain the reason for the field being `Optional` and `Computed`

Example:

```go
	"etag": {
		Type: pluginsdk.TypeString,
		Optional: true,
		// NOTE: O+C Azure generates a new value every time this resource is updated
		Computed: true,
	},
```
