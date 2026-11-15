// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

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

			vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)

			state := BackupInstanceKubernatesClusterModel{
				Name:    id.BackupInstanceName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.Location = location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation)
					state.BackupPolicyId = properties.PolicyInfo.PolicyId
					state.KubernetesClusterId = properties.DataSourceInfo.ResourceID
					state.ProtectionState = pointer.FromEnum(properties.CurrentProtectionState)

					if policyParameters := properties.PolicyInfo.PolicyParameters; policyParameters != nil {
						if dataStorePara := policyParameters.DataStoreParametersList; dataStorePara != nil {
							if dsp := pointer.From(dataStorePara); len(dsp) > 0 {
								if parameter, ok := dsp[0].(backupinstanceresources.AzureOperationalStoreParameters); ok && parameter.ResourceGroupId != nil {
									resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*parameter.ResourceGroupId)
									if err != nil {
										return err
									}
									state.SnapshotResourceGroupName = resourceGroupId.ResourceGroupName
								}
							}
						}
						if backupDsp := policyParameters.BackupDatasourceParametersList; backupDsp != nil {
							if v := flattenBackupDatasourceParameters(*backupDsp); v != nil {
								state.BackupDatasourceParameters = pointer.From(v)
							}
						}
					}
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenBackupDatasourceParameters(input []backupinstanceresources.BackupDatasourceParameters) *[]BackupDatasourceParameters {
	results := make([]BackupDatasourceParameters, 0)
	if len(input) == 0 {
		return &results
	}

	if item, ok := input[0].(backupinstanceresources.KubernetesClusterBackupDatasourceParameters); ok {
		results = append(results, BackupDatasourceParameters{
			ExcludedNamespaces:          pointer.From(item.ExcludedNamespaces),
			ExcludedResourceTypes:       pointer.From(item.ExcludedResourceTypes),
			ClusterScopeResourceEnabled: item.IncludeClusterScopeResources,
			IncludedNamespaces:          pointer.From(item.IncludedNamespaces),
			IncludedResourceTypes:       pointer.From(item.IncludedResourceTypes),
			LabelSelectors:              pointer.From(item.LabelSelectors),
			VolumeSnapshotEnabled:       item.SnapshotVolumes,
		})
	}
	return &results
}

type BackupDatasourceParameters struct {
	IncludedNamespaces          []string `tfschema:"included_namespaces"`
	IncludedResourceTypes       []string `tfschema:"included_resource_types"`
	ExcludedNamespaces          []string `tfschema:"excluded_namespaces"`
	ExcludedResourceTypes       []string `tfschema:"excluded_resource_types"`
	LabelSelectors              []string `tfschema:"label_selectors"`
	VolumeSnapshotEnabled       bool     `tfschema:"volume_snapshot_enabled"`
	ClusterScopeResourceEnabled bool     `tfschema:"cluster_scoped_resources_enabled"`
}
