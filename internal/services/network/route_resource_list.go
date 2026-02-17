// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/routes"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RouteListResource struct{}
type RouteListModel struct {
	RouteTableId types.String `tfsdk:"route_table_id"`
}

var _ sdk.FrameworkListWrappedResource = new(RouteListResource)

func (r RouteListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceRoute()
}

func (r RouteListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_route"
}
func (r RouteListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"route_table_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: routes.ValidateRouteTableID,
					},
				},
			},
		},
	}
}

func (r RouteListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.Routes
	var data RouteListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]routes.Route, 0)

	if !data.RouteTableId.IsNull() {
		routetableId, err := routes.ParseRouteTableID(data.RouteTableId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing network RouteTable ID for `%s`", "azurerm_route"), err)
			return
		}

		resp, err := client.ListComplete(ctx, *routetableId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_route"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, route := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(route.Name)

			rd := resourceRoute().Data(&terraform.InstanceState{})

			id, err := routes.ParseRouteID(pointer.From(route.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Network Route ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceRouteFlatten(rd, id, &route); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_route"), err)
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
