package dataprotection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(DataProtectionBackupVaultListResource)

func (DataProtectionBackupVaultListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureDataProtectionBackupVaultResourceName
}

func (DataProtectionBackupVaultListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceDataProtectionBackupVault()
}

func (DataProtectionBackupVaultListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataProtection.BackupVaultClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []backupvaultresources.BackupVaultResource
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	resp, err := client.BackupVaultsGetInSubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureDataProtectionBackupVaultResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceDataProtectionBackupVault().Data(&terraform.InstanceState{})
			id, err := backupvaultresources.ParseBackupVaultIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureDataProtectionBackupVaultResourceName), err)
				return
			}
			if !data.ResourceGroupName.IsNull() && !strings.EqualFold(id.ResourceGroupName, data.ResourceGroupName.ValueString()) {
				continue
			}
			rd.SetId(id.ID())

			if err := resourceDataProtectionBackupVaultFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureDataProtectionBackupVaultResourceName), err)
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
