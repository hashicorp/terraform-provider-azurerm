package devcenter

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devcenters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

type DevCenterResource struct{}

type DevCenterModel struct {
	Name          string            `tfschema:"name"`
	ResourceGroup string            `tfschema:"resource_group_name"`
	Location      string            `tfschema:"location"`
	Identity      []interface{}     `tfschema:"identity"`
	Tags          map[string]string `tfschema:"tags"`
}

type DevCenterropertiesModel struct {
	DevCenterUri string `tfschema:"dev_center_uri"`
}

var _ sdk.ResourceWithUpdate = DevCenterResource{}

func (r DevCenterResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r DevCenterResource) ModelObject() interface{} {
	return &DevCenterModel{}
}

func (r DevCenterResource) ResourceType() string {
	return "azurerm_dev_center"
}

func (r DevCenterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.DevCenterClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := devcenters.NewDevCenterID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			devCenterIdentity, err := identity.ExpandSystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			devCenter := devcenters.DevCenter{
				Id:       utils.String(id.ID()),
				Identity: devCenterIdentity,
				Location: model.Location,
				Name:     utils.String(model.Name),
				Tags:     pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, devCenter); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.DevCenterClient
			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				state := DevCenterModel{
					Name:          id.DevCenterName,
					ResourceGroup: id.ResourceGroupName,
					Location:      model.Location,
				}
				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
				devCenterIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(devCenterIdentity)

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r DevCenterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.DevCenterClient
			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return devcenters.ValidateDevCenterID
}

func (r DevCenterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.DevCenterClient
			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state DevCenterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterIdentity, err := identity.ExpandSystemAndUserAssignedMap(state.Identity)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChangesExcept("name", "resource_group_name", "location") {
				devCenter := devcenters.DevCenter{
					Identity: devCenterIdentity,
					Tags:     pointer.To(state.Tags),
				}

				if err := client.CreateOrUpdateThenPoll(ctx, *id, devCenter); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}
