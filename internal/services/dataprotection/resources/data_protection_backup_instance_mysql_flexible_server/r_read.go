// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_mysql_flexible_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
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

			state := BackupInstanceMySQLFlexibleServerModel{
				Name:    id.BackupInstanceName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.NormalizeNilable(props.DataSourceInfo.ResourceLocation)

					serverId, err := servers.ParseFlexibleServerIDInsensitively(props.DataSourceInfo.ResourceID)
					if err != nil {
						return err
					}
					state.ServerId = serverId.ID()

					backupPolicyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(props.PolicyInfo.PolicyId)
					if err != nil {
						return err
					}
					state.BackupPolicyId = backupPolicyId.ID()

					state.ProtectionState = pointer.FromEnum(props.CurrentProtectionState)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func dataProtectionBackupInstanceMySQLFlexibleServerStateRefreshFunc(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.BackupInstancesGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if resp.Model == nil {
			return nil, "", fmt.Errorf("polling for %s: `model` was nil", id)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("polling for %s: `properties` was nil", id)
		}

		return resp, pointer.FromEnum(resp.Model.Properties.CurrentProtectionState), nil
	}
}
