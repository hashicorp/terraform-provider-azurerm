// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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

type DataProtectionBackupInstanceElasticSanVolumeGroupListModel struct {
	VaultId types.String `tfsdk:"vault_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(ListResource)

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var data DataProtectionBackupInstanceElasticSanVolumeGroupListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vaultId, err := backupinstances.ParseBackupVaultID(data.VaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `vault_id` for `%s`", Resource{}.ResourceType()), err)
		return
	}

	resp, err := metadata.Client.DataProtection.BackupInstancesClient20260601.ListComplete(ctx, *vaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", Resource{}.ResourceType()), err)
		return
	}

	resource := Resource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			if item.Properties == nil || !strings.EqualFold(pointer.From(item.Properties.DataSourceInfo.DatasourceType), "Microsoft.ElasticSan/elasticSans/volumeGroups") {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := backupinstanceresources.ParseBackupInstanceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Data Protection Backup Instance ID", err)
				return
			}

			properties := item.Properties
			model := BackupInstanceElasticSanVolumeGroupModel{
				Location:                location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation),
				ElasticSanVolumeGroupId: properties.DataSourceInfo.ResourceID,
				BackupPolicyId:          properties.PolicyInfo.PolicyId,
				ProtectionState:         pointer.FromEnum(properties.CurrentProtectionState),
			}
			if policyParameters := properties.PolicyInfo.PolicyParameters; policyParameters != nil {
				if dataStoreParameters := policyParameters.DataStoreParametersList; dataStoreParameters != nil {
					for _, item := range *dataStoreParameters {
						if parameter, ok := item.(backupinstances.AzureOperationalStoreParameters); ok && parameter.ResourceGroupId != nil {
							resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*parameter.ResourceGroupId)
							if err != nil {
								sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing snapshot Resource Group ID", err)
								return
							}
							model.SnapshotResourceGroupName = resourceGroupId.ResourceGroupName
							break
						}
					}
				}
				if datasourceParameters := policyParameters.BackupDatasourceParametersList; datasourceParameters != nil {
					for _, item := range *datasourceParameters {
						if parameter, ok := item.(backupinstances.GenericBackupDatasourceParameters); ok {
							if len(parameter.ResourceSelectors) != 1 {
								sdk.SetErrorDiagnosticAndPushListResult(result, push, "flattening Elastic SAN volume selector", fmt.Errorf("flattening Elastic SAN volume selector: expected one selector but got %d", len(parameter.ResourceSelectors)))
								return
							}
							model.VolumeName = parameter.ResourceSelectors[0]
							break
						}
					}
				}
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
