// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_cosmosdb_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstanceCosmosdbAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			parameters := *existing.Model
			if parameters.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if metadata.ResourceData.HasChange("data_protection_backup_policy_cosmosdb_account_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// the built-in poller is for the LRO to be finished, but the service requires additional time to finish the configure for backup
			// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/41986
			pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, *id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
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
