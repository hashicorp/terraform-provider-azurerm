package resource

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupResource struct{}

type ResourceGroupModel struct {
	Name     string            `tfschema:"name"`
	Location string            `tfschema:"location"`
	Tags     map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = ResourceGroupResource{}

func (r ResourceGroupResource) ModelObject() interface{} {
	return &ResourceGroupModel{}
}

func (r ResourceGroupResource) ResourceType() string {
	return "azurerm_resource_group"
}

func (r ResourceGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ResourceGroupID
}

func (r ResourceGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r ResourceGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			var model ResourceGroupModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := parse.NewResourceGroupID(metadata.Client.Account.SubscriptionId, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing resource group: %+v", err)
				}
			}

			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
			}

			parameters := resources.Group{
				Location: utils.String(model.Location),
				Tags:     tags.FromTypedObject(model.Tags),
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			resp, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.ResourceData.SetId(*resp.ID)

			return nil
		},
	}
}

func (r ResourceGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient
			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state := ResourceGroupModel{
				Name:     utils.NormalizeNilableString(resp.Name),
				Location: location.NormalizeNilable(resp.Location),
				Tags:     tags.ToTypedObject(resp.Tags),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if metadata.Client.Features.ResourceGroup.PreventDeletionIfContainsResources {
				resourceClient := metadata.Client.Resource.ResourcesClient
				err = pluginsdk.Retry(10*time.Minute, func() *pluginsdk.RetryError {
					results, err := resourceClient.ListByResourceGroupComplete(ctx, id.ResourceGroup, "", "provisioningState", utils.Int32(500))
					if err != nil {
						return pluginsdk.NonRetryableError(fmt.Errorf("listing resources in %s: %v", *id, err))
					}
					nestedResourceIds := make([]string, 0)
					for results.NotDone() {
						val := results.Value()
						if val.ID != nil {
							nestedResourceIds = append(nestedResourceIds, *val.ID)
						}

						if err := results.NextWithContext(ctx); err != nil {
							return pluginsdk.NonRetryableError(fmt.Errorf("retrieving next page of nested items for %s: %+v", id, err))
						}
					}

					if len(nestedResourceIds) > 0 {
						time.Sleep(30 * time.Second)
						return pluginsdk.RetryableError(resourceGroupContainsItemsError(id.ResourceGroup, nestedResourceIds))
					}
					return nil
				})
				if err != nil {
					return err
				}
			}

			deleteFuture, err := client.Delete(ctx, id.ResourceGroup, "")
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
			if err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ResourceGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.GroupsClient

			var model ResourceGroupModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := parse.ResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			existing.Tags = tags.FromTypedObject(model.Tags)

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func resourceGroupContainsItemsError(name string, nestedResourceIds []string) error {
	formattedResourceUris := make([]string, 0)
	for _, id := range nestedResourceIds {
		formattedResourceUris = append(formattedResourceUris, fmt.Sprintf("* `%s`", id))
	}
	sort.Strings(formattedResourceUris)

	message := fmt.Sprintf(`deleting Resource Group %[1]q: the Resource Group still contains Resources.

Terraform is configured to check for Resources within the Resource Group when deleting the Resource Group - and
raise an error if nested Resources still exist to avoid unintentionally deleting these Resources.

Terraform has detected that the following Resources still exist within the Resource Group:

%[2]s

This feature is intended to avoid the unintentional destruction of nested Resources provisioned through some
other means (for example, an ARM Template Deployment) - as such you must either remove these Resources, or
disable this behaviour using the feature flag 'prevent_deletion_if_contains_resources' within the 'features'
block when configuring the Provider, for example:

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

When that feature flag is set, Terraform will skip checking for any Resources within the Resource Group and
delete this using the Azure API directly (which will clear up any nested resources).

More information on the 'features' block can be found in the documentation:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features
`, name, strings.Join(formattedResourceUris, "\n"))
	return fmt.Errorf(strings.ReplaceAll(message, "'", "`"))
}
