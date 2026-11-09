// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_kubernetes_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
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

			var model BackupInstanceKubernatesClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			parameters := *existing.Model
			if parameters.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_id") {
				parameters.Properties.IdentityDetails = &backupinstanceresources.IdentityDetails{
					UseSystemAssignedIdentity: pointer.To(true),
				}
				if model.UserAssignedIdentityId != "" {
					identityId, err := commonids.ParseUserAssignedIdentityID(model.UserAssignedIdentityId)
					if err != nil {
						return err
					}
					parameters.Properties.IdentityDetails = &backupinstanceresources.IdentityDetails{
						UseSystemAssignedIdentity:  pointer.To(false),
						UserAssignedIdentityArmURL: pointer.To(identityId.ID()),
					}
				}
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			pollerType := custompollers.NewDataProtectionBackupInstancePoller(client, *id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateUpdatingProtection,
			})
			poller := pollers.NewPoller(pollerType, time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			return nil
		},
	}
}
