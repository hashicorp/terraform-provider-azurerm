// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_network_mapping

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryNetworkMappingRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationnetworkmappings.ParseReplicationNetworkMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.NetworkMappingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery network mapping %q: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving site recovery network mapping %q: `model` was nil", id)
	}
	model := resp.Model

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("source_recovery_fabric_name", id.ReplicationFabricName)
	d.Set("name", id.ReplicationNetworkMappingName)
	if props := model.Properties; props != nil {
		d.Set("source_network_id", props.PrimaryNetworkId)
		d.Set("target_network_id", props.RecoveryNetworkId)

		targetFabricId, err := replicationfabrics.ParseReplicationFabricIDInsensitively(*props.RecoveryFabricArmId)
		if err != nil {
			return err
		}
		d.Set("target_recovery_fabric_name", targetFabricId.ReplicationFabricName)
	}

	return nil
}
