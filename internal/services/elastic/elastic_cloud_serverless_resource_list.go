// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticCloudServerlessListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ElasticCloudServerlessListResource)

func (ElasticCloudServerlessListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = ElasticCloudServerlessResource{}.ResourceType()
}

func (ElasticCloudServerlessListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(ElasticCloudServerlessResource{})
}

func (ElasticCloudServerlessListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Elastic.ServerlessMonitorClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var monitors []elasticmonitorresources.ElasticMonitorResource

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.MonitorsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", ElasticCloudServerlessResource{}.ResourceType()), err)
			return
		}

		monitors = resp.Items
	default:
		resp, err := client.MonitorsListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", ElasticCloudServerlessResource{}.ResourceType()), err)
			return
		}

		monitors = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		serverlessResource := ElasticCloudServerlessResource{}
		for _, monitor := range monitors {
			if monitor.Properties == nil || pointer.From(monitor.Properties.HostingType) != elasticmonitorresources.HostingTypeServerless {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(monitor.Name)

			id, err := elasticmonitorresources.ParseMonitorIDInsensitively(pointer.From(monitor.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Elastic Monitor ID", err)
				return
			}

			resourceMetadata := sdk.NewResourceMetaData(metadata.Client, serverlessResource)
			resourceMetadata.SetID(id)

			if err := serverlessResource.flatten(resourceMetadata, id, &monitor); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", serverlessResource.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, resourceMetadata.ResourceData, &result)
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
