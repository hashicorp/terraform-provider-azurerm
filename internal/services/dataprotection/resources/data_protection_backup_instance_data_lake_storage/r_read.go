// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_data_lake_storage

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

			state := BackupInstanceDataLakeStorageModel{
				Name:                        id.BackupInstanceName,
				DataProtectionBackupVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.NormalizeNilable(props.DataSourceInfo.ResourceLocation)

					storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(props.DataSourceInfo.ResourceID)
					if err != nil {
						return err
					}
					state.StorageAccountId = storageAccountId.ID()

					backupPolicyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(props.PolicyInfo.PolicyId)
					if err != nil {
						return err
					}
					state.BackupPolicyId = backupPolicyId.ID()

					state.ProtectionState = pointer.FromEnum(props.CurrentProtectionState)

					if policyParas := props.PolicyInfo.PolicyParameters; policyParas != nil {
						if dataStoreParas := policyParas.BackupDatasourceParametersList; dataStoreParas != nil {
							if dsp := pointer.From(dataStoreParas); len(dsp) > 0 {
								if parameter, ok := dsp[0].(backupinstanceresources.AdlsBlobBackupDatasourceParameters); ok {
									state.StorageContainerNames = parameter.ContainersList
								}
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
