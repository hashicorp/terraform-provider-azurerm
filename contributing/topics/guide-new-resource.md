# Guide: New Resource

This guide covers adding a new Resource to a Service Package, see [adding a New Service Package](guide-new-service-package.md) if the Service Package doesn't exist yet.

### Related Topics

* [Acceptance Testing](reference-acceptance-testing.md)
* [Our Recommendations for opening a Pull Request](guide-opening-a-pr.md)

### Stages

At this point in time the AzureRM Provider supports both Typed and Untyped Resources - more information can be found [in the High Level Overview](high-level-overview.md).

This guide covers adding a new Typed Resource, which makes uses the Typed SDK within this repository, which requires the following steps:

1. Ensure all the dependencies are installed (see [Building the Provider](building-the-provider.md)).
2. Add an SDK Client (if required).
3. Define the Resource ID.
4. Scaffold an empty/new Resource.
5. Register the new Resource.
6. Add Acceptance Test(s) for this Resource.
7. Run the Acceptance Test(s).
8. Add Documentation for this Resource.
9. Send the Pull Request.

We'll go through each of those steps in turn, presuming that we're creating a Resource for a Resource Group.

### Step 1: Ensure the Tools are installed

See [Building the Provider](building-the-provider.md).

### Step 2: Add an SDK Client (if required)

This section covers how to add and configure the SDK Client.

Determining which SDK Client you should be using is a little complicated unfortunately.

The Client for the Service Package can be found in `./internal/services/{name}/client/client.go` - and we can add an instance of the SDK Client we want to use (here `resources.GroupsClient`) and configure it (adding credentials etc):

```go
package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GroupsClient *resources.GroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	groupsClient := resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)
	
	// ...
	
	return &Client{
		GroupsClient: &groupsClient,
	}
}
```

A few things of note here:

1. The field `GroupsClient` within the struct is a pointer, meaning that if it's not initialized the Provider will crash/panic - which is intentional to avoid using an unconfigured client (which will have no credentials, and cause misleading errors).
2. When creating the client, note that we're using `NewGroupsClientWithBaseURI` (and not `NewGroupsClient`) from the SDK - this is intentional since we want to specify the Resource Manager endpoint for the Azure Environment (e.g. Public, China, US Government etc) that the credentials we're using are connected to.
3. The call to `o.ConfigureClient` configures the authorization token which should be used for this SDK Client - in most cases `ResourceManagerAuthorizer` is the authorizer you want to use.

At this point, this SDK Client should be usable within the Resource via:

```go
client := metadata.Client.{ServicePackage}.{ClientField}
```

For example, in this case:

```go
client := metadata.Client.Resource.GroupsClient
```

### Step 3: Define the Resource ID

Next we're going to generate a Resource ID Struct, Parser and Validator for the specific Azure Resource that we're working with, in this case for a Resource Group.

We have [some automation within the codebase](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/internal/tools/generator-resource-id) which generates all of that using `go:generate` commands - what this means is that we can add a single line to the `resourceids.go` file within the Service Package (in this case `./internal/services/resource/resourceids.go`) to generate these.

An example of this is shown below:

```go
package resource

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupExample -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1
```

In this case, you need to specify the `name` the Resource (in this case `ResourceGroupExample`) and the `id` which is an example of this Resource ID (in this case `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1`).

