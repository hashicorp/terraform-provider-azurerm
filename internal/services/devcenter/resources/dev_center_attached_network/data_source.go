// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_attached_network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type DevCenterAttachedNetworkDataSourceModel struct {
	Name                string `tfschema:"name"`
	DevCenterId         string `tfschema:"dev_center_id"`
	NetworkConnectionId string `tfschema:"network_connection_id"`
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Attached Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devcenters.DevCenterId{}, "Dev Center"),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_connection_id": {
			Description: "The ID of the attached Dev Center Network Connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (DataSource) ModelObject() any {
	return &DevCenterAttachedNetworkDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_dev_center_attached_network"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Dev Center Attached Network."
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.AttachedNetworkConnections
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var state DevCenterAttachedNetworkDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devcenters.ParseDevCenterID(state.DevCenterId)
			if err != nil {
				return err
			}

			id := attachednetworkconnections.NewDevCenterAttachedNetworkID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.AttachedNetworksGetByDevCenter(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.AttachedNetworkName
			state.DevCenterId = attachednetworkconnections.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.NetworkConnectionId = props.NetworkConnectionId
				}
			}

			return metadata.Encode(&state)
		},
	}
}
