// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualwans"
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
	VirtualHubConnectionListResource struct{}
	VirtualHubConnectionListModel    struct {
		VirtualHubId types.String `tfsdk:"virtual_hub_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(VirtualHubConnectionListResource)

func (r VirtualHubConnectionListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceVirtualHubConnection()
}

func (r VirtualHubConnectionListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureVirtualHubConnectionResourceName
}

func (r VirtualHubConnectionListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_hub_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: virtualwans.ValidateVirtualHubID,
					},
				},
			},
		},
	}
}

func (r VirtualHubConnectionListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.VirtualWANs

	var data VirtualHubConnectionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	virtualHubId, err := virtualwans.ParseVirtualHubID(data.VirtualHubId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Virtual Hub ID for `%s`", azureVirtualHubConnectionResourceName), err)
		return
	}

	resp, err := client.HubVirtualNetworkConnectionsListComplete(ctx, *virtualHubId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureVirtualHubConnectionResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceVirtualHubConnection().Data(&terraform.InstanceState{})

			id, err := virtualwans.ParseHubVirtualNetworkConnectionIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Virtual Hub Connection ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceVirtualHubConnectionFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureVirtualHubConnectionResourceName), err)
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
