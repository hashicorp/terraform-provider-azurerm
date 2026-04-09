// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpPrefixIPv6AssociationResource struct{}

var (
	_ sdk.Resource                  = NatGatewayPublicIpPrefixIPv6AssociationResource{}
	_ sdk.ResourceWithCustomizeDiff = NatGatewayPublicIpPrefixIPv6AssociationResource{}
)

type NatGatewayPublicIpPrefixIPv6AssociationModel struct {
	NATGatewayId     string `tfschema:"nat_gateway_id"`
	PublicIPPrefixId string `tfschema:"public_ip_prefix_id"`
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&natgateways.NatGatewayId{}),

		"public_ip_prefix_id": commonschema.ResourceIDReferenceRequiredForceNew(&publicipprefixes.PublicIPPrefixId{}),
	}
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) ModelObject() interface{} {
	return &NatGatewayPublicIpPrefixIPv6AssociationModel{}
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) ResourceType() string {
	return "azurerm_nat_gateway_public_ip_prefix_ipv6_association"
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rawNatGatewayId := metadata.ResourceDiff.GetRawConfig().AsValueMap()["nat_gateway_id"]
			if rawNatGatewayId.IsNull() || !rawNatGatewayId.IsKnown() {
				return nil
			}

			rawPublicIPPrefixId := metadata.ResourceDiff.GetRawConfig().AsValueMap()["public_ip_prefix_id"]
			if rawPublicIPPrefixId.IsNull() || !rawPublicIPPrefixId.IsKnown() {
				return nil
			}

			var model NatGatewayPublicIpPrefixIPv6AssociationModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(model.NATGatewayId)
			if err != nil {
				return err
			}

			natGatewayResp, err := metadata.Client.Network.NatGateways.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGatewayResp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", *natGatewayId, err)
			}

			if natGatewayResp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *natGatewayId)
			}

			natGatewaySku := pointer.From(pointer.From(natGatewayResp.Model.Sku).Name)
			if natGatewaySku == natgateways.NatGatewaySkuNameStandard {
				return errors.New("`nat_gateway_id` must reference a NAT Gateway with SKU `StandardV2`")
			}

			publicIPPrefixId, err := publicipprefixes.ParsePublicIPPrefixID(model.PublicIPPrefixId)
			if err != nil {
				return err
			}

			resp, err := metadata.Client.Network.PublicIPPrefixes.Get(ctx, *publicIPPrefixId, publicipprefixes.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", *publicIPPrefixId, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", *publicIPPrefixId)
			}

			publicIPPrefixSku := pointer.From(pointer.From(resp.Model.Sku).Name)
			version := pointer.From(resp.Model.Properties.PublicIPAddressVersion)
			if natGatewaySku == natgateways.NatGatewaySkuNameStandardVTwo && (publicIPPrefixSku != publicipprefixes.PublicIPPrefixSkuNameStandardVTwo || version != publicipprefixes.IPVersionIPvSix) {
				return errors.New("`public_ip_prefix_id` must reference an `IPv6` Public IP Prefix with SKU `StandardV2`")
			}

			return nil
		},
	}
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			var state NatGatewayPublicIpPrefixIPv6AssociationModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			publicIPPrefixId, err := publicipprefixes.ParsePublicIPPrefixID(state.PublicIPPrefixId)
			if err != nil {
				return err
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(state.NATGatewayId)
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

			id := commonids.NewCompositeResourceID(natGatewayId, publicIPPrefixId)

			existingPrefixes := pointer.From(natGateway.Model.Properties.PublicIPPrefixesV6)
			for _, existingPublicIPPrefix := range existingPrefixes {
				if strings.EqualFold(pointer.From(existingPublicIPPrefix.Id), publicIPPrefixId.ID()) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			existingPrefixes = append(existingPrefixes, natgateways.SubResource{
				Id: pointer.To(publicIPPrefixId.ID()),
			})
			natGateway.Model.Properties.PublicIPPrefixesV6 = pointer.To(existingPrefixes)

			if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) Read() sdk.ResourceFunc {
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

			state := NatGatewayPublicIpPrefixIPv6AssociationModel{
				NATGatewayId:     id.First.ID(),
				PublicIPPrefixId: id.Second.ID(),
			}

			found := false
			if natGateway.Model == nil || natGateway.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id.First)
			}

			for _, publicIPPrefix := range pointer.From(natGateway.Model.Properties.PublicIPPrefixesV6) {
				if strings.EqualFold(pointer.From(publicIPPrefix.Id), id.Second.ID()) {
					found = true
					break
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			return metadata.Encode(&state)
		},
	}
}

func (NatGatewayPublicIpPrefixIPv6AssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
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

			publicIPPrefixesV6 := make([]natgateways.SubResource, 0)
			needToDelete := false
			for _, publicIPPrefix := range pointer.From(natGateway.Model.Properties.PublicIPPrefixesV6) {
				if !strings.EqualFold(pointer.From(publicIPPrefix.Id), id.Second.ID()) {
					publicIPPrefixesV6 = append(publicIPPrefixesV6, publicIPPrefix)
				} else {
					needToDelete = true
				}
			}

			if !needToDelete {
				return nil
			}

			natGateway.Model.Properties.PublicIPPrefixesV6 = pointer.To(publicIPPrefixesV6)

			if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NatGatewayPublicIpPrefixIPv6AssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		if _, err := commonids.ParseCompositeResourceID(input.(string), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{}); err != nil {
			errors = append(errors, err)
		}
		return warnings, errors
	}
}
