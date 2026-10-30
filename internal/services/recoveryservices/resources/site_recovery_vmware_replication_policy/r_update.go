// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryReplicationPolicyVmwareModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			recoveryPoint := plan.RecoveryPointRetentionInMinutes
			appConsistency := plan.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}

			parameters := replicationpolicies.UpdatePolicyInput{
				Properties: &replicationpolicies.UpdatePolicyInputProperties{
					ReplicationProviderSettings: &replicationpolicies.InMageRcmPolicyCreationInput{
						RecoveryPointHistoryInMinutes:     &recoveryPoint,
						AppConsistentFrequencyInMinutes:   &appConsistency,
						EnableMultiVMSync:                 pointer.To(EnableMultiVMSyncEnabled),
						CrashConsistentFrequencyInMinutes: pointer.To(int64(10)),
					},
				},
			}
			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}
