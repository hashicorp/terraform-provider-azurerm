// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_fabric

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryFabricRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationfabrics.ParseReplicationFabricID(d.Id())
	if err != nil {
		return err
	}

	fabricClient := meta.(*clients.Client).RecoveryServices.FabricClient
	client := fabricClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making read request on site recovery fabric %s: %+v", id.String(), err)
	}

	d.Set("name", id.ReplicationFabricName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if details := props.CustomDetails; details != nil {
				fabric, ok := details.(replicationfabrics.AzureFabricSpecificDetails)
				if !ok {
					return fmt.Errorf("expected `details` to be an AzureFabricSpecificDetails but it wasn't: %+v", details)
				}
				d.Set("location", location.NormalizeNilable(fabric.Location))
			}
		}
	}

	return nil
}
