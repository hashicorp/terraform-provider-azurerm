// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	client := meta.(*clients.Client).Network.PublicIPAddresses
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewPublicIPAddressID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = string(pointer.From(sku.Name))
		}
		d.Set("sku", skuName)

		if props := model.Properties; props != nil {
			domainNameLabel := ""
			fqdn := ""
			reverseFqdn := ""
			if dnsSettings := props.DnsSettings; dnsSettings != nil {
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
				d.Set("ddos_protection_mode", string(pointer.From(ddosSetting.ProtectionMode)))
				if subResource := ddosSetting.DdosProtectionPlan; subResource != nil {
					d.Set("ddos_protection_plan_id", subResource.Id)
				}
			}

			d.Set("domain_name_label", domainNameLabel)
			d.Set("fqdn", fqdn)
			d.Set("reverse_fqdn", reverseFqdn)

			d.Set("allocation_method", string(pointer.From(props.PublicIPAllocationMethod)))
			d.Set("ip_address", props.IPAddress)
			d.Set("ip_version", string(pointer.From(props.PublicIPAddressVersion)))
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)

			d.Set("ip_tags", flattenPublicIpPropsIpTags(props.IPTags))
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
