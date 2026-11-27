// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_site

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnSiteRead(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseVpnSiteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnSitesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnSiteName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			deviceModel := ""
			deviceVendor := ""
			if deviceProp := props.DeviceProperties; deviceProp != nil {
				deviceModel = pointer.From(deviceProp.DeviceModel)
				deviceVendor = pointer.From(deviceProp.DeviceVendor)
			}
			d.Set("device_model", deviceModel)
			d.Set("device_vendor", deviceVendor)

			virtualWanId := ""
			if props.VirtualWAN != nil && props.VirtualWAN.Id != nil {
				parsed, err := virtualwans.ParseVirtualWANIDInsensitively(*props.VirtualWAN.Id)
				if err == nil {
					virtualWanId = parsed.ID()
				}
			}
			d.Set("virtual_wan_id", virtualWanId)

			if err := d.Set("address_cidrs", flattenVpnSiteAddressSpace(props.AddressSpace)); err != nil {
				return fmt.Errorf("setting `address_cidrs`: %+v", err)
			}
			if err := d.Set("link", flattenVpnSiteLinks(props.VpnSiteLinks)); err != nil {
				return fmt.Errorf("setting `link`: %+v", err)
			}
			if err := d.Set("o365_policy", flattenVpnSiteO365Policy(props.O365Policy)); err != nil {
				return fmt.Errorf("setting `o365_policy`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}

	return nil
}
