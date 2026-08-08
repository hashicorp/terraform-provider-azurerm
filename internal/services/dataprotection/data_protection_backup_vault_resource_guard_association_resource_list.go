// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	resourceguardproxy "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultResourceGuardAssociationListResource struct{}

type DataProtectionBackupVaultResourceGuardAssociationListModel struct {
	DataProtectionBackupVaultId types.String `tfsdk:"data_protection_backup_vault_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(DataProtectionBackupVaultResourceGuardAssociationListResource)

func (r DataProtectionBackupVaultResourceGuardAssociationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = dataProtectionBackupVaultResourceGuardAssociationResourceType
}

func (r DataProtectionBackupVaultResourceGuardAssociationListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(DataProtectionBackupVaultResourceGuardAssociationResource{})
}

func (r DataProtectionBackupVaultResourceGuardAssociationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"data_protection_backup_vault_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: backupvaultresources.ValidateBackupVaultID,
					},
				},
			},
		},
	}
}

func (r DataProtectionBackupVaultResourceGuardAssociationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataProtection.ResourceGuardProxyClient

	var data DataProtectionBackupVaultResourceGuardAssociationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := backupvaultresources.ParseBackupVaultID(data.DataProtectionBackupVaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Data Protection Backup Vault ID for `%s`", dataProtectionBackupVaultResourceGuardAssociationResourceType), err)
		return
	}

	proxyVaultId := resourceguardproxy.NewBackupVaultID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName)
	resp, err := client.DppResourceGuardProxyListComplete(ctx, proxyVaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", dataProtectionBackupVaultResourceGuardAssociationResourceType), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, proxy := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(proxy.Name)

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(pointer.From(proxy.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Data Protection Backup Vault Resource Guard Association ID", err)
				return
			}

			rd := r.ResourceFunc().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := setDataProtectionBackupVaultResourceGuardAssociationResourceData(rd, *id, &proxy); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", dataProtectionBackupVaultResourceGuardAssociationResourceType), err)
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
