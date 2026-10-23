// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_cosmosdb_account

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstances"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

type ListModel struct {
	VaultId types.String `tfsdk:"data_protection_backup_vault_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ListResource)

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_protection_backup_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: backupinstances.ValidateBackupVaultID,
					},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataProtection.BackupInstancesClient20260601

	var data ListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := backupinstances.ParseBackupVaultID(data.VaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Data Protection Backup Vault ID", err)
		return
	}

	resp, err := client.ListComplete(ctx, *vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", Resource{}.ResourceType()), err)
		return
	}

	r := Resource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			if item.Properties == nil || !strings.EqualFold(pointer.From(item.Properties.DataSourceInfo.ResourceType), cosmosDBAccountDataSourceType) {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := backupinstanceresources.ParseBackupInstanceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", r.ResourceType()), err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			model := BackupInstanceCosmosdbAccountModel{
				Location:          location.NormalizeNilable(item.Properties.DataSourceInfo.ResourceLocation),
				CosmosdbAccountId: item.Properties.DataSourceInfo.ResourceID,
				BackupPolicyId:    item.Properties.PolicyInfo.PolicyId,
				ProtectionState:   pointer.FromEnum(item.Properties.CurrentProtectionState),
			}
			if err := r.flatten(rmd, id, &model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
