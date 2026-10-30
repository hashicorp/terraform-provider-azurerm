// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (h Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ContainerMappingClient

			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			fabricId := replicationfabrics.NewReplicationFabricID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := HyperVReplicationPolicyAssociationModel{
				Name:     id.ReplicationProtectionContainerMappingName,
				FabricId: fabricId.ID(),
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `Properties` was nil", id)
			}

			prop := resp.Model.Properties
			if prop.PolicyId != nil {
				state.PolicyId = helpers.HandleAzureSdkForGoBug2824(*prop.PolicyId)
			}

			return metadata.Encode(&state)
		},
	}
}

const TargetContainerIdAzure = "Microsoft Azure"

func fetchHyperVReplicationPolicyAssociationContainerNameByHostName(ctx context.Context, containerClient *replicationprotectioncontainers.ReplicationProtectionContainersClient, fabricId replicationfabrics.ReplicationFabricId) (string, error) {
	id, err := replicationprotectioncontainers.ParseReplicationFabricID(fabricId.ID())
	if err != nil {
		return "", fmt.Errorf("parsing %s: %+v", fabricId.ID(), err)
	}

	resp, err := containerClient.ListByReplicationFabricsComplete(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("listing containers: %+v", err)
	}

	if len(resp.Items) == 0 || len(resp.Items) > 1 {
		return "", fmt.Errorf("expected one container but got %d", len(resp.Items))
	}

	if resp.Items[0].Id == nil {
		return "", fmt.Errorf("container id is nil")
	}

	return helpers.HandleAzureSdkForGoBug2824(*resp.Items[0].Id), nil
}
