// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(KeyVaultListResource)

func (r KeyVaultListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceKeyVault()
}

func (r KeyVaultListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = keyVaultResourceName
}

func (r KeyVaultListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.KeyVault.VaultsClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]vaults.Vault, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), vaults.ListByResourceGroupOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", keyVaultResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID), vaults.ListBySubscriptionOperationOptions{})
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", keyVaultResourceName), err)
			return
		}

		results = resp.Items
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, vault := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(vault.Name)

			id, err := commonids.ParseKeyVaultIDInsensitively(pointer.From(vault.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Key Vault ID", err)
				return
			}

			rd := resourceKeyVault().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceKeyVaultFlatten(ctx, metadata.Client.KeyVault.ManagementClient, rd, id, &vault, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", keyVaultResourceName), err)
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
