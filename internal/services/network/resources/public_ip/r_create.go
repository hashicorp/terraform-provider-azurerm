// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewPublicIPAddressID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, publicipaddresses.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_public_ip", id.ID())
		}
	}

	sku := d.Get("sku").(string)
	ipAllocationMethod := d.Get("allocation_method").(string)

	if strings.EqualFold(sku, string(publicipaddresses.PublicIPAddressSkuNameStandard)) || strings.EqualFold(sku, string(publicipaddresses.PublicIPAddressSkuNameStandardVTwo)) {
		if !strings.EqualFold(ipAllocationMethod, string(publicipaddresses.IPAllocationMethodStatic)) {
			return fmt.Errorf("`allocation_method` must be set to `Static` when `sku` is set to `Standard` or `StandardV2`")
		}
	}

	ddosProtectionMode := d.Get("ddos_protection_mode").(string)

	publicIp := publicipaddresses.PublicIPAddress{
		Name:             pointer.To(id.PublicIPAddressesName),
		ExtendedLocation: expandEdgeZoneNew(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Sku: &publicipaddresses.PublicIPAddressSku{
			Name: pointer.ToEnum[publicipaddresses.PublicIPAddressSkuName](sku),
			Tier: pointer.ToEnum[publicipaddresses.PublicIPAddressSkuTier](d.Get("sku_tier").(string)),
		},
		Properties: &publicipaddresses.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: pointer.ToEnum[publicipaddresses.IPAllocationMethod](ipAllocationMethod),
			PublicIPAddressVersion:   pointer.ToEnum[publicipaddresses.IPVersion](d.Get("ip_version").(string)),
			IdleTimeoutInMinutes:     pointer.To(int64(d.Get("idle_timeout_in_minutes").(int))),
			DdosSettings: &publicipaddresses.DdosSettings{
				ProtectionMode: pointer.ToEnum[publicipaddresses.DdosSettingsProtectionMode](ddosProtectionMode),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	ddosProtectionPlanId, planOk := d.GetOk("ddos_protection_plan_id")
	if planOk {
		if !strings.EqualFold(ddosProtectionMode, "enabled") {
			return fmt.Errorf("ddos protection plan id can only be set when ddos protection is enabled")
		}
		publicIp.Properties.DdosSettings.DdosProtectionPlan = &publicipaddresses.SubResource{
			Id: pointer.To(ddosProtectionPlanId.(string)),
		}
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		publicIp.Zones = &zones
	}

	publicIpPrefixId, publicIpPrefixIdOk := d.GetOk("public_ip_prefix_id")

	if publicIpPrefixIdOk {
		publicIpPrefix := publicipaddresses.SubResource{}
		publicIpPrefix.Id = pointer.To(publicIpPrefixId.(string))
		publicIp.Properties.PublicIPPrefix = &publicIpPrefix
	}

	dnl, dnlOk := d.GetOk("domain_name_label")
	rfqdn, rfqdnOk := d.GetOk("reverse_fqdn")
	dnlc, dnlcOk := d.GetOk("domain_name_label_scope")

	if dnlOk || rfqdnOk || dnlcOk {
		dnsSettings := publicipaddresses.PublicIPAddressDnsSettings{}

		if rfqdnOk {
			dnsSettings.ReverseFqdn = pointer.To(rfqdn.(string))
		}

		if dnlOk {
			dnsSettings.DomainNameLabel = pointer.To(dnl.(string))
		}

		if dnlcOk {
			dnsSettings.DomainNameLabelScope = pointer.ToEnum[publicipaddresses.PublicIPAddressDnsSettingsDomainNameLabelScope](dnlc.(string))
		}

		publicIp.Properties.DnsSettings = &dnsSettings
	}

	if v, ok := d.GetOk("ip_tags"); ok {
		ipTags := v.(map[string]any)
		newIpTags := []publicipaddresses.IPTag{}

		for key, val := range ipTags {
			ipTag := publicipaddresses.IPTag{
				IPTagType: pointer.To(key),
				Tag:       pointer.To(val.(string)),
			}
			newIpTags = append(newIpTags, ipTag)
		}

		publicIp.Properties.IPTags = &newIpTags
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, publicIp, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourcePublicIpRead(d, meta)
}
