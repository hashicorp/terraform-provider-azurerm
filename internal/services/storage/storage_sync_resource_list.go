// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageSyncListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(StorageSyncListResource)

func (r StorageSyncListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceStorageSync()
}

func (r StorageSyncListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = storageSyncResourceName
}

func (r StorageSyncListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.SyncServiceClient

	var data sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	results := make([]storagesyncservicesresource.StorageSyncService, 0)

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.StorageSyncServicesListByResourceGroup(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", storageSyncResourceName), err)
			return
		}

		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	default:
		resp, err := client.StorageSyncServicesListBySubscription(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", storageSyncResourceName), err)
			return
		}

		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := storagesyncservicesresource.ParseStorageSyncServiceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Sync Service ID", err)
				return
			}

			rd := resourceStorageSync().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceStorageSyncFlatten(ctx, rd, id, &item, metadata.Client.Storage.SyncRegisteredServerClient, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", storageSyncResourceName), err)
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
