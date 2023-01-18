package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
)

type ManagerStaticMemberModel struct {
	Name           string `tfschema:"name"`
	NetworkGroupId string `tfschema:"network_group_id"`
	TargetVNetId   string `tfschema:"target_virtual_network_id"`
	Region         string `tfschema:"region"`
}

type ManagerStaticMemberResource struct{}

var _ sdk.Resource = ManagerStaticMemberResource{}

func (r ManagerStaticMemberResource) ResourceType() string {
	return "azurerm_network_manager_static_member"
}

func (r ManagerStaticMemberResource) ModelObject() interface{} {
	return &ManagerStaticMemberModel{}
}

func (r ManagerStaticMemberResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerStaticMemberID
}

func (r ManagerStaticMemberResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NetworkManagerNetworkGroupID,
		},

		"target_virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualNetworkID,
		},
	}
}

func (r ManagerStaticMemberResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"region": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerStaticMemberResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerStaticMemberModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerStaticMembersClient
			networkGroupId, err := parse.NetworkManagerNetworkGroupID(model.NetworkGroupId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerStaticMemberID(networkGroupId.SubscriptionId, networkGroupId.ResourceGroup, networkGroupId.NetworkManagerName, networkGroupId.NetworkGroupName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			staticMember := &network.StaticMember{
				StaticMemberProperties: &network.StaticMemberProperties{
					ResourceID: &model.TargetVNetId,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *staticMember, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerStaticMemberResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerStaticMembersClient

			id, err := parse.NetworkManagerStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.StaticMemberProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerStaticMemberModel{
				Name:           id.StaticMemberName,
				NetworkGroupId: parse.NewNetworkManagerNetworkGroupID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName).ID(),
			}

			if properties.Region != nil {
				state.Region = *properties.Region
			}

			if properties.ResourceID != nil {
				state.TargetVNetId = *properties.ResourceID
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerStaticMemberResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerStaticMembersClient

			id, err := parse.NetworkManagerStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
