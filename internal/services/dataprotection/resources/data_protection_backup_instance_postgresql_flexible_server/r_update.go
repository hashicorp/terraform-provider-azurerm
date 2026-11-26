// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_postgresql_flexible_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

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
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			parameters := *existing.Model

			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			// Service will update the protection after the resource is updated and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is updated completely
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{string(backupinstanceresources.CurrentProtectionStateUpdatingProtection)},
				Target:       []string{string(backupinstanceresources.CurrentProtectionStateProtectionConfigured)},
				Refresh:      dataProtectionBackupInstancePostgreSQLFlexibleServerStateRefreshFunc(ctx, client, *id),
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
