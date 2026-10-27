// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone_virtual_network_link

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePrivateDnsZoneVirtualNetworkLinkFlatten(d *pluginsdk.ResourceData, id *virtualnetworklinks.VirtualNetworkLinkId, model *virtualnetworklinks.VirtualNetworkLink) error {
	d.Set("name", id.VirtualNetworkLinkName)
	d.Set("private_dns_zone_id", privatezones.NewPrivateDnsZoneID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName).ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("registration_enabled", props.RegistrationEnabled)
			d.Set("resolution_policy", pointer.From(props.ResolutionPolicy))

			if network := props.VirtualNetwork; network != nil {
				d.Set("virtual_network_id", network.Id)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
