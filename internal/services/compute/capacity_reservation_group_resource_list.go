package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CapacityReservationGroupListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CapacityReservationGroupListResource)

func (CapacityReservationGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureCapacityReservationGroupResourceName
}

func (CapacityReservationGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCapacityReservationGroup()
}

func (CapacityReservationGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Compute.CapacityReservationGroupsClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []capacityreservationgroups.CapacityReservationGroup
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), capacityreservationgroups.DefaultListByResourceGroupOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCapacityReservationGroupResourceName), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID), capacityreservationgroups.DefaultListBySubscriptionOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCapacityReservationGroupResourceName), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceCapacityReservationGroup().Data(&terraform.InstanceState{})
			id, err := capacityreservationgroups.ParseCapacityReservationGroupIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureCapacityReservationGroupResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceCapacityReservationGroupFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureCapacityReservationGroupResourceName), err)
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
