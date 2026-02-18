// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDnsZoneVirtualNetworkLinkListResource struct{}
type PrivateDnsZoneVirtualNetworkLinkListModel struct {
	PrivateDnsZoneId types.String `tfsdk:"private_dns_zone_id"`
}

var _ sdk.FrameworkListWrappedResource = new(PrivateDnsZoneVirtualNetworkLinkListResource)

func (r PrivateDnsZoneVirtualNetworkLinkListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePrivateDnsZoneVirtualNetworkLink()
}

func (r PrivateDnsZoneVirtualNetworkLinkListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_private_dns_zone_virtual_network_link"
}
func (r PrivateDnsZoneVirtualNetworkLinkListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"private_dns_zone_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: virtualnetworklinks.ValidatePrivateDnsZoneID,
					},
				},
			},
		},
	}
}

func (r PrivateDnsZoneVirtualNetworkLinkListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDns.VirtualNetworkLinksClient
	var data PrivateDnsZoneVirtualNetworkLinkListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]virtualnetworklinks.VirtualNetworkLink, 0)

	if !data.PrivateDnsZoneId.IsNull() {
		privatednszoneId, err := virtualnetworklinks.ParsePrivateDnsZoneID(data.PrivateDnsZoneId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing privatedns PrivateDnsZone ID for `%s`", "azurerm_private_dns_zone_virtual_network_link"), err)
			return
		}

		resp, err := client.ListComplete(ctx, *privatednszoneId, virtualnetworklinks.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_private_dns_zone_virtual_network_link"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, privatednszonevirtualnetworklink := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(privatednszonevirtualnetworklink.Name)

			rd := resourcePrivateDnsZoneVirtualNetworkLink().Data(&terraform.InstanceState{})

			id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(pointer.From(privatednszonevirtualnetworklink.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing PrivateDnsZoneVirtualNetworkLink ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourcePrivateDnsZoneVirtualNetworkLinkFlatten(rd, id, &privatednszonevirtualnetworklink); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_private_dns_zone_virtual_network_link"), err)
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
