// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubnetRouteTableAssociationListResource struct{}

type SubnetRouteTableAssociationListModel struct {
	VirtualNetworkId types.String `tfsdk:"virtual_network_id"`
}

var _ sdk.FrameworkListWrappedResource = new(SubnetRouteTableAssociationListResource)

func (SubnetRouteTableAssociationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = subnetRouteTableAssociation
}

func (SubnetRouteTableAssociationListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceSubnetRouteTableAssociation()
}

func (SubnetRouteTableAssociationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_network_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateVirtualNetworkID},
				},
			},
		},
	}
}

func (SubnetRouteTableAssociationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.Subnets

	var data SubnetRouteTableAssociationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := commonids.ParseVirtualNetworkID(data.VirtualNetworkId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", subnetRouteTableAssociation), err)
		return
	}

	resp, err := client.ListComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", subnetRouteTableAssociation), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			if item.Properties == nil || item.Properties.RouteTable == nil {
				continue
			}

			rd := resourceSubnetRouteTableAssociation().Data(&terraform.InstanceState{})
			id, err := commonids.ParseSubnetIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", subnetRouteTableAssociation), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceSubnetRouteTableAssociationFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", subnetRouteTableAssociation), err)
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
