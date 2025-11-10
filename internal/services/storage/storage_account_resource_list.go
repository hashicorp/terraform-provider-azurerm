// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.FrameworkListWrappedResource = &StorageAccountListResource{}

type StorageAccountListResource struct{}

func (r StorageAccountListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceStorageAccount()
}

func (r StorageAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = storageAccountResourceName
}

func (r StorageAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	storageClient := metadata.Client.Storage.ResourceManager
	client := storageClient.StorageAccounts

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	listResults := make([]storageaccounts.StorageAccount, 0)
	subscriptionID := metadata.SubscriptionId
	if data.SubscriptionId.ValueString() != "" {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case data.ResourceGroupName.ValueString() != "":
		resourceGroupId := commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString())
		resp, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", storageAccountResourceName), err)
			return
		}

		listResults = resp.Items

	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", storageAccountResourceName), err)
			return
		}

		listResults = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		// TODO - Do we need to handle limiting the results to ListRequest.Limit?
		variableTimeout := time.Duration(5*len(listResults)) * time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), variableTimeout)
		defer cancel()
		for _, account := range listResults {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(account.Name)
			id, err := commonids.ParseStorageAccountID(*account.Id)
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Storage Account ID", err)
				return
			}

			saResource := resourceStorageAccount()

			rd := saResource.Data(&terraform.InstanceState{})

			rd.SetId(id.ID())

			if err := resourceStorageAccountFlatten(ctx, rd, *id, pointer.To(account), metadata.Client); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "encoding Resource data", err)
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity data", err)
				return
			}

			tfTypeResource, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State data", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
