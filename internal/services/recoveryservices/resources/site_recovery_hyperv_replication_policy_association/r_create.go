// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (h Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan HyperVReplicationPolicyAssociationModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			parsedFabricId, err := replicationfabrics.ParseReplicationFabricID(plan.FabricId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.FabricId, err)
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			containerId, err := fetchHyperVReplicationPolicyAssociationContainerNameByHostName(ctx, metadata.Client.RecoveryServices.ProtectionContainerClient, *parsedFabricId)
			if err != nil {
				return fmt.Errorf("fetching container id: %+v", err)
			}

			parsedContainerId, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", containerId, err)
			}

			id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID(subscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, plan.Name)

			param := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
					PolicyId:                    &plan.PolicyId,
					TargetProtectionContainerId: pointer.To(TargetContainerIdAzure),
					ProviderSpecificInput:       replicationprotectioncontainermappings.BaseReplicationProviderSpecificContainerMappingInputImpl{},
				},
			}
			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
