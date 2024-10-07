// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/staticmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	return staticmembers.ValidateStaticMemberID
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
			ValidateFunc: staticmembers.ValidateNetworkGroupID,
		},

		"target_virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
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

			client := metadata.Client.Network.StaticMembers
			networkGroupId, err := staticmembers.ParseNetworkGroupID(model.NetworkGroupId)
			if err != nil {
				return err
			}

			id := staticmembers.NewStaticMemberID(networkGroupId.SubscriptionId, networkGroupId.ResourceGroupName, networkGroupId.NetworkManagerName, networkGroupId.NetworkGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			staticMember := staticmembers.StaticMember{
				Properties: &staticmembers.StaticMemberProperties{
					ResourceId: &model.TargetVNetId,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, staticMember); err != nil {
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
			client := metadata.Client.Network.StaticMembers

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties
			state := ManagerStaticMemberModel{
				Name:           id.StaticMemberName,
				NetworkGroupId: staticmembers.NewNetworkGroupID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.NetworkGroupName).ID(),
			}

			if properties.Region != nil {
				state.Region = *properties.Region
			}

			if properties.ResourceId != nil {
				state.TargetVNetId = *properties.ResourceId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerStaticMemberResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.StaticMembers

			id, err := staticmembers.ParseStaticMemberID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
