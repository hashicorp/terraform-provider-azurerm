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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = NATGatewayPublicIpPrefixV6AssociationResource{}

type NATGatewayPublicIpPrefixV6AssociationResource struct{}

type NATGatewayPublicIpPrefixV6AssociationModel struct {
	NATGatewayId     string `tfschema:"nat_gateway_id"`
	PublicIPPrefixId string `tfschema:"public_ip_prefix_id"`
}

func (NATGatewayPublicIpPrefixV6AssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&natgateways.NatGatewayId{}),

		"public_ip_prefix_id": commonschema.ResourceIDReferenceRequiredForceNew(&publicipprefixes.PublicIPPrefixId{}),
	}
}

func (NATGatewayPublicIpPrefixV6AssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NATGatewayPublicIpPrefixV6AssociationResource) ModelObject() interface{} {
	return &NATGatewayPublicIpPrefixV6AssociationModel{}
}

func (NATGatewayPublicIpPrefixV6AssociationResource) ResourceType() string {
	return "azurerm_nat_gateway_public_ip_prefix_v6_association"
}

func (r NATGatewayPublicIpPrefixV6AssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			var state NATGatewayPublicIpPrefixV6AssociationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			publicIpPrefixId, err := publicipprefixes.ParsePublicIPPrefixID(state.PublicIPPrefixId)
			if err != nil {
				return err
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(state.NATGatewayId)
			if err != nil {
				return err
			}

			locks.ByName(natGatewayId.NatGatewayName, natGatewayResourceName)
			defer locks.UnlockByName(natGatewayId.NatGatewayName, natGatewayResourceName)

			natGateway, err := client.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *natGatewayId)
				}
				return fmt.Errorf("retrieving %s: %+v", *natGatewayId, err)
			}

			id := commonids.NewCompositeResourceID(natGatewayId, publicIpPrefixId)

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					publicIpPrefixesV6 := make([]natgateways.SubResource, 0)

					if v := props.PublicIPPrefixesV6; v != nil {
						for _, existingPublicIPPrefix := range *v {
							if existingPublicIPPrefix.Id == nil {
								continue
							}

							if strings.EqualFold(*existingPublicIPPrefix.Id, publicIpPrefixId.ID()) {
								return metadata.ResourceRequiresImport(r.ResourceType(), id)
							}

							publicIpPrefixesV6 = append(publicIpPrefixesV6, existingPublicIPPrefix)
						}
					}

					publicIpPrefixesV6 = append(publicIpPrefixesV6, natgateways.SubResource{
						Id: pointer.To(publicIpPrefixId.ID()),
					})
					props.PublicIPPrefixesV6 = pointer.To(publicIpPrefixesV6)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *natGatewayId, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (NATGatewayPublicIpPrefixV6AssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
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

			state := NATGatewayPublicIpPrefixV6AssociationModel{
				NATGatewayId:     id.First.ID(),
				PublicIPPrefixId: id.Second.ID(),
			}

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.PublicIPPrefixesV6 == nil {
						return metadata.MarkAsGone(id)
					}

					found := false
					for _, pipPrefix := range *props.PublicIPPrefixesV6 {
						if pipPrefix.Id == nil {
							continue
						}

						if strings.EqualFold(*pipPrefix.Id, id.Second.ID()) {
							found = true
							break
						}
					}

					if !found {
						return metadata.MarkAsGone(id)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (NATGatewayPublicIpPrefixV6AssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
			if err != nil {
				return err
			}

			locks.ByName(id.First.NatGatewayName, natGatewayResourceName)
			defer locks.UnlockByName(id.First.NatGatewayName, natGatewayResourceName)

			natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *id.First)
				}
				return fmt.Errorf("retrieving %s: %+v", *id.First, err)
			}

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					publicIpPrefixesV6 := make([]natgateways.SubResource, 0)

					if v := props.PublicIPPrefixesV6; v != nil {
						for _, publicIPPrefix := range *v {
							if publicIPPrefix.Id == nil {
								continue
							}

							if !strings.EqualFold(*publicIPPrefix.Id, id.Second.ID()) {
								publicIpPrefixesV6 = append(publicIpPrefixesV6, publicIPPrefix)
							}
						}
					}
					props.PublicIPPrefixesV6 = pointer.To(publicIpPrefixesV6)
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
				return fmt.Errorf("removing association between %s and %s: %+v", *id.First, *id.Second, err)
			}

			return nil
		},
	}
}

func (r NATGatewayPublicIpPrefixV6AssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		_, err := commonids.ParseCompositeResourceID(input.(string), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
		if err != nil {
			errors = append(errors, err)
		}
		return warnings, errors
	}
}
