// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/ipampools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerIpamPoolDataSource struct{}

var _ sdk.DataSource = ManagerIpamPoolDataSource{}

func (r ManagerIpamPoolDataSource) ResourceType() string {
	return "azurerm_network_manager_ipam_pool"
}

func (r ManagerIpamPoolDataSource) ModelObject() interface{} {
	return &ManagerIpamPoolDataSourceModel{}
}

type ManagerIpamPoolDataSourceModel struct {
	AddressPrefixes  []string          `tfschema:"address_prefixes"`
	Description      string            `tfschema:"description"`
	DisplayName      string            `tfschema:"display_name"`
	Location         string            `tfschema:"location"`
	Name             string            `tfschema:"name"`
	NetworkManagerId string            `tfschema:"network_manager_id"`
	ParentPoolName   string            `tfschema:"parent_pool_name"`
	Tags             map[string]string `tfschema:"tags"`
}

func (r ManagerIpamPoolDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": commonschema.ResourceIDReferenceRequired(&ipampools.NetworkManagerId{}),
	}
}

func (r ManagerIpamPoolDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"address_prefixes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"parent_pool_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ManagerIpamPoolDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.IPamPools

			var model ManagerIpamPoolDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := ipampools.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := ipampools.NewIPamPoolID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagerIpamPoolDataSourceModel{
				Name:             id.IpamPoolName,
				NetworkManagerId: ipampools.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties
				state.AddressPrefixes = props.AddressPrefixes
				state.Description = pointer.From(props.Description)
				state.DisplayName = pointer.From(props.DisplayName)
				state.ParentPoolName = pointer.From(props.ParentPoolName)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
