// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageActionsTaskDefinitionListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(StorageActionsTaskDefinitionListResource)

func (StorageActionsTaskDefinitionListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(StorageActionsTaskDefinitionResource{})
}

func (StorageActionsTaskDefinitionListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = StorageActionsTaskDefinitionResource{}.ResourceType()
}

func (StorageActionsTaskDefinitionListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.StorageTasksClient

	var data sdk.DefaultListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	r := StorageActionsTaskDefinitionResource{}

	var results []storagetasks.StorageTask
	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := storagetasks.ParseStorageTaskIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Task ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
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
