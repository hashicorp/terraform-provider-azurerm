// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_virtual_network_swift_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotVirtualNetworkSwiftConnectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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

	if _, err := client.DeleteSwiftVirtualNetworkSlot(ctx, appSlotID); err != nil {
		return fmt.Errorf("deleting Swift Network Connection for %s: %w", appSlotID, err)
	}

	return nil
}
