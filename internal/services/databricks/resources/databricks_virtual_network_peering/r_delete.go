// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	// Block all changes to any resource of this type...
	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Databricks %s: %+v", *id, err)
	}

	return nil
}
