// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_virtual_network_swift_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotVirtualNetworkSwiftConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.SlotVirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}

	appSlotID := webapps.NewSlotID(id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName)

	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, appSlotID)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appSlotID, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil || existing.Model.Properties.SubnetResourceId == nil {
		d.SetId("")
		return nil
	}

	subnetID, err := commonids.ParseSubnetID(*existing.Model.Properties.SubnetResourceId)
	if err != nil {
		return err
	}

	d.Set("subnet_id", subnetID.ID())
	d.Set("app_service_id", appID.ID())
	d.Set("slot_name", id.SlotName)

	return nil
}
