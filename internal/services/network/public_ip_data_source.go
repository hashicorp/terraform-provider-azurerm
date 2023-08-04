// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourcePublicIP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePublicIPRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"allocation_method": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ddos_protection_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ddos_protection_plan_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"domain_name_label": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"reverse_fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_tags": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"zones": commonschema.ZonesMultipleComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePublicIPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPublicIpAddressID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	if resp.PublicIPAddressPropertiesFormat == nil {
		return fmt.Errorf("retreving %s: `properties` was nil", id)
	}

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku", skuName)

	props := *resp.PublicIPAddressPropertiesFormat

	domainNameLabel := ""
	fqdn := ""
	reverseFqdn := ""
	if dnsSettings := props.DNSSettings; dnsSettings != nil {
		if dnsSettings.DomainNameLabel != nil {
			domainNameLabel = *dnsSettings.DomainNameLabel
		}
		if dnsSettings.Fqdn != nil {
			fqdn = *dnsSettings.Fqdn
		}
		if dnsSettings.ReverseFqdn != nil {
			reverseFqdn = *dnsSettings.ReverseFqdn
		}
	}

	if ddosSetting := props.DdosSettings; ddosSetting != nil {
		d.Set("ddos_protection_mode", string(ddosSetting.ProtectionMode))
		if subResource := ddosSetting.DdosProtectionPlan; subResource != nil {
			d.Set("ddos_protection_plan_id", subResource.ID)
		}
	}

	d.Set("domain_name_label", domainNameLabel)
	d.Set("fqdn", fqdn)
	d.Set("reverse_fqdn", reverseFqdn)

	d.Set("allocation_method", string(props.PublicIPAllocationMethod))
	d.Set("ip_address", props.IPAddress)
	d.Set("ip_version", string(props.PublicIPAddressVersion))
	d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)

	d.Set("ip_tags", flattenPublicIpPropsIpTags(props.IPTags))

	return tags.FlattenAndSet(d, resp.Tags)
}
