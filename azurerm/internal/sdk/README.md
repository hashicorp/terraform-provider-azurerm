## SDK for Strongly-Typed Resources

This package is a prototype for creating strongly-typed Data Sources and Resources - and in future will likely form the foundation for Terraform Data Sources and Resources in this Provider going forward.

## Should I use this package to build resources?

Not at this time - please use Terraform's Plugin SDK instead - reference examples can be found in `./azurerm/internal/services/notificationhub`.

More documentation for this package will ship in the future when this is ready for general use.

---

## What's the long-term intention for this package?

Each Service Package contains the following:

* Client - giving reference to the SDK Client which should be used to interact with Azure
* ID Parsers, Formatters and a Validator - giving a canonical ID for each Resource 
* Validation functions specific to this service package, for example for the Name

This package can be used to tie these together in a more strongly typed fashion, for example:

```
package example

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ResourceGroup struct {
	Name     string            `tfschema:"name"`
	Location string            `tfschema:"location"`
	Tags     map[string]string `tfschema:"tags"`
}

type ResourceGroupResource struct {
}

func (r ResourceGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": location.Schema(),

		"tags": tags.Schema(),
	}
}

func (r ResourceGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceGroupResource) ResourceType() string {
	return "azurerm_example"
}

func (r ResourceGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ResourceGroup
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("creating Resource Group %q..", state.Name)
			client := metadata.Client.Resource.GroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewResourceGroupID(subscriptionId, state.Name)
			existing, err := client.Get(ctx, state.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of an existing Resource Group %q: %+v", state.Name, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := resources.Group{
				Location: utils.String(state.Location),
				Tags:     tags.FromTypedObject(state.Tags),
			}
			if _, err := client.CreateOrUpdate(ctx, state.Name, input); err != nil {
				return fmt.Errorf("creating Resource Group %q: %+v", state.Name, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient
			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving Resource Group %q..", id.ResourceGroup)
			group, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(group.Response) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return metadata.Encode(&ResourceGroup{
				Name:     id.ResourceGroup,
				Location: location.NormalizeNilable(group.Location),
				Tags:     tags.ToTypedObject(group.Tags),
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ResourceGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state..")
			var state ResourceGroup
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Resource.GroupsClient

			input := resources.GroupPatchable{
				Tags: tags.FromTypedObject(state.Tags),
			}

			if _, err := client.Update(ctx, id.ResourceGroup, input); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient
			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			future, err := client.Delete(ctx, id.ResourceGroup, "")
			if err != nil {
				if response.WasNotFound(future.Response()) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			metadata.Logger.Infof("waiting for the deletion of %s..", *id)
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupID
}

func (r ResourceGroupResource) ModelObject() interface{} {
	return ResourceGroup{}
}
```

The end result being the removal of a lot of common bugs by moving to a convention - for example:

* The Context object passed into each method _always_ has a deadline/timeout attached to it
* The Read function is automatically called at the end of a Create and Update function - meaning users don't have to do this 
* Each Resource has to have an ID Formatter and Validation Function
* The Model Object is validated via unit tests to ensure it contains the relevant struct tags (TODO: also confirming these exist in the state and are of the correct type, so no Set errors occur)

Ultimately this allows bugs to be caught by the Compiler (for example if a Read function is unimplemented) - or Unit Tests (for example should the `tfschema` struct tags be missing) - rather than during Provider Initialization, which reduces the feedback loop.
