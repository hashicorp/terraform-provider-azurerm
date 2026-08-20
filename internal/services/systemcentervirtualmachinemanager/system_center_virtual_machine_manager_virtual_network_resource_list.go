// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SystemCenterVirtualMachineManagerVirtualNetworkListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(SystemCenterVirtualMachineManagerVirtualNetworkListResource)

func (SystemCenterVirtualMachineManagerVirtualNetworkListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SystemCenterVirtualMachineManagerVirtualNetworkResource{}.ResourceType()
}

func (SystemCenterVirtualMachineManagerVirtualNetworkListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(SystemCenterVirtualMachineManagerVirtualNetworkResource{})
}

func (SystemCenterVirtualMachineManagerVirtualNetworkListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SystemCenterVirtualMachineManager.VirtualNetworks

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []virtualnetworks.VirtualNetwork

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SystemCenterVirtualMachineManagerVirtualNetworkResource{}.ResourceType()), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SystemCenterVirtualMachineManagerVirtualNetworkResource{}.ResourceType()), err)
			return
		}

		results = resp.Items
	}

	r := SystemCenterVirtualMachineManagerVirtualNetworkResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, virtualNetwork := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(virtualNetwork.Name)

			id, err := virtualnetworks.ParseVirtualNetworkIDInsensitively(pointer.From(virtualNetwork.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing SCVMM Virtual Network ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &virtualNetwork); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
