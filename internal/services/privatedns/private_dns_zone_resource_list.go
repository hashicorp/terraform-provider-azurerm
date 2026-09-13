package privatedns

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDnsZoneListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(PrivateDnsZoneListResource)

func (r PrivateDnsZoneListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePrivateDnsZone()
}

func (r PrivateDnsZoneListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = privateDnsZoneResourceName
}

func (r PrivateDnsZoneListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDns.PrivateZonesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]privatezones.PrivateZone, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroup(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), privatezones.DefaultListByResourceGroupOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", privateDnsZoneResourceName), err)
			return
		}

		results = pointer.From(resp.Model)
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), privatezones.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", privateDnsZoneResourceName), err)
			return
		}

		results = resp.Items
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, privateZone := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(privateZone.Name)

			id, err := privatezones.ParsePrivateDnsZoneID(pointer.From(privateZone.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Private DNS Zone ID", err)
				return
			}

			rd := resourcePrivateDnsZone().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourcePrivateDnsZoneFlatten(ctx, rd, metadata.Client, id, &privateZone, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", privateDnsZoneResourceName), err)
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
