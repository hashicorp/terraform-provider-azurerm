// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = NetworkSecurityPerimeterAssociationResource{}

type NetworkSecurityPerimeterAssociationResource struct{}

type NetworkSecurityPerimeterAssociationResourceModel struct {
	Name       string `tfschema:"name"`
	ProfileId  string `tfschema:"profile_id"`
	ResourceId string `tfschema:"resource_id"`
	AccessMode string `tfschema:"access_mode"`
}

func (NetworkSecurityPerimeterAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
			ForceNew:     true,
		},

		"profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networksecurityperimeterprofiles.ValidateProfileID,
		},

		"access_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				networksecurityperimeterassociations.PossibleValuesForAssociationAccessMode(),
				false,
			),
		},
	}
}

func (NetworkSecurityPerimeterAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NetworkSecurityPerimeterAssociationResource) ModelObject() interface{} {
	return &NetworkSecurityPerimeterAssociationResourceModel{}
}

func (NetworkSecurityPerimeterAssociationResource) ResourceType() string {
	return "azurerm_network_security_perimeter_association"
}

func (r NetworkSecurityPerimeterAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{

		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAssociationsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config NetworkSecurityPerimeterAssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			profileId, err := networksecurityperimeterprofiles.ParseProfileID(config.ProfileId)
			if err != nil {
				return err
			}

			nspId := networksecurityperimeters.NewNetworkSecurityPerimeterID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.NetworkSecurityPerimeterName)

			id := networksecurityperimeterassociations.NewResourceAssociationID(subscriptionId, nspId.ResourceGroupName, nspId.NetworkSecurityPerimeterName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := networksecurityperimeterassociations.NspAssociation{
				Properties: &networksecurityperimeterassociations.NspAssociationProperties{
					Profile: &networksecurityperimeterassociations.SubResource{
						Id: pointer.To(profileId.ID()),
					},
					PrivateLinkResource: &networksecurityperimeterassociations.SubResource{
						Id: pointer.To(config.ResourceId),
					},
					AccessMode: pointer.To(networksecurityperimeterassociations.AssociationAccessMode(config.AccessMode)),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetworkSecurityPerimeterAssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAssociationsClient

			id, err := networksecurityperimeterassociations.ParseResourceAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkSecurityPerimeterAssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("profile_id") {
				existing.Model.Properties.Profile = &networksecurityperimeterassociations.SubResource{
					Id: pointer.To(config.ProfileId),
				}
			}
			if metadata.ResourceData.HasChange("resource_id") {

				existing.Model.Properties.PrivateLinkResource = &networksecurityperimeterassociations.SubResource{
					Id: pointer.To(config.ResourceId),
				}

			}
			if metadata.ResourceData.HasChange("access_mode") {
				existing.Model.Properties.AccessMode = pointer.To(networksecurityperimeterassociations.AssociationAccessMode(config.AccessMode))

			}

			param := networksecurityperimeterassociations.NspAssociation{
				Properties: existing.Model.Properties,
			}
			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (NetworkSecurityPerimeterAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAssociationsClient

			id, err := networksecurityperimeterassociations.ParseResourceAssociationID(metadata.ResourceData.Id())
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

			state := NetworkSecurityPerimeterAssociationResourceModel{
				Name:       id.ResourceAssociationName,
				ProfileId:  pointer.From(resp.Model.Properties.Profile.Id),
				ResourceId: pointer.From(resp.Model.Properties.PrivateLinkResource.Id),
				AccessMode: string(pointer.From(resp.Model.Properties.AccessMode)),
			}

			return metadata.Encode(&state)
		},
	}
}

func (NetworkSecurityPerimeterAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NetworkSecurityPerimeterAssociationsClient

			id, err := networksecurityperimeterassociations.ParseResourceAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			time.Sleep(5 * time.Second)

			// stateConf := &pluginsdk.StateChangeConf{
			// 	Pending: []string{"Present"},
			// 	Target:  []string{"Deleted"},
			// 	Refresh: func() (interface{}, string, error) {
			// 		resp, err := client.Get(ctx, *id)
			// 		if err != nil {
			// 			if resp.HttpResponse == nil {
			// 				return nil, "Present", err
			// 			}

			// 			if response.WasNotFound(resp.HttpResponse) || response.WasNotFound(resp.HttpResponse) {
			// 				return nil, "Deleted", nil
			// 			}
			// 			return nil, "Present", err
			// 		}
			// 		return resp, "Present", nil
			// 	},

			// 	Timeout:    15 * time.Minute,
			// 	MinTimeout: 15 * time.Second,
			// }

			// _, err = stateConf.WaitForStateContext(ctx)
			// if err != nil {
			// 	return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
			// }
			return nil
		},
	}
}

func (NetworkSecurityPerimeterAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networksecurityperimeterassociations.ValidateResourceAssociationID
}
