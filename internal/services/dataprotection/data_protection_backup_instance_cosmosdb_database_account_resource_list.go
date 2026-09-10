// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

type DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource struct{}

type DataProtectionBackupInstanceCosmosDBDatabaseAccountListModel struct {
	VaultId types.String `tfsdk:"vault_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource)

func (DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}.ResourceType()
}

func (DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{})
}

func (DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: backupinstances.ValidateBackupVaultID},
				},
			},
		},
	}
}

func (DataProtectionBackupInstanceCosmosDBDatabaseAccountListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var data DataProtectionBackupInstanceCosmosDBDatabaseAccountListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := backupinstances.ParseBackupVaultID(data.VaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `vault_id` for `%s`", DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}.ResourceType()), err)
		return
	}

	resp, err := metadata.Client.DataProtection.BackupInstancesClient20260601.ListComplete(ctx, *vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}.ResourceType()), err)
		return
	}

	resource := DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			if item.Properties == nil || !strings.EqualFold(pointer.From(item.Properties.DataSourceInfo.DatasourceType), "Microsoft.DocumentDB/databaseAccounts") {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := backupinstanceresources.ParseBackupInstanceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Data Protection Backup Instance ID", err)
				return
			}

			model, err := convertDataProtectionModel[backupinstanceresources.BackupInstanceResource](item)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "converting Data Protection Backup Instance", err)
				return
			}

			resourceMetadata := sdk.NewResourceMetaData(metadata.Client, resource)
			resourceMetadata.SetID(id)
			if err := resource.flatten(resourceMetadata, id, model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resource.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, resourceMetadata.ResourceData, &result)
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
