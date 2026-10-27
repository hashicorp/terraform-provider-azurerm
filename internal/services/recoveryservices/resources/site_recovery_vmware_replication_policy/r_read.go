// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy

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
				return fmt.Errorf("reading %s : %+v", id, err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

			state := SiteRecoveryReplicationPolicyVmwareModel{
				Name:            id.ReplicationPolicyName,
				RecoveryVaultID: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if resp.Model.Properties != nil && resp.Model.Properties.ProviderSpecificDetails != nil {
					if inMageRcm, ok := resp.Model.Properties.ProviderSpecificDetails.(replicationpolicies.InMageRcmPolicyDetails); ok {
						if inMageRcm.RecoveryPointHistoryInMinutes != nil {
							state.RecoveryPointRetentionInMinutes = *inMageRcm.RecoveryPointHistoryInMinutes
						}
						if inMageRcm.AppConsistentFrequencyInMinutes != nil {
							state.ApplicationConsistentSnapshotFrequencyInMinutes = *inMageRcm.AppConsistentFrequencyInMinutes
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const EnableMultiVMSyncEnabled string = "True"
