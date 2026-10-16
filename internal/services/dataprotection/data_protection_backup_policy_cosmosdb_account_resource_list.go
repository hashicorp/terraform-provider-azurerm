// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupPolicyCosmosdbAccountListModel struct {
	DataProtectionBackupVaultId types.String `tfsdk:"data_protection_backup_vault_id"`
}

type DataProtectionBackupPolicyCosmosdbAccountListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(DataProtectionBackupPolicyCosmosdbAccountListResource)

func (DataProtectionBackupPolicyCosmosdbAccountListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(DataProtectionBackupPolicyCosmosdbAccountResource{})
}

func (DataProtectionBackupPolicyCosmosdbAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = DataProtectionBackupPolicyCosmosdbAccountResource{}.ResourceType()
}

func (DataProtectionBackupPolicyCosmosdbAccountListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_protection_backup_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: basebackuppolicyresources.ValidateBackupVaultID,
					},
				},
			},
		},
	}
}

func (DataProtectionBackupPolicyCosmosdbAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataProtection.BackupPolicyClient20260601

	var data DataProtectionBackupPolicyCosmosdbAccountListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := basebackuppolicyresources.ParseBackupVaultID(data.DataProtectionBackupVaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `data_protection_backup_vault_id`", err)
		return
	}

	resp, err := client.BackupPoliciesListComplete(ctx, *vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", DataProtectionBackupPolicyCosmosdbAccountResource{}.ResourceType()), err)
		return
	}

	r := DataProtectionBackupPolicyCosmosdbAccountResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, policy := range resp.Items {
			properties, ok := policy.Properties.(basebackuppolicyresources.BackupPolicy)
			if !ok || !backupPolicySupportsCosmosdbAccounts(properties) {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(policy.Name)

			id, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(pointer.From(policy.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Data Protection Backup Policy ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, id, &policy); err != nil {
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

func backupPolicySupportsCosmosdbAccounts(input basebackuppolicyresources.BackupPolicy) bool {
	for _, datasourceType := range input.DatasourceTypes {
		if strings.EqualFold(datasourceType, "Microsoft.DocumentDB/databaseAccounts") {
			return true
		}
	}
	return false
}
