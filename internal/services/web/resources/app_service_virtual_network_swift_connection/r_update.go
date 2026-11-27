// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_virtual_network_swift_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceVirtualNetworkSwiftConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appID, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: model was nil", appID)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: model was nil", appID)
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetID.VirtualNetworkName, network.VirtualNetworkResourceName)

	locks.ByName(subnetID.SubnetName, network.SubnetResourceName)
	defer locks.UnlockByName(subnetID.SubnetName, network.SubnetResourceName)

	existing.Model.Properties.SubnetResourceId = pointer.To(d.Get("subnet_id").(string))

	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(ctx, appID, *existing.Model); err != nil {
		return fmt.Errorf("creating association between %s and %s: %w", appID, subnetID, err)
	}

	pollerType := custompollers.NewVirtualNetworkAndSubnetProvisioningSucceededPoller(meta.(*clients.Client).Network, subnetID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for completion of association between %s and %s: %w", appID, subnetID, err)
	}

	return resourceAppServiceVirtualNetworkSwiftConnectionRead(d, meta)
}
