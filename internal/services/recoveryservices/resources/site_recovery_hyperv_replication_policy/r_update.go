// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ReplicationPolicyHyperVModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			parameters := replicationpolicies.UpdatePolicyInput{
				Properties: &replicationpolicies.UpdatePolicyInputProperties{
					ReplicationProviderSettings: &replicationpolicies.HyperVReplicaAzurePolicyInput{
						RecoveryPointHistoryDuration:                  &model.RecoveryPointRetentionInHours,
						ApplicationConsistentSnapshotFrequencyInHours: &model.ApplicationConsistentSnapshotFrequencyInHours,
						ReplicationInterval:                           &model.CopyFrequency,
					},
				},
			}

			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
