package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/routetables"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RouteTableListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(RouteTableListResource)

func (r RouteTableListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceRouteTable()
}

func (r RouteTableListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = routeTableResourceName
}

func (r RouteTableListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.RouteTables

	ctx, cancel := context.WithTimeout(ctx, time.Minute*60)
	defer cancel()

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]routetables.RouteTable, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", routeTableResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", routeTableResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, table := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(table.Name)

			id, err := routetables.ParseRouteTableID(pointer.From(table.Id))
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "parsing Route Table ID", err)
				return
			}

			rd := resourceRouteTable().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceRouteTableFlatten(rd, id, &table); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("encoding `%s` resource data", routeTableResourceName), err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting Identity Data", err)
				return
			}

			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Resource State", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting Resource Data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
