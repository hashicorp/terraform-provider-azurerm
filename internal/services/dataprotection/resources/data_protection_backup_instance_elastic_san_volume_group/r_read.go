// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := resp.Model.Properties
			model := BackupInstanceElasticSanVolumeGroupModel{
				Location:                location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation),
				ElasticSanVolumeGroupId: properties.DataSourceInfo.ResourceID,
				BackupPolicyId:          properties.PolicyInfo.PolicyId,
				ProtectionState:         pointer.FromEnum(properties.CurrentProtectionState),
			}
			if policyParameters := properties.PolicyInfo.PolicyParameters; policyParameters != nil {
				if dataStoreParameters := policyParameters.DataStoreParametersList; dataStoreParameters != nil {
					for _, item := range *dataStoreParameters {
						if parameter, ok := item.(backupinstanceresources.AzureOperationalStoreParameters); ok && parameter.ResourceGroupId != nil {
							resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*parameter.ResourceGroupId)
							if err != nil {
								return err
							}
							model.SnapshotResourceGroupName = resourceGroupId.ResourceGroupName
							break
						}
					}
				}
				if datasourceParameters := policyParameters.BackupDatasourceParametersList; datasourceParameters != nil {
					for _, item := range *datasourceParameters {
						if parameter, ok := item.(backupinstanceresources.GenericBackupDatasourceParameters); ok {
							if len(parameter.ResourceSelectors) != 1 {
								return fmt.Errorf("flattening Elastic SAN volume selector: expected one selector but got %d", len(parameter.ResourceSelectors))
							}
							model.VolumeName = parameter.ResourceSelectors[0]
							break
						}
					}
				}
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *backupinstanceresources.BackupInstanceId, model BackupInstanceElasticSanVolumeGroupModel) error {
	model.Name = id.BackupInstanceName
	model.VaultId = backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID()

	volumeGroupId, err := volumegroups.ParseVolumeGroupIDInsensitively(model.ElasticSanVolumeGroupId)
	if err != nil {
		return err
	}
	model.ElasticSanVolumeGroupId = volumeGroupId.ID()

	policyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(model.BackupPolicyId)
	if err != nil {
		return err
	}
	model.BackupPolicyId = policyId.ID()

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&model)
}
