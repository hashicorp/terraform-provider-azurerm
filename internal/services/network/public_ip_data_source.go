// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePublicIP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Public IP Address.",
		Read:        dataSourcePublicIPRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the public IP address.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"domain_name_label_scope": {
				Description: "Scope for the domain name label.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"edge_zone": commonschema.EdgeZoneComputed(),

			"public_ip_prefix_id": {
				Description: "If specified then public IP address allocated will be provided from the public IP prefix resource.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"sku": {
				Description: "The SKU of the Public IP.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"sku_tier": {
				Description: "The SKU Tier that is used for the Public IP.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"allocation_method": {
				Description: "The allocation method for this IP address.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ddos_protection_mode": {
				Description: "The DDoS protection mode of the public IP.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ddos_protection_plan_id": {
				Description: "The ID of DDoS protection plan associated with the public IP.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ip_version": {
				Description: "The IP version being used, for example `IPv4` or `IPv6`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"domain_name_label": {
				Description: "The label for the Domain Name.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"idle_timeout_in_minutes": {
				Description: "The timeout for the TCP idle connection.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"fqdn": {
				Description: "Fully qualified domain name of the A DNS record associated with the public IP. This is the concatenation of the domainNameLabel and the regionalized DNS zone.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"reverse_fqdn": {
				Description: "The fully qualified domain name that resolves to this public IP address.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ip_address": {
				Description: "The IP address value that was allocated.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"ip_tags": {
				Description: "A mapping of tags to assigned to the resource.",
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"zones": commonschema.ZonesMultipleComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePublicIPRead(d *pluginsdk.ResourceData, meta any) error {
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
		d.Set("edge_zone", flattenEdgeZoneNew(model.ExtendedLocation))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = pointer.FromEnum(sku.Name)
			d.Set("sku_tier", pointer.FromEnum(sku.Tier))
		}
		d.Set("sku", skuName)

		if props := model.Properties; props != nil {
			domainNameLabel := ""
			fqdn := ""
			reverseFqdn := ""
			domainNameLabelScope := ""
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
				domainNameLabelScope = pointer.FromEnum(dnsSettings.DomainNameLabelScope)
			}

			if ddosSetting := props.DdosSettings; ddosSetting != nil {
				d.Set("ddos_protection_mode", pointer.FromEnum(ddosSetting.ProtectionMode))
				if subResource := ddosSetting.DdosProtectionPlan; subResource != nil {
					d.Set("ddos_protection_plan_id", subResource.Id)
				}
			}

			if publicIpPrefix := props.PublicIPPrefix; publicIpPrefix != nil {
				d.Set("public_ip_prefix_id", publicIpPrefix.Id)
			}

			d.Set("domain_name_label", domainNameLabel)
			d.Set("domain_name_label_scope", domainNameLabelScope)
			d.Set("fqdn", fqdn)
			d.Set("reverse_fqdn", reverseFqdn)

			d.Set("allocation_method", pointer.FromEnum(props.PublicIPAllocationMethod))
			d.Set("ip_address", props.IPAddress)
			d.Set("ip_version", pointer.FromEnum(props.PublicIPAddressVersion))
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)

			d.Set("ip_tags", flattenPublicIpPropsIpTags(props.IPTags))
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
