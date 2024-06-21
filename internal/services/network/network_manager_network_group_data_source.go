// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerNetworkGroupDataSource struct{}

var _ sdk.DataSource = ManagerNetworkGroupDataSource{}

func (r ManagerNetworkGroupDataSource) ResourceType() string {
	return "azurerm_network_manager_network_group"
}

func (r ManagerNetworkGroupDataSource) ModelObject() interface{} {
	return &ManagerNetworkGroupModel{}
}

func (r ManagerNetworkGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networkgroups.ValidateNetworkManagerID,
		},
	}
}

func (r ManagerNetworkGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagerNetworkGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkGroups

			var model ManagerNetworkGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := networkgroups.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := networkgroups.NewNetworkGroupID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", id)
			}

			properties := existing.Model.Properties
			state := ManagerNetworkGroupModel{
				Name:             id.NetworkGroupName,
				NetworkManagerId: networkgroups.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
				Description:      pointer.From(properties.Description),
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
