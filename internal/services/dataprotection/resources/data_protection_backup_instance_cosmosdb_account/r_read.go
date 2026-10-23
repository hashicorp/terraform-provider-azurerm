// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_cosmosdb_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
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

			var model *BackupInstanceCosmosdbAccountModel
			if resp.Model != nil && resp.Model.Properties != nil {
				props := resp.Model.Properties
				model = &BackupInstanceCosmosdbAccountModel{
					Location:          location.NormalizeNilable(props.DataSourceInfo.ResourceLocation),
					CosmosdbAccountId: props.DataSourceInfo.ResourceID,
					BackupPolicyId:    props.PolicyInfo.PolicyId,
					ProtectionState:   pointer.FromEnum(props.CurrentProtectionState),
				}
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *backupinstanceresources.BackupInstanceId, model *BackupInstanceCosmosdbAccountModel) error {
	state := BackupInstanceCosmosdbAccountModel{
		Name:                        id.BackupInstanceName,
		DataProtectionBackupVaultId: backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID(),
	}

	if model != nil {
		state.Location = model.Location

		accountId, err := cosmosdb.ParseDatabaseAccountIDInsensitively(model.CosmosdbAccountId)
		if err != nil {
			return err
		}
		state.CosmosdbAccountId = accountId.ID()

		backupPolicyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(model.BackupPolicyId)
		if err != nil {
			return err
		}
		state.BackupPolicyId = backupPolicyId.ID()

		state.ProtectionState = model.ProtectionState
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
