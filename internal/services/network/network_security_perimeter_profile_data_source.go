// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = NetworkSecurityPerimeterProfileDataSource{}

type NetworkSecurityPerimeterProfileDataSource struct{}

type NetworkSecurityPerimeterProfileDataSourceModel struct {
	Name        string `tfschema:"name"`
	PerimeterId string `tfschema:"perimeter_id"`
}

func (NetworkSecurityPerimeterProfileDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"perimeter_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networksecurityperimeterprofiles.ValidateNetworkSecurityPerimeterID,
		},
	}
}

func (NetworkSecurityPerimeterProfileDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkSecurityPerimeterProfileDataSource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterProfileDataSourceModel{}
}

func (NetworkSecurityPerimeterProfileDataSource) ResourceType() string {
	return "azurerm_network_security_perimeter_profile"
}

func (NetworkSecurityPerimeterProfileDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterProfilesClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state NetworkSecurityPerimeterProfileDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			nspId, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(state.PerimeterId)
			if err != nil {
				return err
			}

			id := networksecurityperimeterprofiles.NewProfileID(subscriptionId, nspId.ResourceGroupName, nspId.NetworkSecurityPerimeterName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
