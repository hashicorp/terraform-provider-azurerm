// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/ddosprotectionplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipaddresses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePublicIp() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpCreate,
		Read:   resourcePublicIpRead,
		Update: resourcePublicIpUpdate,
		Delete: resourcePublicIpDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParsePublicIPAddressID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"allocation_method": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipaddresses.IPAllocationMethodStatic),
					string(publicipaddresses.IPAllocationMethodDynamic),
				}, false),
			},

			// Optional
			"ddos_protection_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipaddresses.DdosSettingsProtectionModeDisabled),
					string(publicipaddresses.DdosSettingsProtectionModeEnabled),
					string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited),
				}, false),
				Default: string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited),
			},

			"ddos_protection_plan_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
			},

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			"ip_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(publicipaddresses.IPVersionIPvFour),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipaddresses.IPVersionIPvFour),
					string(publicipaddresses.IPVersionIPvSix),
				}, false),
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default: func() interface{} {
					// https://azure.microsoft.com/en-us/updates/upgrade-to-standard-sku-public-ip-addresses-in-azure-by-30-september-2025-basic-sku-will-be-retired/
					if !features.FourPointOhBeta() {
						return string(publicipaddresses.PublicIPAddressSkuNameBasic)
					}
					return string(publicipaddresses.PublicIPAddressSkuNameStandard)
				}(),
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipaddresses.PublicIPAddressSkuNameBasic),
					string(publicipaddresses.PublicIPAddressSkuNameStandard),
				}, false),
			},

			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(publicipaddresses.PublicIPAddressSkuTierRegional),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipaddresses.PublicIPAddressSkuTierGlobal),
					string(publicipaddresses.PublicIPAddressSkuTierRegional),
				}, false),
			},

			"idle_timeout_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      4,
				ValidateFunc: validation.IntBetween(4, 30),
			},

			"domain_name_label": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.PublicIpDomainNameLabel,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"reverse_fqdn": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_ip_prefix_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
			},

			"ip_tags": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePublicIpCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP creation.")

	id := commonids.NewPublicIPAddressID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_public_ip", id.ID())
	}

	sku := d.Get("sku").(string)
	ipAllocationMethod := d.Get("allocation_method").(string)

	if strings.EqualFold(sku, "standard") {
		if !strings.EqualFold(ipAllocationMethod, "static") {
			return fmt.Errorf("static IP allocation must be used when creating Standard SKU public IP addresses")
		}
	}

	ddosProtectionMode := d.Get("ddos_protection_mode").(string)

	publicIp := publicipaddresses.PublicIPAddress{
		Name:             pointer.To(id.PublicIPAddressesName),
		ExtendedLocation: expandEdgeZoneNew(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Sku: &publicipaddresses.PublicIPAddressSku{
			Name: pointer.To(publicipaddresses.PublicIPAddressSkuName(sku)),
			Tier: pointer.To(publicipaddresses.PublicIPAddressSkuTier(d.Get("sku_tier").(string))),
		},
		Properties: &publicipaddresses.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: pointer.To(publicipaddresses.IPAllocationMethod(ipAllocationMethod)),
			PublicIPAddressVersion:   pointer.To(publicipaddresses.IPVersion(d.Get("ip_version").(string))),
			IdleTimeoutInMinutes:     pointer.To(int64(d.Get("idle_timeout_in_minutes").(int))),
			DdosSettings: &publicipaddresses.DdosSettings{
				ProtectionMode: pointer.To(publicipaddresses.DdosSettingsProtectionMode(ddosProtectionMode)),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	if dnlOk || rfqdnOk {
		dnsSettings := publicipaddresses.PublicIPAddressDnsSettings{}

		if rfqdnOk {
			dnsSettings.ReverseFqdn = pointer.To(rfqdn.(string))
		}

		if dnlOk {
			dnsSettings.DomainNameLabel = pointer.To(dnl.(string))
		}

		publicIp.Properties.DnsSettings = &dnsSettings
	}

	if v, ok := d.GetOk("ip_tags"); ok {
		ipTags := v.(map[string]interface{})
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

	if err := client.CreateOrUpdateThenPoll(ctx, id, publicIp); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePublicIpRead(d, meta)
}

func resourcePublicIpUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP update.")

	id, err := commonids.ParsePublicIPAddressID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("allocation_method") {
		payload.Properties.PublicIPAllocationMethod = pointer.To(publicipaddresses.IPAllocationMethod(d.Get("allocation_method").(string)))
	}

	if d.HasChange("ddos_protection_mode") {
		if payload.Properties.DdosSettings == nil {
			payload.Properties.DdosSettings = &publicipaddresses.DdosSettings{}
		}
		payload.Properties.DdosSettings.ProtectionMode = pointer.To(publicipaddresses.DdosSettingsProtectionMode(d.Get("ddos_protection_mode").(string)))
	}

	if d.HasChange("ddos_protection_plan_id") {
		if !strings.EqualFold(string(*payload.Properties.DdosSettings.ProtectionMode), "enabled") {
			return fmt.Errorf("ddos protection plan id can only be set when ddos protection is enabled")
		}
		if payload.Properties.DdosSettings == nil {
			payload.Properties.DdosSettings = &publicipaddresses.DdosSettings{}
		}
		payload.Properties.DdosSettings.DdosProtectionPlan = &publicipaddresses.SubResource{
			Id: pointer.To(d.Get("ddos_protection_plan_id").(string)),
		}
	}

	if d.HasChange("idle_timeout_in_minutes") {
		payload.Properties.IdleTimeoutInMinutes = utils.Int64(int64(d.Get("idle_timeout_in_minutes").(int)))
	}

	if d.HasChange("domain_name_label") {
		if payload.Properties.DnsSettings == nil {
			payload.Properties.DnsSettings = &publicipaddresses.PublicIPAddressDnsSettings{}
		}
		payload.Properties.DnsSettings.DomainNameLabel = utils.String(d.Get("domain_name_label").(string))
	}

	if d.HasChange("reverse_fqdn") {
		if payload.Properties.DnsSettings == nil {
			payload.Properties.DnsSettings = &publicipaddresses.PublicIPAddressDnsSettings{}
		}
		payload.Properties.DnsSettings.ReverseFqdn = utils.String(d.Get("reverse_fqdn").(string))
	}

	if d.HasChanges("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourcePublicIpRead(d, meta)
}

func resourcePublicIpRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParsePublicIPAddressID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.PublicIPAddressesName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneNew(model.ExtendedLocation))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(pointer.From(sku.Name)))
			d.Set("sku_tier", string(pointer.From(sku.Tier)))
		}
		if props := model.Properties; props != nil {
			d.Set("allocation_method", string(pointer.From(props.PublicIPAllocationMethod)))
			d.Set("ip_version", string(pointer.From(props.PublicIPAddressVersion)))

			if publicIpPrefix := props.PublicIPPrefix; publicIpPrefix != nil {
				d.Set("public_ip_prefix_id", publicIpPrefix.Id)
			}

			if settings := props.DnsSettings; settings != nil {
				d.Set("fqdn", settings.Fqdn)
				d.Set("reverse_fqdn", settings.ReverseFqdn)
				d.Set("domain_name_label", settings.DomainNameLabel)
			}

			ddosProtectionMode := string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited)
			if ddosSetting := props.DdosSettings; ddosSetting != nil {
				ddosProtectionMode = string(pointer.From(ddosSetting.ProtectionMode))
				if subResource := ddosSetting.DdosProtectionPlan; subResource != nil {
					d.Set("ddos_protection_plan_id", subResource.Id)
				}
			}
			d.Set("ddos_protection_mode", ddosProtectionMode)

			d.Set("ip_tags", flattenPublicIpPropsIpTags(props.IPTags))

			d.Set("ip_address", props.IPAddress)
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourcePublicIpDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParsePublicIPAddressID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenPublicIpPropsIpTags(input *[]publicipaddresses.IPTag) map[string]interface{} {
	out := make(map[string]interface{})

	if input != nil {
		for _, tag := range *input {
			if tag.IPTagType != nil {
				out[*tag.IPTagType] = tag.Tag
			}
		}
	}

	return out
}
