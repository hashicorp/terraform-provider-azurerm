// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	VirtualNetworkPeeringListResource struct{}
	VirtualNetworkPeeringListModel    struct {
		VirtualNetworkId types.String `tfsdk:"virtual_network_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(VirtualNetworkPeeringListResource)

func (r VirtualNetworkPeeringListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceVirtualNetworkPeering()
}

func (r VirtualNetworkPeeringListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_virtual_network_peering"
}

func (r VirtualNetworkPeeringListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_network_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateVirtualNetworkID,
					},
				},
			},
		},
	}
}

func (r VirtualNetworkPeeringListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.VirtualNetworkPeerings
	var data VirtualNetworkPeeringListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]virtualnetworkpeerings.VirtualNetworkPeering, 0)

	if !data.VirtualNetworkId.IsNull() {
		virtualnetworkId, err := commonids.ParseVirtualNetworkIDInsensitively(data.VirtualNetworkId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Virtual Network ID for `%s`", "azurerm_virtual_network_peering"), err)
			return
		}

		resp, err := client.ListComplete(ctx, *virtualnetworkId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_virtual_network_peering"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, virtualnetworkpeering := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(virtualnetworkpeering.Name)

			rd := resourceVirtualNetworkPeering().Data(&terraform.InstanceState{})

			id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(pointer.From(virtualnetworkpeering.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Virtual Network Peering ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceVirtualNetworkPeeringFlatten(rd, id, &virtualnetworkpeering); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_virtual_network_peering"), err)
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
