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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIPv6AssociationResource struct{}

var (
	_ sdk.Resource                  = NatGatewayPublicIPv6AssociationResource{}
	_ sdk.ResourceWithCustomizeDiff = NatGatewayPublicIPv6AssociationResource{}
)

type NatGatewayPublicIPv6AssociationModel struct {
	NatGatewayId      string `tfschema:"nat_gateway_id"`
	PublicIpAddressId string `tfschema:"public_ip_address_id"`
}

func (r NatGatewayPublicIPv6AssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&natgateways.NatGatewayId{}),

		"public_ip_address_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.PublicIPAddressId{}),
	}
}

func (r NatGatewayPublicIPv6AssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NatGatewayPublicIPv6AssociationResource) ModelObject() interface{} {
	return &NatGatewayPublicIPv6AssociationModel{}
}

func (r NatGatewayPublicIPv6AssociationResource) ResourceType() string {
	return "azurerm_nat_gateway_public_ipv6_association"
}

func (r NatGatewayPublicIPv6AssociationResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rawNatGatewayId := metadata.ResourceDiff.GetRawConfig().AsValueMap()["nat_gateway_id"]
			if rawNatGatewayId.IsNull() || !rawNatGatewayId.IsKnown() {
				return nil
			}

			rawPublicIPAddressId := metadata.ResourceDiff.GetRawConfig().AsValueMap()["public_ip_address_id"]
			if rawPublicIPAddressId.IsNull() || !rawPublicIPAddressId.IsKnown() {
				return nil
			}

			var model NatGatewayPublicIPv6AssociationModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(model.NatGatewayId)
			if err != nil {
				return err
			}

			natGatewayResp, err := metadata.Client.Network.NatGateways.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *natGatewayId, err)
			}

			if natGatewayResp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *natGatewayId)
			}

			natGatewaySku := pointer.From(pointer.From(natGatewayResp.Model.Sku).Name)
			if natGatewaySku == natgateways.NatGatewaySkuNameStandard {
				return fmt.Errorf("`nat_gateway_id` with SKU `Standard` does not support IPv6")
			}

			publicIPAddressId, err := commonids.ParsePublicIPAddressID(model.PublicIpAddressId)
			if err != nil {
				return err
			}

			resp, err := metadata.Client.Network.PublicIPAddresses.Get(ctx, *publicIPAddressId, publicipaddresses.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *publicIPAddressId, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", *publicIPAddressId)
			}

			publicIPAddressSku := pointer.From(pointer.From(resp.Model.Sku).Name)
			if natGatewaySku == natgateways.NatGatewaySkuNameStandardVTwo && publicIPAddressSku != publicipaddresses.PublicIPAddressSkuNameStandardVTwo {
				return fmt.Errorf("`public_ip_address_id` must use SKU `StandardV2` when `nat_gateway_id` uses SKU `StandardV2`, got `%s`", publicIPAddressSku)
			}

			if version := pointer.From(resp.Model.Properties.PublicIPAddressVersion); version != publicipaddresses.IPVersionIPvSix {
				return fmt.Errorf("`public_ip_address_id` must use `IPv6`, got `%s`", version)
			}

			return nil
		},
	}
}

func (r NatGatewayPublicIPv6AssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			var state NatGatewayPublicIPv6AssociationModel
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

			existingIPs := pointer.From(natGateway.Model.Properties.PublicIPAddressesV6)
			for _, existingPublicIPAddress := range existingIPs {
				if strings.EqualFold(pointer.From(existingPublicIPAddress.Id), publicIpAddressId.ID()) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			existingIPs = append(existingIPs, natgateways.SubResource{
				Id: pointer.To(state.PublicIpAddressId),
			})
			natGateway.Model.Properties.PublicIPAddressesV6 = pointer.To(existingIPs)

			if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NatGatewayPublicIPv6AssociationResource) Read() sdk.ResourceFunc {
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

			state := NatGatewayPublicIPv6AssociationModel{
				NatGatewayId:      id.First.ID(),
				PublicIpAddressId: id.Second.ID(),
			}

			publicIPAddressFound := false
			if natGateway.Model == nil || natGateway.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id.First)
			}

			for _, pip := range pointer.From(natGateway.Model.Properties.PublicIPAddressesV6) {
				if strings.EqualFold(pointer.From(pip.Id), id.Second.ID()) {
					publicIPAddressFound = true
					break
				}
			}

			if !publicIPAddressFound {
				return metadata.MarkAsGone(id)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NatGatewayPublicIPv6AssociationResource) Delete() sdk.ResourceFunc {
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
			needToDelete := false
			for _, publicIPAddress := range pointer.From(natGateway.Model.Properties.PublicIPAddressesV6) {
				if !strings.EqualFold(pointer.From(publicIPAddress.Id), id.Second.ID()) {
					publicIpAddressesV6 = append(publicIpAddressesV6, publicIPAddress)
				} else {
					needToDelete = true
				}
			}

			if !needToDelete {
				return nil
			}
			natGateway.Model.Properties.PublicIPAddressesV6 = pointer.To(publicIpAddressesV6)

			if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NatGatewayPublicIPv6AssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		if _, err := commonids.ParseCompositeResourceID(input.(string), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{}); err != nil {
			errors = append(errors, err)
		}
		return warnings, errors
	}
}