> The segments of the Resource ID should be camelCased (e.g. `resourceGroups` rather than `resourcegroups`) per the Azure API Specification - see [Azure Resource IDs in the Glossary](reference-glossary.md#azure-resource-ids) for more information.

You can generate the Resource ID Struct, Parser and Validation functions by running `make generate` - which will output the following files:

* `./internal/service/resource/parse/resource_group_example.go` - contains the Resource ID Struct, Formatter and Parser.
* `./internal/service/resource/parse/resource_group_example_test.go` - contains tests for those ^.
* `./internal/service/resource/validate/resource_group_example_id.go` - contains Terraform validation functions for the Resource ID.

These types can then be used in the Resource we're creating below.

### Step 4: Scaffold an empty/new Resource

Since we're creating a Resource for a Resource Group, which is a part of the Resources API - we'll want to create an empty Go file within the Service Package for Resources, which is located at `./internal/services/resource`.

In this case, this'd be a file called `resource_group_example_resource.go`, which we'll start out with the following:

> **Note:** We'd normally name this file `resource_group_resource.go` - but there's an existing Resource for Resource Groups, so we're appending `example` to the name throughout this guide.

```go
package resource

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.Resource = ResourceGroupExampleResource{}

type ResourceGroupExampleResource struct {}
```

> **Note:** Your editor may show a suggestion to implement the methods defined in `sdk.Resource` for the `ResourceGroupExampleResource` struct - we'd recommend holding off the first time around to explain each of the methods.

In this case the interface `sdk.Resource` defines all of the methods required for a Resource which the newly created struct for the Resource Group Resource need to implement, which are:

```go
type Resource interface {
    Arguments() map[string]*schema.Schema
    Attributes() map[string]*schema.Schema
    ModelObject() interface{}
    ResourceType() string
	Create() ResourceFunc
	Read() ResourceFunc
	Delete() ResourceFunc
	IDValidationFunc() pluginsdk.SchemaValidateFunc
}
```

To go through these in turn:

* `Arguments` returns a list of schema fields which are user-specifiable - either Required or Optional.
* `Attributes` returns a list of schema fields which are Computed (read-only).
* `ModelObject` returns a reference to a Go struct which is used as the Model for this Resource (this can also return `nil` if there's no model).
* `ResourceType` returns the name of this resource within the Provider (for example `azurerm_resource_group_example`).
* `Create` returns a function defining both the Timeout and the Create function (which creates this Resource Group using the Azure API) for this Resource.
* `Read` returns a function defining both the Timeout and the Read function (which retrieves information from the Azure API) for this Resource.
* `Delete` returns a function defining both the Timeout and the Delete function (which deletes this Resource Group using the Azure API) for this Resource.
* `IDValidationFunc` returns a function which validates the Resource ID provided during `terraform import` to ensure it matches what we expect for this Resource.

```go
func (ResourceGroupExampleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		
		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (ResourceGroupExampleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ResourceGroupExampleResource) ModelObject() interface{} {
	return nil
}

func (ResourceGroupExampleResource) ResourceType() string {
	return "azurerm_resource_group_example"
}
```

> In this case we're using the resource type `azurerm_resource_group_example` as [an existing Resource for `azurerm_resource_group` exists](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group) and the names need to be unique.

These functions define a Resource called `azurerm_resource_group_example`, which has two Required arguments (`name` and `location`) and one Optional argument (`tags`). We'll come back to `ModelObject` later.

---

Let's start by implementing the Create function:

```go
func (r ResourceGroupExampleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        // the Timeout is how long Terraform should wait for this function to run before returning an error
        // whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
        // limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			// retrieve the Name for this Resource Group from the Terraform Config
			// and then create a Resource ID for this Resource Group
			// using the Subscription ID & name 
			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)
			id := parse.NewResourceGroupID(subscriptionId, name)
			
			// then we want to check for the presence of an existing resource with this name
			// this is because the Azure API uses the `name` as a unique idenfitier and Upserts
			// so we don't want to unintentionally adopt this resource by using the same name
			existing, err := client.Get(ctx, id.ResourceGroup)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}
			
			// create the Resource Group
			param := resources.Group{
				Location: utils.String(location.Normalize(metadata.ResourceData.Get("location").(string))),
				Tags:     tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}
			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// set the Resource ID, meaning that we track this resource
			metadata.SetID(id)
			return nil
		},
	}
}
```


Let's implement the Update function:

```go
func (r ResourceGroupExampleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        // the Timeout is how long Terraform should wait for this function to run before returning an error
        // whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
        // limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			// parse the existing Resource ID from the State
			id, err := parse.ResourceGroupID(metadata.ResourceData.Get("id").(string))
			if err != nil {
				return err
			}
			
			// update the Resource Group
			// NOTE: for a more complex resource we'd recommend retrieving the existing Resource from the
			// API and then conditionally updating it when fields in the config have been updated, which
			// can be determined by using `d.HasChanges` - for example:
			//
			//   existing, err := client.Get(ctx, id.ResourceGroup)
			//   if err != nil {
			//     return fmt.Errorf("retrieving existing %s: %+v", id, err)
			//   }
			//   if d.HasChanges("tags") {
			//     existing.Tags = tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{}))
			//   }
			//
			// doing so allows users to take advantage of Terraform's `ignore_changes` functionality.
			//
			// However since a Resource Group only has one field which is updatable (tags) so in this case we'll only
			// enter the update function if `tags` has been updated.
			param := resources.Group{
				Location: utils.String(location.Normalize(metadata.ResourceData.Get("location").(string))),
				Tags:     tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}
			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// set the Resource ID, meaning that we track this resource
			metadata.SetID(id)
			return nil
		},
	}
}
```


---

Next up, let's implement the Read function - which retrieves the information about the Resource Group from Azure:

```go
func (ResourceGroupExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 5 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 5 * time.Minute,
		
		// the Func returns a function which looks up the state of the Resource Group and sets it into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            client := metadata.Client.Resource.GroupsClient

			// parse the Resource Group ID from the `id` field
            id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			
			// then retrieve the Resource Group by its Name
            resp, err := client.Get(ctx, id.ResourceGroup)
            if err != nil {
				// if the Resource Group doesn't exist (e.g. we get a 404 Not Found)
				// since this is a Resource (e.g. we created it/it was imported into the state)
				// it previously existed - so we must mark this as "gone" for Terraform
                if utils.ResponseWasNotFound(resp.Response) {
                    return metadata.MarkAsGone(id)
                }
				
                // otherwise it's a genuine error (auth/api error etc) so raise it
				// there should be enough context for the user to interpret the error
				// or raise a bug report if there's something we should handle
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }
			
			// at this point we can set information about this Resource Group into the State
			// identifier fields such as the name, resource group name to name a few need to be sourced
			// from the Resource ID instead of the API response
			metadata.ResourceData.Set("name", id.ResourceGroup)
			
			// the SDK will return a Model as well as a nested Properties object for the resource
			// for readability and consistency we assign the Model to a variable and nil check as shown below.
			// since the SDK accounts for responses where the Model is nil we do not need to throw an error if
			// the Model is nil since this will be caught earlier on. We still nil check to prevent the provider from
			// crashing.
			if model := resp.Model; model != nil {
                            // the Location and Tags fields are a little different - and we have a couple of normalization
                            // functions for these.
                            
                            // whilst this may seem like a weird thing to call out in an example, because these two fields
                            // are present on the majority of resources, we hope it explains why they're a little different
                            
                            // in this case the Location can be returned in various different forms, for example
                            // "West Europe", "WestEurope" or "westeurope" - as such we normalize these into a
                            // lower-cased singular word with no spaces (e.g. "westeurope") so this is consistent
                            // for users
                            metadata.ResourceData.Set("location", location.NormalizeNilable(model.Location))
							
                            if props := model.Properties; props != nil {
                                // if there are properties to set into state do that here
                            }
                            
                            // (as above) Tags are a little different, so we have a dedicated helper function
                            // to flatten these consistently across the Provider
							if err := tags.FlattenAndSet(metadata.ResourceData, model.Tags); err != nil {
                                return fmt.Errorf("setting `tags`: %+v", err)
                            }
                        }       
			return nil
		},
	}
}
```

---

Next we can add the Delete function:

```go
func (ResourceGroupExampleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, it can take a while to delete the nested items
		// particularly if we're rate-limited - but users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which deletes the Resource Group
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// NOTE: this is an optional parameter in the SDK we're not concerned with, so we pass an empty string instead
			forceDeletionTypes := ""
			
			// trigger the deletion of the Resource Group
			future, err := client.Delete(ctx, id.ResourceGroup, forceDeletionTypes)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			
			// keep polling until the Resource Group has been deleted
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}
```

---

Finally we can add the `IDValidationFunc` function:

```go
func (ResourceGroupExampleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
    return validate.ResourceGroupID
}
```

---

At this point the finished Resource should look like (including imports):

```go
package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.Resource = ResourceGroupExampleResource{}

type ResourceGroupExampleResource struct{}

func (ResourceGroupExampleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (ResourceGroupExampleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ResourceGroupExampleResource) ModelObject() interface{} {
	return nil
}

func (ResourceGroupExampleResource) ResourceType() string {
	return "azurerm_resource_group_example"
}

func (r ResourceGroupExampleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			// retrieve the Name for this Resource Group from the Terraform Config
			// and then create a Resource ID for this Resource Group
			// using the Subscription ID & name
			subscriptionId := metadata.Client.Account.SubscriptionId
			name := metadata.ResourceData.Get("name").(string)
			id := parse.NewResourceGroupID(subscriptionId, name)

			// then we want to check for the presence of an existing resource with this name
			// this is because the Azure API uses the `name` as a unique idenfitier and Upserts
			// so we don't want to unintentionally adopt this resource by using the same name
			existing, err := client.Get(ctx, id.ResourceGroup)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// create the Resource Group
			param := resources.Group{
				Location: utils.String(location.Normalize(metadata.ResourceData.Get("location").(string))),
				Tags:     tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}
			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// set the Resource ID, meaning that we track this resource
			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceGroupExampleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        // the Timeout is how long Terraform should wait for this function to run before returning an error
        // whilst 30 minutes may initially seem excessive, we set this as a default to account for rate
        // limiting - but having this here means that users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which retrieves the current state of the Resource Group into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			// parse the existing Resource ID from the State
			id, err := parse.ResourceGroupID(metadata.ResourceData.Get("id").(string))
			if err != nil {
				return err
			}
			
			// update the Resource Group
			// NOTE: for a more complex resource we'd recommend retrieving the existing Resource from the
			// API and then conditionally updating it when fields in the config have been updated, which
			// can be determined by using `d.HasChanges` - for example:
			//
			//   existing, err := client.Get(ctx, id.ResourceGroup)
			//   if err != nil {
			//     return fmt.Errorf("retrieving existing %s: %+v", id, err)
			//   }
			//   if d.HasChanges("tags") {
			//     existing.Tags = tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{}))
			//   }
			//
			// doing so allows users to take advantage of Terraform's `ignore_changes` functionality.
			//
			// However since a Resource Group only has one field which is updatable (tags) so in this case we'll only
			// enter the update function if `tags` has been updated.
			param := resources.Group{
				Location: utils.String(location.Normalize(metadata.ResourceData.Get("location").(string))),
				Tags:     tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}
			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// set the Resource ID, meaning that we track this resource
			metadata.SetID(id)
			return nil
		},
	}
}

func (ResourceGroupExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 5 minutes may initially seem excessive, we set this as a default to account for rate
		// limiting - but having this here means that users can override this in their config as necessary
		Timeout: 5 * time.Minute,

		// the Func returns a function which looks up the state of the Resource Group and sets it into the state
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			// parse the Resource Group ID from the `id` field
			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// then retrieve the Resource Group by its Name
			resp, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				// if the Resource Group doesn't exist (e.g. we get a 404 Not Found)
				// since this is a Resource (e.g. we created it/it was imported into the state)
				// it previously existed - so we must mark this as "gone" for Terraform
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				// otherwise it's a genuine error (auth/api error etc) so raise it
				// there should be enough context for the user to interpret the error
				// or raise a bug report if there's something we should handle
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			// at this point we can set information about this Resource Group into the State
			metadata.ResourceData.Set("name", id.ResourceGroup)

			// the Location and Tags fields are a little different - and we have a couple of normalization
			// functions for these.
			//
			// whilst this may seem like a weird thing to call out in an example, because these two fields
			// are present on the majority of resources, we hope it explains why they're a little different
			//
			// in this case the Location can be returned in various different forms, for example
			// "West Europe", "WestEurope" or "westeurope" - as such we normalize these into a
			// lower-cased singular word with no spaces (e.g. "westeurope") so this is consistent
			// for users
			metadata.ResourceData.Set("location", location.NormalizeNilable(resp.Location))

			// (as above) Tags are a little different, so we have a dedicated helper function
			// to flatten these consistently across the Provider
			return tags.FlattenAndSet(metadata.ResourceData, resp.Tags)
		},
	}
}

func (ResourceGroupExampleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// the Timeout is how long Terraform should wait for this function to run before returning an error
		// whilst 30 minutes may initially seem excessive, it can take a while to delete the nested items
		// particularly if we're rate-limited - but users can override this in their config as necessary
		Timeout: 30 * time.Minute,

		// the Func returns a function which deletes the Resource Group
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// NOTE: this is an optional parameter in the SDK we're not concerned with, so we pass an empty string instead
			forceDeletionTypes := ""

			// trigger the deletion of the Resource Group
			future, err := client.Delete(ctx, id.ResourceGroup, forceDeletionTypes)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			// keep polling until the Resource Group has been deleted
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (ResourceGroupExampleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupID
}
```

At this point in time this Resource is now code-complete - there's an optional extension to make this cleaner by using a Typed Model, however this isn't necessary.

### Step 5: Register the new Resource

Resources are registered within the `registration.go` within each Service Package - and should look something like this:

```go
package resource

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

// ...

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}
```

---

> **Note:** It's possible that the Service Registration (above) doesn't currently support Typed Resources, in which case you may need to add the following:

```go
var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct {
}

func (Registration) Name() string {
	return "Some Service"
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

func (Registration) WebsiteCategories() []string {
	return []string{
		"Some Service",
	}
}
```

> In this case you'll also need to add a line to register this Service Registration [in the list of Typed Service Registrations](https://github.com/hashicorp/terraform-provider-azurerm/blob/bd7c755b789fa131778ef93824cf3bae5caccf56/internal/provider/services.go#L109).

---

To register the Resource we need to add an instance of the struct used for the Resource to the list of Resources, for example:

```go
// Resources returns a list of Resources supported by this Service
func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ResourceGroupExampleResource{},	
	}
}
```

At this point the Resource is registered, as when the Azure Provider builds up a list of supported Resources during initialization, it parses each of the Service Registrations to put together a definitive list of the Resources that we support.

This means that if you [Build the Provider](building-the-provider.md), at this point you should be able to apply the following Resource:

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group_example" "test" {
  name     = "example-resources"
  location = "West Europe"
}

output "id" {
  value = azurerm_resource_group_example.test.id
}
```

### Step 6: Add Acceptance Test(s) for this Resource

We're going to test the Resource that we've just built by dynamically provisioning a Resource Group using the new `azurerm_resource_group_example` Resource.

In Go tests are expected to be in a file name in the format `{original_file_name}_test.go` - in our case that'd be `resource_group_example_resource_test.go`, into which we'll want to add:

```go
package resource_test

import (
    "context"
    "fmt"
    "testing"
    
    "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupExampleTestResource struct{}

func TestAccResourceGroupExample_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_example", "test")
	testResource := ResourceGroupExampleTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceGroupExample_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_example", "test")
	testResource := ResourceGroupExampleTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.RequiresImportErrorStep(testResource.requiresImportConfig),
	})
}

