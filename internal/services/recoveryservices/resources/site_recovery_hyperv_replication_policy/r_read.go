// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationpolicies.ParseReplicationPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s : %+v", *id, err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

			state := ReplicationPolicyHyperVModel{
				Name:            id.ReplicationPolicyName,
				RecoveryVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if detail, isA2A := expandHyperVToAzurePolicyDetail(resp.Model); isA2A {
					if detail.ApplicationConsistentSnapshotFrequencyInHours != nil {
						state.ApplicationConsistentSnapshotFrequencyInHours = *detail.ApplicationConsistentSnapshotFrequencyInHours
					}
					if detail.ReplicationInterval != nil {
						state.CopyFrequency = *detail.ReplicationInterval
					}
					if detail.RecoveryPointHistoryDurationInHours != nil {
						state.RecoveryPointRetentionInHours = *detail.RecoveryPointHistoryDurationInHours
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}
