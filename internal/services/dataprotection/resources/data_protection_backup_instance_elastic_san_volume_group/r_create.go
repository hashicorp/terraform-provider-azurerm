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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceElasticSanVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			vaultId, err := backupvaultresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}
			id := backupinstanceresources.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupInstancesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			volumeGroupId, err := volumegroups.ParseVolumeGroupID(model.ElasticSanVolumeGroupId)
			if err != nil {
				return err
			}
			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			elasticSanId := volumegroups.NewElasticSanID(volumeGroupId.SubscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName)
			snapshotResourceGroupId := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, model.SnapshotResourceGroupName)
			resourceLocation := pointer.To(location.Normalize(model.Location))
			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       volumeGroupId.ID(),
						ResourceLocation: resourceLocation,
						ResourceName:     pointer.To(volumeGroupId.VolumeGroupName),
						ResourceType:     pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ResourceUri:      pointer.To(volumeGroupId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       elasticSanId.ID(),
						ResourceLocation: resourceLocation,
						ResourceName:     pointer.To(elasticSanId.ElasticSanName),
						ResourceType:     pointer.To("Microsoft.ElasticSan/elasticSans"),
						ResourceUri:      pointer.To(elasticSanId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					ObjectType:   "BackupInstance",
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstanceresources.PolicyParameters{
							BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
								backupinstanceresources.GenericBackupDatasourceParameters{
									ObjectType:        "GenericBackupDatasourceParameters",
									ResourceSelectors: []string{model.VolumeName},
								},
							},
							DataStoreParametersList: &[]backupinstanceresources.DataStoreParameters{
								backupinstanceresources.AzureOperationalStoreParameters{
									ResourceGroupId: pointer.To(snapshotResourceGroupId.ID()),
									DataStoreType:   backupinstanceresources.DataStoreTypesOperationalStore,
									ObjectType:      "AzureOperationalStoreParameters",
								},
							},
						},
					},
				},
			}

			if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return pollElasticSanVolumeGroupBackupInstance(ctx, client, id, backupinstanceresources.CurrentProtectionStateConfiguringProtection)
		},
	}
}
