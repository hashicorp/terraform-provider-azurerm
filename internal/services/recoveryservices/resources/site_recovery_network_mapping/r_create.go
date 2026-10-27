// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_network_mapping

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryNetworkMappingCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	targetFabricName := d.Get("target_recovery_fabric_name").(string)
	sourceNetworkId := d.Get("source_network_id").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// get network name from id
	parsedSourceNetworkId, err := commonids.ParseVirtualNetworkID(sourceNetworkId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_network_id '%s' (network mapping %s): %+v", sourceNetworkId, name, err)
	}

	id := replicationnetworkmappings.NewReplicationNetworkMappingID(subscriptionId, resGroup, vaultName, fabricName, parsedSourceNetworkId.VirtualNetworkName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) && !helpers.WasBadRequestWithNotExist(existing.HttpResponse, err) {
				return fmt.Errorf("checking for presence of existing site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_network_mapping", *existing.Model.Id)
		}
	}

	parameters := replicationnetworkmappings.CreateNetworkMappingInput{
		Properties: replicationnetworkmappings.CreateNetworkMappingInputProperties{
			RecoveryNetworkId:  targetNetworkId,
			RecoveryFabricName: &targetFabricName,
			FabricSpecificDetails: replicationnetworkmappings.AzureToAzureCreateNetworkMappingInput{
				PrimaryNetworkId: sourceNetworkId,
			},
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating site recovery network mapping %s (vault %s): %+v", name, vaultName, err)
	}
	d.SetId(id.ID())

	return resourceSiteRecoveryNetworkMappingRead(d, meta)
}
