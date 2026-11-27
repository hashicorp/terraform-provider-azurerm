// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_virtual_network_swift_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceVirtualNetworkSwiftConnectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil || existing.Model.Properties.SubnetResourceId == nil {
		// assume deleted
		return nil
	}

	subnetID, err := commonids.ParseSubnetID(*existing.Model.Properties.SubnetResourceId)
	if err != nil {
		return err
	}

	locks.ByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetID.SubnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetID.SubnetName, network.SubnetResourceName)

	if _, err := client.DeleteSwiftVirtualNetwork(ctx, appID); err != nil {
		return fmt.Errorf("deleting Swift Network Connection for %s: %w", id, err)
	}

	return nil
}
