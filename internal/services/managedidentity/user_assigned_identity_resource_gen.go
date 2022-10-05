package managedidentity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentities"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type UserAssignedIdentityResource struct {
}

func (r UserAssignedIdentityResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 128),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r UserAssignedIdentityResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"principal_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"client_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r UserAssignedIdentityResource) ModelObject() interface{} {
	return nil
}

func (r UserAssignedIdentityResource) ResourceType() string {
	return "azurerm_user_assigned_identity"
}

func (r UserAssignedIdentityResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.ManagedIdentities
			subscriptionId := metadata.Client.Account.SubscriptionId

			t := metadata.ResourceData.Get("tags").(map[string]interface{})

			resourceId := commonids.NewUserAssignedIdentityID(subscriptionId, metadata.ResourceData.Get("resource_group_name").(string), metadata.ResourceData.Get("name").(string))
			if metadata.ResourceData.IsNewResource() {
				existing, err := client.UserAssignedIdentitiesGet(ctx, resourceId)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", resourceId, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_user_assigned_identity", resourceId.ID())
				}
			}

			identity := managedidentities.Identity{
				Location: location.Normalize(metadata.ResourceData.Get("location").(string)),
				Tags:     tags.Expand(t),
			}

			if _, err := client.UserAssignedIdentitiesCreateOrUpdate(ctx, resourceId, identity); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", resourceId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r UserAssignedIdentityResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.ManagedIdentities
			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.UserAssignedIdentitiesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.ResourceName)
			metadata.ResourceData.Set("resource_group_name", id.ResourceGroupName)

			if model := resp.Model; model != nil {
				metadata.ResourceData.Set("location", location.Normalize(model.Location))

				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("client_id", props.ClientId)
					metadata.ResourceData.Set("principal_id", props.PrincipalId)
					metadata.ResourceData.Set("tenant_id", props.TenantId)
				}

				if err := tags.FlattenAndSet(metadata.ResourceData, model.Tags); err != nil {
					return err
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r UserAssignedIdentityResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.ManagedIdentities
			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.UserAssignedIdentitiesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r UserAssignedIdentityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateUserAssignedIdentityID
}
