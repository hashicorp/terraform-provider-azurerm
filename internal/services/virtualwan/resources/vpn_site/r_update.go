// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_site

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnSiteUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnSiteID(d.Id())
	if err != nil {
		return err
	}
	existing, err := client.VpnSitesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("address_cidrs") {
		payload.Properties.AddressSpace = expandVpnSiteAddressSpace(d.Get("address_cidrs").(*pluginsdk.Set).List())
	}

	if d.HasChanges("device_vendor", "device_model") {
		payload.Properties.DeviceProperties = expandVpnSiteDeviceProperties(d.Get("device_vendor").(string), d.Get("device_model").(string))
	}

	if d.HasChange("link") {
		payload.Properties.VpnSiteLinks = expandVpnSiteLinks(d.Get("link").([]any))
	}

	if d.HasChange("o365_policy") {
		payload.Properties.O365Policy = expandVpnSiteO365Policy(d.Get("o365_policy").([]any))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.VpnSitesCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnSiteRead(d, meta)
}
