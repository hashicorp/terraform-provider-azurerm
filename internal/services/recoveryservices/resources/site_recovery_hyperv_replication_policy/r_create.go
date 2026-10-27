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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan ReplicationPolicyHyperVModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.ReplicationPoliciesClient

			parsedVaultId, err := vaults.ParseVaultID(plan.RecoveryVaultId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.RecoveryVaultId, err)
			}

			id := replicationpolicies.NewReplicationPolicyID(parsedVaultId.SubscriptionId, parsedVaultId.ResourceGroupName, parsedVaultId.VaultName, plan.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
					if !response.WasNotFound(existing.HttpResponse) && !helpers.WasBadRequestWithNotExist(existing.HttpResponse, err) {
						return fmt.Errorf("checking presence %s: %+v", plan.Name, err)
					}
				}

				if existing.Model != nil {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := replicationpolicies.CreatePolicyInput{
				Properties: &replicationpolicies.CreatePolicyInputProperties{
					ProviderSpecificInput: &replicationpolicies.HyperVReplicaAzurePolicyInput{
						RecoveryPointHistoryDuration:                  &plan.RecoveryPointRetentionInHours,
						ApplicationConsistentSnapshotFrequencyInHours: &plan.ApplicationConsistentSnapshotFrequencyInHours,
						ReplicationInterval:                           &plan.CopyFrequency,
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

func expandHyperVToAzurePolicyDetail(input *replicationpolicies.Policy) (out *replicationpolicies.HyperVReplicaAzurePolicyDetails, isA2A bool) {
	if input.Properties == nil {
		return nil, false
	}
	if input.Properties.ProviderSpecificDetails == nil {
		return nil, false
	}
	detail, isA2A := input.Properties.ProviderSpecificDetails.(replicationpolicies.HyperVReplicaAzurePolicyDetails)
	if isA2A {
		out = &detail
	}
	return
}
