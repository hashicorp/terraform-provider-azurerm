// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/pools"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HorizondbClusterListResource struct{}

var _ sdk.FrameworkListWrappedResource = HorizondbClusterListResource{}

func (HorizondbClusterListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(HorizondbClusterResource{})
}

func (HorizondbClusterListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = HorizondbClusterResource{}.ResourceType()
}

func (HorizondbClusterListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var config sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &config); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	subscriptionId := metadata.SubscriptionId
	if !config.SubscriptionId.IsNull() {
		subscriptionId = config.SubscriptionId.ValueString()
	}
	client := metadata.Client.HorizonDB.ClustersClient
	var items []clusters.Cluster
	if !config.ResourceGroupName.IsNull() {
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionId, config.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, "listing HorizonDB Clusters", err)
			return
		}
		items = resp.Items
	} else {
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionId))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, "listing HorizonDB Clusters", err)
			return
		}
		items = resp.Items
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range items {
			result := request.NewListResult(ctx)
			id, err := clusters.ParseClusterIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing HorizonDB Cluster ID", err)
				return
			}
			result.DisplayName = id.ClusterName
			resourceMetadata := sdk.NewResourceMetaData(metadata.Client, HorizondbClusterResource{})
			resourceMetadata.ResourceData.SetId(id.ID())
			var pool *pools.Pool
			var administratorLogin string
			if request.IncludeResource {
				pool, err = readHorizondbClusterPool(ctx, resourceMetadata, *id, &item)
				if err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "retrieving HorizonDB Cluster Pool", err)
					return
				}
				administratorLogin, err = readHorizondbClusterAdministratorLogin(ctx, metadata.Client.HorizonDB.PasswordAuthenticationsClient, *id)
				if err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "retrieving HorizonDB Cluster administrator login", err)
					return
				}
			}
			if err := (HorizondbClusterResource{}).flatten(resourceMetadata, id, &item, pool, administratorLogin); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", HorizondbClusterResource{}.ResourceType()), err)
				return
			}
			sdk.EncodeListResult(ctx, resourceMetadata.ResourceData, &result)
			if !push(result) || result.Diagnostics.HasError() {
				return
			}
		}
	}
}
