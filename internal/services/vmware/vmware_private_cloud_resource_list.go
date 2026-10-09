// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/privateclouds"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VmwarePrivateCloudListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(VmwarePrivateCloudListResource)

func (r VmwarePrivateCloudListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceVmwarePrivateCloud()
}

func (r VmwarePrivateCloudListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureVmwarePrivateCloudResourceName
}

func (r VmwarePrivateCloudListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Vmware.PrivateCloudClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]privateclouds.PrivateCloud, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureVmwarePrivateCloudResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListInSubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureVmwarePrivateCloudResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, privateCloud := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(privateCloud.Name)

			id, err := privateclouds.ParsePrivateCloudID(pointer.From(privateCloud.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing VMware Private Cloud ID", err)
				return
			}

			rd := resourceVmwarePrivateCloud().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceVmwarePrivateCloudFlatten(rd, id, &privateCloud); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureVmwarePrivateCloudResourceName), err)
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
