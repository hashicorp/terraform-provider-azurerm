// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpV6AssociationResource struct{}

var _ sdk.Resource = NatGatewayPublicIpV6AssociationResource{}

type NatGatewayPublicIpV6AssociationModel struct {
	NatGatewayId      string `tfschema:"nat_gateway_id"`
	PublicIpAddressId string `tfschema:"public_ip_address_id"`
}

func (r NatGatewayPublicIpV6AssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&natgateways.NatGatewayId{}),

		"public_ip_address_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.PublicIPAddressId{}),
	}
}

func (r NatGatewayPublicIpV6AssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NatGatewayPublicIpV6AssociationResource) ModelObject() interface{} {
	return &NatGatewayPublicIpV6AssociationModel{}
}

func (r NatGatewayPublicIpV6AssociationResource) ResourceType() string {
	return "azurerm_nat_gateway_public_ip_v6_association"
}

func (r NatGatewayPublicIpV6AssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			var state NatGatewayPublicIpV6AssociationModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			publicIpAddressId, err := commonids.ParsePublicIPAddressID(state.PublicIpAddressId)
			if err != nil {
				return err
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(state.NatGatewayId)
			if err != nil {
				return err
			}

			locks.ByID(natGatewayId.ID())
			defer locks.UnlockByID(natGatewayId.ID())

			natGateway, err := client.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *natGatewayId)
				}
				return fmt.Errorf("retrieving %s: %+v", *natGatewayId, err)
			}

			if natGateway.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *natGatewayId)
			}
			if natGateway.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *natGatewayId)
			}

			id := commonids.NewCompositeResourceID(natGatewayId, publicIpAddressId)

			publicIpAddressesV6 := make([]natgateways.SubResource, 0)

			if v := natGateway.Model.Properties.PublicIPAddressesV6; v != nil {
				for _, existingPublicIPAddress := range *v {
					if existingPublicIPAddress.Id == nil {
						continue
					}

					if strings.EqualFold(*existingPublicIPAddress.Id, publicIpAddressId.ID()) {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}

					publicIpAddressesV6 = append(publicIpAddressesV6, existingPublicIPAddress)
				}
			}

			publicIpAddressesV6 = append(publicIpAddressesV6, natgateways.SubResource{
				Id: pointer.To(state.PublicIpAddressId),
			})
			natGateway.Model.Properties.PublicIPAddressesV6 = pointer.To(publicIpAddressesV6)

			if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NatGatewayPublicIpV6AssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
			if err != nil {
				return err
			}

			natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id.First, err)
			}

			state := NatGatewayPublicIpV6AssociationModel{
				NatGatewayId:      id.First.ID(),
				PublicIpAddressId: id.Second.ID(),
			}

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.PublicIPAddressesV6 == nil {
						return metadata.MarkAsGone(id)
					}

					publicIPAddressFound := false
					for _, pip := range *props.PublicIPAddressesV6 {
						if pip.Id == nil {
							continue
						}

						if strings.EqualFold(*pip.Id, id.Second.ID()) {
							publicIPAddressFound = true
							break
						}
					}

					if !publicIPAddressFound {
						return metadata.MarkAsGone(id)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NatGatewayPublicIpV6AssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
			if err != nil {
				return err
			}

			locks.ByID(id.First.ID())
			defer locks.UnlockByID(id.First.ID())

			natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *id.First)
				}
				return fmt.Errorf("retrieving %s: %+v", *id.First, err)
			}

			if natGateway.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id.First)
			}
			if natGateway.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id.First)
			}

			publicIpAddressesV6 := make([]natgateways.SubResource, 0)

			if v := natGateway.Model.Properties.PublicIPAddressesV6; v != nil {
				for _, publicIPAddress := range *v {
					if publicIPAddress.Id == nil {
						continue
					}

					if !strings.EqualFold(*publicIPAddress.Id, id.Second.ID()) {
						publicIpAddressesV6 = append(publicIpAddressesV6, publicIPAddress)
					}
				}
			}
			natGateway.Model.Properties.PublicIPAddressesV6 = pointer.To(publicIpAddressesV6)

			if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NatGatewayPublicIpV6AssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		if _, err := commonids.ParseCompositeResourceID(input.(string), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{}); err != nil {
			errors = append(errors, err)
		}
		return warnings, errors
	}
}
