// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := protectionpolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if m := existing.Model; m != nil {
				props, _ := m.Properties.(protectionpolicies.AzureVMWorkloadProtectionPolicy)

				if metadata.ResourceData.HasChange("settings") {
					props.Settings = expandBackupProtectionPolicyVMWorkloadSettings(model.Settings)
				}

				if metadata.ResourceData.HasChange("protection_policy") {
					protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies, model.WorkloadType)
					if err != nil {
						return err
					}

					props.SubProtectionPolicy = protectionPolicy
				}

				m.Properties = props

				if _, err := client.CreateOrUpdate(ctx, *id, *m, protectionpolicies.DefaultCreateOrUpdateOperationOptions()); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
