// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedHatHcpOpenShiftClusterListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(RedHatHcpOpenShiftClusterListResource)

func (RedHatHcpOpenShiftClusterListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = RedHatHcpOpenShiftCluster{}.ResourceType()
}

func (RedHatHcpOpenShiftClusterListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(RedHatHcpOpenShiftCluster{})
}

func (r RedHatHcpOpenShiftClusterListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient

	var data sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []hcpopenshiftclusters.HcpOpenShiftCluster

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	clusterResource := RedHatHcpOpenShiftCluster{}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", clusterResource.ResourceType()), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", clusterResource.ResourceType()), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, cluster := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(cluster.Name)

			id, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(pointer.From(cluster.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing HCP Open Shift Cluster ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, clusterResource)
			meta.SetID(id)

			var config RedHatHcpOpenShiftClusterModel
			if err := meta.Decode(&config); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "decoding HCP Open Shift Cluster state", err)
				return
			}

			if err := clusterResource.flatten(meta, *id, &cluster); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", clusterResource.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
