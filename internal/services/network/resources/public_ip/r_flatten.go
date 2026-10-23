// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourcePublicIpFlatten(d *pluginsdk.ResourceData, id *commonids.PublicIPAddressId, model *publicipaddresses.PublicIPAddress) error {
	d.Set("name", id.PublicIPAddressesName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneNew(model.ExtendedLocation))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if sku := model.Sku; sku != nil {
			d.Set("sku", pointer.FromEnum(sku.Name))
			d.Set("sku_tier", pointer.FromEnum(sku.Tier))
		}
		if props := model.Properties; props != nil {
			d.Set("allocation_method", pointer.FromEnum(props.PublicIPAllocationMethod))
			d.Set("ip_version", pointer.FromEnum(props.PublicIPAddressVersion))

			if publicIpPrefix := props.PublicIPPrefix; publicIpPrefix != nil {
				d.Set("public_ip_prefix_id", publicIpPrefix.Id)
			}

			fqdn := ""
			reverseFqdn := ""
			domainNameLabel := ""
			domainNameLabelScope := ""
			if settings := props.DnsSettings; settings != nil {
				fqdn = pointer.From(settings.Fqdn)
				reverseFqdn = pointer.From(settings.ReverseFqdn)
				domainNameLabel = pointer.From(settings.DomainNameLabel)
				domainNameLabelScope = pointer.FromEnum(settings.DomainNameLabelScope)
			}

			d.Set("fqdn", fqdn)
			d.Set("reverse_fqdn", reverseFqdn)
			d.Set("domain_name_label", domainNameLabel)
			d.Set("domain_name_label_scope", domainNameLabelScope)

			ddosProtectionMode := string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited)
			if ddosSetting := props.DdosSettings; ddosSetting != nil {
				ddosProtectionMode = pointer.FromEnum(ddosSetting.ProtectionMode)
				if subResource := ddosSetting.DdosProtectionPlan; subResource != nil {
					d.Set("ddos_protection_plan_id", subResource.Id)
				}
			}
			d.Set("ddos_protection_mode", ddosProtectionMode)

			d.Set("ip_tags", flattenPublicIpPropsIpTags(props.IPTags))

			d.Set("ip_address", props.IPAddress)
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenPublicIpPropsIpTags(input *[]publicipaddresses.IPTag) map[string]any {
	out := make(map[string]any)

	if input != nil {
		for _, tag := range *input {
			if tag.IPTagType != nil {
				out[*tag.IPTagType] = tag.Tag
			}
		}
	}

	return out
}
