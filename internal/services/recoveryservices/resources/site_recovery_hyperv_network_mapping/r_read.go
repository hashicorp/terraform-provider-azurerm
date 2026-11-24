// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_network_mapping

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id: %+v", err)
			}

			client := metadata.Client.RecoveryServices.NetworkMappingClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			vaultId := replicationnetworkmappings.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
			state := HyperVNetworkMappingModel{
				Name:    id.ReplicationNetworkMappingName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.TargetNetworkId = pointer.From(props.RecoveryNetworkId)
					state.SCVMMname = pointer.From(props.PrimaryFabricFriendlyName)
					state.NetworkName = pointer.From(props.PrimaryNetworkFriendlyName)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

const HyperVNetworkMappingRecoveryFabricName = "Microsoft Azure"

func fetchFabricIdByFriendlyName(ctx context.Context, fabricClient *replicationfabrics.ReplicationFabricsClient, vaultId, friendlyName string) (string, error) {
	parsedVaultId, err := replicationfabrics.ParseVaultID(vaultId)
	if err != nil {
		return "", fmt.Errorf("parsing vault id: %+v", err)
	}

	fabrics, err := fabricClient.ListComplete(ctx, *parsedVaultId)
	if err != nil {
		return "", fmt.Errorf("listing fabrics: %+v", err)
	}

	for _, fabric := range fabrics.Items {
		if fabric.Properties != nil && fabric.Properties.FriendlyName != nil && *fabric.Properties.FriendlyName == friendlyName && fabric.Id != nil {
			return helpers.HandleAzureSdkForGoBug2824(*fabric.Id), nil
		}
	}

	return "", fmt.Errorf("fabric not found")
}

func fetchReplicationNetworkIdByFriendlyName(ctx context.Context, networksClient *replicationnetworks.ReplicationNetworksClient, fabricId, friendlyName string) (string, error) {
	parsedFabricId, err := replicationnetworks.ParseReplicationFabricID(fabricId)
	if err != nil {
		return "", fmt.Errorf("parsing fabric id: %+v", err)
	}

	networks, err := networksClient.ListByReplicationFabricsComplete(ctx, *parsedFabricId)
	if err != nil {
		return "", fmt.Errorf("listing networks: %+v", err)
	}

	for _, network := range networks.Items {
		if network.Properties != nil && network.Properties.FriendlyName != nil && *network.Properties.FriendlyName == friendlyName {
			return helpers.HandleAzureSdkForGoBug2824(*network.Id), nil
		}
	}

	return "", fmt.Errorf("replication network not found")
}
