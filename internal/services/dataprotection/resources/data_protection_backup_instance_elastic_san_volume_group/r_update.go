// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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

			var model BackupInstanceElasticSanVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := *existing.Model
			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}
			if metadata.ResourceData.HasChange("volume_name") {
				if parameters.Properties.PolicyInfo.PolicyParameters == nil {
					parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{}
				}
				parameters.Properties.PolicyInfo.PolicyParameters.BackupDatasourceParametersList = &[]backupinstanceresources.BackupDatasourceParameters{
					backupinstanceresources.GenericBackupDatasourceParameters{
						ObjectType:        "GenericBackupDatasourceParameters",
						ResourceSelectors: []string{model.VolumeName},
					},
				}
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return pollElasticSanVolumeGroupBackupInstance(ctx, client, *id, backupinstanceresources.CurrentProtectionStateUpdatingProtection)
		},
	}
}