func TestAccResourceGroupExample_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_example", "test")
	testResource := ResourceGroupExampleTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.completeConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.completeConfig, testResource),
		data.ImportStep(),
	})
}

func (ResourceGroupExampleTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ResourceGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.GroupsClient.Get(ctx, id.ResourceGroup)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (ResourceGroupExampleTestResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group_example" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupExampleTestResource) requiresImportConfig(data acceptance.TestData) string {
	template := r.basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_example" "import" {
  name     = azurerm_resource_group_example.test.name
  location = azurerm_resource_group_example.test.location
}
`, template)
}

func (ResourceGroupExampleTestResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group_example" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    Hello = "World"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
```

There's a more detailed breakdown of how this works [in the Acceptance Testing reference](reference-acceptance-testing.md) - but to summarize what's going on here:

1. Test Terraform Configurations are defined as methods on the struct `ResourceGroupExampleResource` so that they're easily accessible (this helps to avoid them being unintentionally used in other resources).
2. The `acceptance.TestData` object contains a number of helpers, including both random integers, strings and the Azure Locations where resources should be provisioned - which are used to ensure when tests are run in parallel that we provision unique resources for testing purposes.
3. The `ApplyStep`'s apply the Terraform Configuration specified and then assert there's no changes after (e.g. `terraform apply` and then checking that `terraform plan` shows no changes).
4. The `ImportStep` takes the Resource ID for the Resource and runs `terraform import azurerm_resource_group_example.test {resourceId}`, checking that the fields defined in the state match the fields returned from the Read function.
6. We append `_test` to the Go package name (e.g. `resource_test`) since we need to be able to access both the `resource` package and the `acceptance` package (which is a circular reference, otherwise).

At this point we should be able to run this test.

### Step 7: Run the Acceptance Test(s)

Detailed [instructions on Running the Tests can be found in this guide](running-the-tests.md) - when a Service Principal is configured you can run the test above using:

```sh
make acctests SERVICE='resource' TESTARGS='-run=TestAccResourceGroupExample_' TESTTIMEOUT='60m'
```

> **Note:** We're using the test prefix `TestAccResourceGroupExample_` and not the name of an individual test, but you can do that too by specifying `"(TestName1|TestName2)"` etc

Which should output:

```sh
==> Checking that code complies with gofmt requirements...
==> Checking that Custom Timeouts are used...
==> Checking that acceptance test packages are used...
TF_ACC=1 go test -v ./internal/services/resource -run=TestAccResourceGroupExample_ -timeout 60m -ldflags="-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc"
=== RUN   TestAccResourceGroupExample_basic
=== PAUSE TestAccResourceGroupExample_basic
=== CONT  TestAccResourceGroupExample_basic
--- PASS: TestAccResourceGroupExample_basic (88.15s)
=== RUN   TestAccResourceGroupExample_complete
=== PAUSE TestAccResourceGroupExample_complete
=== CONT  TestAccResourceGroupExample_complete
--- PASS: TestAccResourceGroupExample_complete (120.23s)
=== RUN   TestAccResourceGroupExample_requiresImport
=== PAUSE TestAccResourceGroupExample_requiresImport
=== CONT  TestAccResourceGroupExample_requiresImport
--- PASS: TestAccResourceGroupExample_requiresImport (116.15s)
PASS
ok  	github.com/hashicorp/terraform-provider-azurerm/internal/services/resource	324.753s
```

### Step 8: Add Documentation for this Resource

At this point in time documentation for each Resource (and Data Source) is written manually, located within the `./website` folder - in this case this will be located at `./website/docs/d/resource_group_example.html.markdown`.

There is a tool within the repository to help scaffold the documentation for a Resource - the documentation for this Resource can be scaffolded via the following command:

```sh
$ make scaffold-website BRAND_NAME="Resource Group Example" RESOURCE_NAME="azurerm_resource_group_example" RESOURCE_TYPE="resource" RESOURCE_ID="/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1"
```

The documentation should look something like below - containing both an example usage and the required, optional and computed fields:

> **Note:** In the example below you'll need to replace each `[]` with a backtick "`" - as otherwise this gets rendered incorrectly, unfortunately.

```markdown
---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group_example"
description: |-
  Manages a Resource Group.
---

# azurerm_resource_group_example

Manages a Resource Group.

## Example Usage

[][][]hcl
resource "azurerm_resource_group_example" "example" {
  name     = "example"
  location = "West Europe"
}
[][][]

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Resource Group should exist. Changing this forces a new Resource Group to be created.

* `name` - (Required) The Name which should be used for this Resource Group. Changing this forces a new Resource Group to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Group.

## Import

Resource Groups can be imported using the `resource id`, e.g.

[][][]shell
terraform import azurerm_resource_group_example.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example
[][][]
```

> **Note:** In the example above you'll need to replace each `[]` with a backtick "`" - as otherwise this gets rendered incorrectly, unfortunately.

### Step 9: Send the Pull Request

See [our recommendations for opening a Pull Request](guide-opening-a-pr.md).
