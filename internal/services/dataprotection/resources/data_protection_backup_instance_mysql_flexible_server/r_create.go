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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceMySQLFlexibleServerModel
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

			serverId, err := servers.ParseFlexibleServerID(model.ServerId)
			if err != nil {
				return err
			}

			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.DBforMySQL/flexibleServers"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       serverId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(serverId.FlexibleServerName),
						ResourceType:     pointer.To("Microsoft.DBforMySQL/flexibleServers"),
						ResourceUri:      pointer.To(serverId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.DBforMySQL/flexibleServers"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       serverId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(serverId.FlexibleServerName),
						ResourceType:     pointer.To("Microsoft.DBforMySQL/flexibleServers"),
						ResourceUri:      pointer.To(serverId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
					},
				},
			}

			if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			// Service will continue to configure the protection after the resource is created and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is configured completely
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{string(backupinstanceresources.CurrentProtectionStateConfiguringProtection)},
				Target:       []string{string(backupinstanceresources.CurrentProtectionStateProtectionConfigured)},
				Refresh:      dataProtectionBackupInstanceMySQLFlexibleServerStateRefreshFunc(ctx, client, id),
				PollInterval: 1 * time.Minute,
				Timeout:      time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			return nil
		},
	}
}
