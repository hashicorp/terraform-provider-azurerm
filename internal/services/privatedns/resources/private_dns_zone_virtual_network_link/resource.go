// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone_virtual_network_link

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name private_dns_zone_virtual_network_link -properties "name" -compare-values "subscription_id:private_dns_zone_id,resource_group_name:private_dns_zone_id,private_dns_zone_name:private_dns_zone_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Read:     resourcePrivateDnsZoneVirtualNetworkLinkRead,
		Update:   resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate,
		Delete:   resourcePrivateDnsZoneVirtualNetworkLinkDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&virtualnetworklinks.VirtualNetworkLinkId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualnetworklinks.VirtualNetworkLinkId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: privateDnsZoneVirtualNetworkLinkSchema(),
	}
}
