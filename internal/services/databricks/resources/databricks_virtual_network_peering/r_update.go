// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksVirtualNetworkPeeringUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
	}

	// these are the only updatable values, so everything else in the existing.Model should still be unchanged...
	if d.HasChange("allow_forwarded_traffic") {
		existing.Model.Properties.AllowForwardedTraffic = pointer.To(d.Get("allow_forwarded_traffic").(bool))
	}

	if d.HasChange("allow_gateway_transit") {
		existing.Model.Properties.AllowGatewayTransit = pointer.To(d.Get("allow_gateway_transit").(bool))
	}

	if d.HasChange("allow_virtual_network_access") {
		existing.Model.Properties.AllowVirtualNetworkAccess = pointer.To(d.Get("allow_virtual_network_access").(bool))
	}

	if d.HasChange("use_remote_gateways") {
		existing.Model.Properties.UseRemoteGateways = pointer.To(d.Get("use_remote_gateways").(bool))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating Databricks %s: %+v", *id, err)
	}

	return resourceDatabricksVirtualNetworkPeeringRead(d, meta)
}
