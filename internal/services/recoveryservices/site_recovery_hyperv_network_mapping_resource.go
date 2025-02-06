// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const HyperVNetworkMappingRecoveryFabricName = "Microsoft Azure"

type HyperVNetworkMappingResource struct{}

type HyperVNetworkMappingModel struct {
	Name            string `tfschema:"name"`
	VaultId         string `tfschema:"recovery_vault_id"`
	SCVMMname       string `tfschema:"source_system_center_virtual_machine_manager_name"`
	NetworkName     string `tfschema:"source_network_name"`
	TargetNetworkId string `tfschema:"target_network_id"`
}

var _ sdk.Resource = HyperVNetworkMappingResource{}

func (s HyperVNetworkMappingResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},

		"source_system_center_virtual_machine_manager_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_network_name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_network_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     azure.ValidateResourceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func (s HyperVNetworkMappingResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s HyperVNetworkMappingResource) ModelObject() interface{} {
	return &HyperVNetworkMappingModel{}
}

func (s HyperVNetworkMappingResource) ResourceType() string {
	return "azurerm_site_recovery_hyperv_network_mapping"
}

func (s HyperVNetworkMappingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationnetworkmappings.ValidateReplicationNetworkMappingID
}

func (s HyperVNetworkMappingResource) Create() sdk.ResourceFunc {
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
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking presence of network mapping: %+v", err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_site_recovery_hyperv_network_mapping", id.ID())
			}

			err = client.CreateThenPoll(ctx, id, replicationnetworkmappings.CreateNetworkMappingInput{
				Properties: replicationnetworkmappings.CreateNetworkMappingInputProperties{
					RecoveryNetworkId:     plan.TargetNetworkId,
					RecoveryFabricName:    pointer.To(HyperVNetworkMappingRecoveryFabricName),
					FabricSpecificDetails: replicationnetworkmappings.VMmToAzureCreateNetworkMappingInput{},
				},
			})
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s HyperVNetworkMappingResource) Read() sdk.ResourceFunc {
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

func (s HyperVNetworkMappingResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id: %+v", err)
			}

			client := metadata.Client.RecoveryServices.NetworkMappingClient

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

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
			return handleAzureSdkForGoBug2824(*fabric.Id), nil
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
			return handleAzureSdkForGoBug2824(*network.Id), nil
		}
	}

	return "", fmt.Errorf("replication network not found")
}
