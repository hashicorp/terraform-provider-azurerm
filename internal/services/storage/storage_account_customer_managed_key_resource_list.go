package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageAccountCustomerManagedKeyListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(StorageAccountCustomerManagedKeyListResource)

func (StorageAccountCustomerManagedKeyListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceStorageAccountCustomerManagedKey()
}

func (r StorageAccountCustomerManagedKeyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = storageAccountCustomerManagedKeyResourceName
}

func (StorageAccountCustomerManagedKeyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Storage.ResourceManager.StorageAccounts

	// retrieve the deadline from the supplied context
	deadline, ok := ctx.Deadline()
	if !ok {
		// This *should* never happen given the List Wrapper instantiates a context with a timeout
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]storageaccounts.StorageAccount, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", storageAccountCustomerManagedKeyResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", storageAccountCustomerManagedKeyResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, storageAccount := range results {
			// results contain all storage accounts including without CMK
			// skip any accounts that do not have CMK
			if storageAccount.Properties.Encryption == nil || pointer.From(storageAccount.Properties.Encryption.KeySource) != "Microsoft.Keyvault" {
				continue
			}

			result := request.NewListResult(deadlineCtx)
			result.DisplayName = pointer.From(storageAccount.Name)
			rd := resourceStorageAccountCustomerManagedKey().Data(&terraform.InstanceState{})

			id, err := commonids.ParseStorageAccountID(pointer.From(storageAccount.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Storage Account ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceStorageAccountCustomerManagedKeyFlatten(deadlineCtx, metadata.Client, rd, id, &storageAccount, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", storageAccountCustomerManagedKeyResourceName), err)
				return
			}

			sdk.EncodeListResult(deadlineCtx, rd, &result)
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
