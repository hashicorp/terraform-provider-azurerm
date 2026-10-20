// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_kubernetes_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceKubernatesClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient

			vaultId, err := backupvaultresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}

			id := backupinstanceresources.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupInstancesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			aksId, err := commonids.ParseKubernetesClusterID(model.KubernetesClusterId)
			if err != nil {
				return err
			}

			snapshotResourceGroupId := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, model.SnapshotResourceGroupName)
			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.ContainerService/managedClusters"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       aksId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(aksId.ManagedClusterName),
						ResourceType:     pointer.To("Microsoft.ContainerService/managedClusters"),
						ResourceUri:      pointer.To(aksId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.ContainerService/managedClusters"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       aksId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(aksId.ManagedClusterName),
						ResourceType:     pointer.To("Microsoft.ContainerService/managedClusters"),
						ResourceUri:      pointer.To(aksId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					ObjectType:   "BackupInstance",
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstanceresources.PolicyParameters{
							DataStoreParametersList: &[]backupinstanceresources.DataStoreParameters{
								backupinstanceresources.AzureOperationalStoreParameters{
									ResourceGroupId: pointer.To(snapshotResourceGroupId.ID()),
									DataStoreType:   backupinstanceresources.DataStoreTypesOperationalStore,
								},
							},
							BackupDatasourceParametersList: expandBackupDatasourceParameters(model.BackupDatasourceParameters),
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
			return nil
		},
	}
}

func expandBackupDatasourceParameters(input []BackupDatasourceParameters) *[]backupinstanceresources.BackupDatasourceParameters {
	if len(input) == 0 {
		return nil
	}
	results := make([]backupinstanceresources.BackupDatasourceParameters, 0)
	results = append(results, backupinstanceresources.KubernetesClusterBackupDatasourceParameters{
		ExcludedNamespaces:           pointer.To(input[0].ExcludedNamespaces),
		ExcludedResourceTypes:        pointer.To(input[0].ExcludedResourceTypes),
		IncludeClusterScopeResources: input[0].ClusterScopeResourceEnabled,
		IncludedNamespaces:           pointer.To(input[0].IncludedNamespaces),
		IncludedResourceTypes:        pointer.To(input[0].IncludedResourceTypes),
		LabelSelectors:               pointer.To(input[0].LabelSelectors),
		SnapshotVolumes:              input[0].VolumeSnapshotEnabled,
	})
	return &results
}
