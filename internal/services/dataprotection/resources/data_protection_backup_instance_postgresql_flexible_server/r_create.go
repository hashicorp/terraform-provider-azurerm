// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_postgresql_flexible_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstancePostgreSQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient20260701

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
						DatasourceType:   pointer.To("Microsoft.DBforPostgreSQL/flexibleServers"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       serverId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(serverId.FlexibleServerName),
						ResourceType:     pointer.To("Microsoft.DBforPostgreSQL/flexibleServers"),
						ResourceUri:      pointer.To(serverId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.DBforPostgreSQL/flexibleServers"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       serverId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(serverId.FlexibleServerName),
						ResourceType:     pointer.To("Microsoft.DBforPostgreSQL/flexibleServers"),
						ResourceUri:      pointer.To(serverId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					ObjectType:   "BackupInstance",
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstanceresources.PolicyParameters{
							BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
								backupinstanceresources.PostgreSqlFlexibleServerBackupDatasourceParameters{
									BackupSolutionType: pointer.ToEnum[backupinstanceresources.BackupSolutionType](model.BackupSolutionType),
								},
							},
						},
					},
				},
			}

			if !features.SixPointOh() {
				if model.BackupSolutionType == "" {
					parameters.Properties.PolicyInfo.PolicyParameters = nil
				}
			}

			if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			// Service will continue to configure the protection after the resource is created and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is configured completely
			pollerType := custompollers.NewDataProtectionBackupInstance20260701Poller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateConfiguringProtection,
			})
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			return nil
		},
	}
}
