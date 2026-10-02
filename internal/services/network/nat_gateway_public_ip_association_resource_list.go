// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NatGatewayPublicIpAssociationListResource struct{}

type NatGatewayPublicIpAssociationListModel struct {
	NatGatewayId      types.String `tfsdk:"nat_gateway_id"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
}

var _ sdk.FrameworkListWrappedResource = &NatGatewayPublicIpAssociationListResource{}

func (r NatGatewayPublicIpAssociationListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNATGatewayPublicIpAssociation()
}

func (r NatGatewayPublicIpAssociationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_nat_gateway_public_ip_association"
}

func (r NatGatewayPublicIpAssociationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_group_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
					stringvalidator.ConflictsWith(path.MatchRoot("nat_gateway_id")),
				},
			},
			"nat_gateway_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: natgateways.ValidateNatGatewayID,
					},
					stringvalidator.ConflictsWith(path.MatchRoot("resource_group_name")),
				},
			},
		},
	}
}

func (r NatGatewayPublicIpAssociationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.NatGateways

	var data NatGatewayPublicIpAssociationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var gateways []natgateways.NatGateway

	switch {
	case !data.NatGatewayId.IsNull():
		natGatewayId, err := natgateways.ParseNatGatewayID(data.NatGatewayId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, "parsing `nat_gateway_id`", err)
			return
		}

		natGateway, err := client.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving %s", natGatewayId), err)
			return
		}

		if natGateway.Model != nil {
			gateways = append(gateways, *natGateway.Model)
		}

	case !data.ResourceGroupName.IsNull():
		rgId := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, data.ResourceGroupName.ValueString())
		resp, err := client.ListComplete(ctx, rgId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing NAT Gateways in %s", rgId), err)
			return
		}

		gateways = resp.Items

	default:
		subId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
		resp, err := client.ListAllComplete(ctx, subId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing NAT Gateways in %s", subId), err)
			return
		}

		gateways = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, gateway := range gateways {
			if gateway.Id == nil || gateway.Properties == nil {
				continue
			}

			gatewayId, err := natgateways.ParseNatGatewayID(pointer.From(gateway.Id))
			if err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing NAT Gateway ID", err)
				return
			}

			var publicIpAddresses []natgateways.SubResource
			publicIpAddresses = append(publicIpAddresses, pointer.From(gateway.Properties.PublicIPAddresses)...)
			publicIpAddresses = append(publicIpAddresses, pointer.From(gateway.Properties.PublicIPAddressesV6)...)

			for _, publicIPAddress := range publicIpAddresses {
				rawPublicIPId := pointer.From(publicIPAddress.Id)
				publicIpId, err := commonids.ParsePublicIPAddressID(rawPublicIPId)
				if err != nil {
					result := request.NewListResult(ctx)
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Public IP Address ID", err)
					return
				}

				id := commonids.NewCompositeResourceID(gatewayId, publicIpId)

				result := request.NewListResult(ctx)
				result.DisplayName = fmt.Sprintf("%s - %s", gatewayId.NatGatewayName, publicIpId.PublicIPAddressesName)

				rd := resourceNATGatewayPublicIpAssociation().Data(&terraform.InstanceState{})
				rd.SetId(id.ID())
				rd.Set("nat_gateway_id", id.First.ID())
				rd.Set("public_ip_address_id", id.Second.ID())

				if err := pluginsdk.SetCompositeResourceIdentityData(rd, id, "nat_gateway_id", "public_ip_address_id"); err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting resource identity for `azurerm_nat_gateway_public_ip_association`", err)
					return
				}

				sdk.EncodeListResult(ctx, rd, &result)
				if result.Diagnostics.HasError() {
					push(result)
					return
				}

				if !push(result) {
					return
				}
			}
		}
	}
}
