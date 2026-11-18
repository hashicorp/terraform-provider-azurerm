// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_postgresql_flexible_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260701

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstancePostgreSQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("reading %s: `model` was nil", *id)
			}

			parameters := *existing.Model
			if parameters.Properties == nil {
				return fmt.Errorf("reading %s: `properties` was nil", *id)
			}

			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if metadata.ResourceData.HasChange("backup_solution_type") {
				if parameters.Properties.PolicyInfo.PolicyParameters == nil {
					parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{}
				}
				parameters.Properties.PolicyInfo.PolicyParameters.BackupDatasourceParametersList = &[]backupinstanceresources.BackupDatasourceParameters{
					backupinstanceresources.PostgreSqlFlexibleServerBackupDatasourceParameters{
						BackupSolutionType: pointer.ToEnum[backupinstanceresources.BackupSolutionType](model.BackupSolutionType),
					},
				}
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// Service will update the protection after the resource is updated and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is updated completely
			pollerType := custompollers.NewDataProtectionBackupInstance20260701Poller(client, *id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateConfiguringProtection,
				backupinstanceresources.CurrentProtectionStateUpdatingProtection,
			})
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", *id, err)
			}

			return nil
		},
	}
}
