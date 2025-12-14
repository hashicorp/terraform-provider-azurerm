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

var _ sdk.Resource = NetworkSecurityPerimeterProfileResource{}

type NetworkSecurityPerimeterProfileResource struct{}

type NetworkSecurityPerimeterProfileResourceModel struct {
	Name        string `tfschema:"name"`
	PerimeterId string `tfschema:"perimeter_id"`
}

func (NetworkSecurityPerimeterProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"perimeter_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: networksecurityperimeterprofiles.ValidateNetworkSecurityPerimeterID,
			ForceNew:     true,
		},
	}
}

func (NetworkSecurityPerimeterProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkSecurityPerimeterProfileResource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterProfileResourceModel{}
}

func (NetworkSecurityPerimeterProfileResource) ResourceType() string {
	return "azurerm_network_security_perimeter_profile"
}

func (r NetworkSecurityPerimeterProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterProfilesClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config NetworkSecurityPerimeterProfileResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			nspId, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(config.PerimeterId)
			if err != nil {
				return err
			}

			id := networksecurityperimeterprofiles.NewProfileID(subscriptionId, nspId.ResourceGroupName, nspId.NetworkSecurityPerimeterName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := networksecurityperimeterprofiles.NspProfile{}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (NetworkSecurityPerimeterProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterProfilesClient

			id, err := networksecurityperimeterprofiles.ParseProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			perimeterId := networksecurityperimeters.NewNetworkSecurityPerimeterID(id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName)

			state := NetworkSecurityPerimeterProfileResourceModel{
				Name:        id.ProfileName,
				PerimeterId: perimeterId.ID(),
			}

			return metadata.Encode(&state)
		},
	}
}

func (NetworkSecurityPerimeterProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterProfilesClient

			id, err := networksecurityperimeterprofiles.ParseProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (NetworkSecurityPerimeterProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networksecurityperimeterprofiles.ValidateProfileID
}
