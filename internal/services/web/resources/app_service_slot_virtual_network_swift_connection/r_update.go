// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_virtual_network_swift_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotVirtualNetworkSwiftConnectionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient

	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.SlotVirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	appSlotID := webapps.NewSlotID(id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName)

	existing, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, appSlotID)
	if err != nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appSlotID, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: model was nil", appSlotID)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Swift Network Connection for %s: properties was nil", appSlotID)
	}

	subnetID, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}
	existing.Model.Properties.SubnetResourceId = pointer.To(subnetID.ID())

	if _, err = client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckSlot(ctx, appSlotID, *existing.Model); err != nil {
		return fmt.Errorf("updating association between %s and %s: %w", appSlotID, subnetID, err)
	}

	pollerType := custompollers.NewVirtualNetworkAndSubnetProvisioningSucceededPoller(networkClient.From(meta), subnetID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for completion of association between %s and %s: %w", appSlotID, subnetID, err)
	}

	return resourceAppServiceSlotVirtualNetworkSwiftConnectionRead(d, meta)
}
