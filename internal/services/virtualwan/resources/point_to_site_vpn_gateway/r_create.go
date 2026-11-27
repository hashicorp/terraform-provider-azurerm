// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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

func resourcePointToSiteVPNGatewayCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := commonids.NewVirtualWANP2SVPNGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.P2sVpnGatewaysGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_point_to_site_vpn_gateway", id.ID())
		}
	}

	parameters := virtualwans.P2SVpnGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.P2SVpnGatewayProperties{
			IsRoutingPreferenceInternet: pointer.To(d.Get("routing_preference_internet_enabled").(bool)),
			P2SConnectionConfigurations: expandPointToSiteVPNGatewayConnectionConfiguration(d.Get("connection_configuration").([]any)),
			VpnServerConfiguration: &virtualwans.SubResource{
				Id: pointer.To(d.Get("vpn_server_configuration_id").(string)),
			},
			VirtualHub: &virtualwans.SubResource{
				Id: pointer.To(d.Get("virtual_hub_id").(string)),
			},
			VpnGatewayScaleUnit: pointer.To(int64(d.Get("scale_unit").(int))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}
	customDNSServers := pluginsdk.ExpandStringSlice(d.Get("dns_servers").([]any))
	if len(*customDNSServers) != 0 {
		parameters.Properties.CustomDnsServers = customDNSServers
	}

	if err := client.P2sVpnGatewaysCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePointToSiteVPNGatewayRead(d, meta)
}
