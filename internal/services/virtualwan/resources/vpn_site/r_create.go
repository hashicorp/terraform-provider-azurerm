// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_site

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVpnSiteCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := virtualwans.NewVpnSiteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.VpnSitesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_vpn_site", id.ID())
		}
	}

	payload := virtualwans.VpnSite{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VpnSiteProperties{
			VirtualWAN: &virtualwans.SubResource{
				Id: pointer.To(d.Get("virtual_wan_id").(string)),
			},
			DeviceProperties: expandVpnSiteDeviceProperties(d.Get("device_vendor").(string), d.Get("device_model").(string)),
			AddressSpace:     expandVpnSiteAddressSpace(d.Get("address_cidrs").(*pluginsdk.Set).List()),
			VpnSiteLinks:     expandVpnSiteLinks(d.Get("link").([]any)),
			O365Policy:       expandVpnSiteO365Policy(d.Get("o365_policy").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.VpnSitesCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnSiteRead(d, meta)
}
