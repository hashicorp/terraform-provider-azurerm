// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryReplicationPolicyVmwareModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			parsedVaultId, err := vaults.ParseVaultID(plan.RecoveryVaultID)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.RecoveryVaultID, err)
			}

			id := replicationpolicies.NewReplicationPolicyID(parsedVaultId.SubscriptionId, parsedVaultId.ResourceGroupName, parsedVaultId.VaultName, plan.Name)

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
					if !response.WasNotFound(existing.HttpResponse) && !helpers.WasBadRequestWithNotExist(existing.HttpResponse, err) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_site_recovery_vmware_replication_policy", id.ID())
				}
			}

			recoveryPoint := plan.RecoveryPointRetentionInMinutes
			appConsistency := plan.ApplicationConsistentSnapshotFrequencyInMinutes
			if appConsistency > recoveryPoint {
				return fmt.Errorf("the value of `application_consistent_snapshot_frequency_in_minutes` must be less than or equal to the value of `recovery_point_retention_in_minutes`")
			}

			parameters := replicationpolicies.CreatePolicyInput{
				Properties: &replicationpolicies.CreatePolicyInputProperties{
					ProviderSpecificInput: &replicationpolicies.InMageRcmPolicyCreationInput{
						RecoveryPointHistoryInMinutes:     &recoveryPoint,
						AppConsistentFrequencyInMinutes:   &appConsistency,
						CrashConsistentFrequencyInMinutes: pointer.To(int64(10)),
						EnableMultiVMSync:                 pointer.To(EnableMultiVMSyncEnabled),
					},
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
