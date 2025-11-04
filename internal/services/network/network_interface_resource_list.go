package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkinterfaces"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkInterfaceListResource struct{}

type NetworkInterfaceListModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
}

var _ sdk.FrameworkListWrappedResource = new(NetworkInterfaceListResource)

func (r NetworkInterfaceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = networkInterfaceResourceName
}

func (r NetworkInterfaceListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceNetworkInterface()
}

func (r NetworkInterfaceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.NetworkInterfaces

	ctx, cancel := context.WithTimeout(ctx, time.Minute*60)
	defer cancel()

	var data NetworkInterfaceListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]networkinterfaces.NetworkInterface, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", networkInterfaceResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", networkInterfaceResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, ni := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(ni.Name)

			id, err := commonids.ParseNetworkInterfaceID(pointer.From(ni.Id))
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "parsing Network Interface ID", err)
				return
			}

			rd := resourceNetworkInterface().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceNetworkInterfaceFlatten(rd, id, &ni); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("encoding `%s` resource data", networkInterfaceResourceName), err)
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
