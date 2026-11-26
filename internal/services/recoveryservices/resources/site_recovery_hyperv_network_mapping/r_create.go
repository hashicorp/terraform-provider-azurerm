// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_network_mapping

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan HyperVNetworkMappingModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.NetworkMappingClient
			fabricClient := metadata.Client.RecoveryServices.FabricClient
			networksClient := metadata.Client.RecoveryServices.ReplicationNetworksClient

			vaultId, err := replicationnetworkmappings.ParseVaultID(plan.VaultId)
			if err != nil {
				return fmt.Errorf("parsing vault id: %+v", err)
			}

			fabricId, err := fetchFabricIdByFriendlyName(ctx, fabricClient, vaultId.ID(), plan.SCVMMname)
			if err != nil {
				return fmt.Errorf("fetching fabric id: %+v", err)
			}

			networkId, err := fetchReplicationNetworkIdByFriendlyName(ctx, networksClient, fabricId, plan.NetworkName)
			if err != nil {
				return fmt.Errorf("fetching network id: %+v", err)
			}

			parsedNetworkId, err := replicationnetworkmappings.ParseReplicationNetworkID(networkId)
			if err != nil {
				return fmt.Errorf("parsing network id: %+v", err)
			}

			id := replicationnetworkmappings.NewReplicationNetworkMappingID(parsedNetworkId.SubscriptionId, parsedNetworkId.ResourceGroupName, parsedNetworkId.VaultName, parsedNetworkId.ReplicationFabricName, parsedNetworkId.ReplicationNetworkName, plan.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking presence of network mapping: %+v", err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_site_recovery_hyperv_network_mapping", id.ID())
				}
			}

			payload := replicationnetworkmappings.CreateNetworkMappingInput{
				Properties: replicationnetworkmappings.CreateNetworkMappingInputProperties{
					RecoveryNetworkId:     plan.TargetNetworkId,
					RecoveryFabricName:    pointer.To(HyperVNetworkMappingRecoveryFabricName),
					FabricSpecificDetails: replicationnetworkmappings.VMmToAzureCreateNetworkMappingInput{},
				},
			}

			if err = client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
