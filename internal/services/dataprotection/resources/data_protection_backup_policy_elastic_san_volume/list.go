// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

import (
	"context"
	"fmt"

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

type ListModel struct {
	DataProtectionBackupVaultId types.String `tfsdk:"data_protection_backup_vault_id"`
}

type ListResource struct{}

var _ sdk.FrameworkListWrappedResource = ListResource{}

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
						Func: basebackuppolicyresources.ValidateBackupVaultID,
					},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var data ListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	vaultId, err := basebackuppolicyresources.ParseBackupVaultIDInsensitively(data.DataProtectionBackupVaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `data_protection_backup_vault_id`", err)
		return
	}
	resp, err := metadata.Client.DataProtection.BackupPolicyClient20260601.BackupPoliciesListComplete(ctx, *vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", Resource{}.ResourceType()), err)
		return
	}

	r := Resource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, policy := range resp.Items {
			if !backupPolicySupportsElasticSanVolumes(policy.Properties) {
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
