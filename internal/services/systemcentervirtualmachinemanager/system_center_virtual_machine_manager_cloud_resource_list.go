// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SystemCenterVirtualMachineManagerCloudListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(SystemCenterVirtualMachineManagerCloudListResource)

func (SystemCenterVirtualMachineManagerCloudListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SystemCenterVirtualMachineManagerCloudResource{}.ResourceType()
}

func (SystemCenterVirtualMachineManagerCloudListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(SystemCenterVirtualMachineManagerCloudResource{})
}

func (SystemCenterVirtualMachineManagerCloudListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.SystemCenterVirtualMachineManager.Clouds

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []clouds.Cloud

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SystemCenterVirtualMachineManagerCloudResource{}.ResourceType()), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", SystemCenterVirtualMachineManagerCloudResource{}.ResourceType()), err)
			return
		}

		results = resp.Items
	}

	r := SystemCenterVirtualMachineManagerCloudResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, cloud := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(cloud.Name)

			id, err := clouds.ParseCloudIDInsensitively(pointer.From(cloud.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing SCVMM Cloud ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &cloud); err != nil {
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
