// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_virtual_network_swift_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceVirtualNetworkSwiftConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkSwiftConnectionID(d.Id())
	if err != nil {
		return err
	}
	appID := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

	existing, err := client.GetSwiftVirtualNetworkConnection(ctx, appID)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Swift Network Connection for %s: %w", appID, err)
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

	return nil
}
