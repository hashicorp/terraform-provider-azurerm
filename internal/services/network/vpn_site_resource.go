// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVpnSite() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVpnSiteCreate,
		Read:   resourceVpnSiteRead,
		Update: resourceVpnSiteUpdate,
		Delete: resourceVpnSiteDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVpnSiteID(id)
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
				ValidateFunc: validate.VpnSiteName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"virtual_wan_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualWANID,
			},

			"address_cidrs": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},

			"device_vendor": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"device_model": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"link": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"provider_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"speed_in_mbps": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},
						"ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"fqdn": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"bgp": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"asn": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"peering_address": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
								},
							},
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"o365_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"traffic_category": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allow_endpoint_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"default_endpoint_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"optimize_endpoint_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVpnSiteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVpnSiteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.VpnSitesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vpn_site", id.ID())
	}

	payload := virtualwans.VpnSite{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VpnSiteProperties{
			VirtualWAN: &virtualwans.SubResource{
				Id: utils.String(d.Get("virtual_wan_id").(string)),
			},
			DeviceProperties: expandVpnSiteDeviceProperties(d.Get("device_vendor").(string), d.Get("device_model").(string)),
			AddressSpace:     expandVpnSiteAddressSpace(d.Get("address_cidrs").(*pluginsdk.Set).List()),
			VpnSiteLinks:     expandVpnSiteLinks(d.Get("link").([]interface{})),
			O365Policy:       expandVpnSiteO365Policy(d.Get("o365_policy").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.VpnSitesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnSiteRead(d, meta)
}

func resourceVpnSiteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnSiteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnSitesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnSiteName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			deviceModel := ""
			deviceVendor := ""
			if deviceProp := props.DeviceProperties; deviceProp != nil {
				deviceModel = pointer.From(deviceProp.DeviceModel)
				deviceVendor = pointer.From(deviceProp.DeviceVendor)
			}
			d.Set("device_model", deviceModel)
			d.Set("device_vendor", deviceVendor)

			virtualWanId := ""
			if props.VirtualWAN != nil && props.VirtualWAN.Id != nil {
				parsed, err := virtualwans.ParseVirtualWANIDInsensitively(*props.VirtualWAN.Id)
				if err == nil {
					virtualWanId = parsed.ID()
				}
			}
			d.Set("virtual_wan_id", virtualWanId)

			if err := d.Set("address_cidrs", flattenVpnSiteAddressSpace(props.AddressSpace)); err != nil {
				return fmt.Errorf("setting `address_cidrs`: %+v", err)
			}
			if err := d.Set("link", flattenVpnSiteLinks(props.VpnSiteLinks)); err != nil {
				return fmt.Errorf("setting `link`: %+v", err)
			}
			if err := d.Set("o365_policy", flattenVpnSiteO365Policy(props.O365Policy)); err != nil {
				return fmt.Errorf("setting `o365_policy`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceVpnSiteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnSiteID(d.Id())
	if err != nil {
		return err
	}
	existing, err := client.VpnSitesGet(ctx, *id)
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

	if d.HasChange("address_cidrs") {
		payload.Properties.AddressSpace = expandVpnSiteAddressSpace(d.Get("address_cidrs").(*pluginsdk.Set).List())
	}

	if d.HasChange("device_vendor") || d.HasChange("device_model") {
		payload.Properties.DeviceProperties = expandVpnSiteDeviceProperties(d.Get("device_vendor").(string), d.Get("device_model").(string))
	}

	if d.HasChange("link") {
		payload.Properties.VpnSiteLinks = expandVpnSiteLinks(d.Get("link").([]interface{}))
	}

	if d.HasChange("o365_policy") {
		payload.Properties.O365Policy = expandVpnSiteO365Policy(d.Get("o365_policy").([]interface{}))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.VpnSitesCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVpnSiteRead(d, meta)
}

func resourceVpnSiteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnSiteID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VpnSitesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVpnSiteDeviceProperties(vendor, model string) *virtualwans.DeviceProperties {
	if vendor == "" && model == "" {
		return nil
	}
	output := &virtualwans.DeviceProperties{}
	if vendor != "" {
		output.DeviceVendor = &vendor
	}
	if model != "" {
		output.DeviceModel = &model
	}

	return output
}

func expandVpnSiteAddressSpace(input []interface{}) *virtualwans.AddressSpace {
	if len(input) == 0 {
		return nil
	}

	addressPrefixes := make([]string, 0)
	for _, addr := range input {
		addressPrefixes = append(addressPrefixes, addr.(string))
	}

	return &virtualwans.AddressSpace{
		AddressPrefixes: &addressPrefixes,
	}
}

func flattenVpnSiteAddressSpace(input *virtualwans.AddressSpace) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	return utils.FlattenStringSlice(input.AddressPrefixes)
}

func expandVpnSiteLinks(input []interface{}) *[]virtualwans.VpnSiteLink {
	if len(input) == 0 {
		return nil
	}

	result := make([]virtualwans.VpnSiteLink, 0)
	for _, e := range input {
		if e == nil {
			continue
		}
		e := e.(map[string]interface{})
		link := virtualwans.VpnSiteLink{
			Name: utils.String(e["name"].(string)),
			Properties: &virtualwans.VpnSiteLinkProperties{
				LinkProperties: &virtualwans.VpnLinkProviderProperties{
					LinkSpeedInMbps: pointer.To(int64(e["speed_in_mbps"].(int))),
				},
			},
		}

		if v, ok := e["provider_name"]; ok {
			link.Properties.LinkProperties.LinkProviderName = pointer.To(v.(string))
		}
		if v, ok := e["ip_address"]; ok {
			link.Properties.IPAddress = pointer.To(v.(string))
		}
		if v, ok := e["fqdn"]; ok {
			link.Properties.Fqdn = pointer.To(v.(string))
		}
		if v, ok := e["bgp"]; ok {
			link.Properties.BgpProperties = expandVpnSiteVpnLinkBgpSettings(v.([]interface{}))
		}

		result = append(result, link)
	}

	return &result
}

func flattenVpnSiteLinks(input *[]virtualwans.VpnSiteLink) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var name string
		if e.Name != nil {
			name = *e.Name
		}

		id := ""
		if e.Id != nil {
			id = *e.Id
		}

		var (
			ipAddress        string
			fqdn             string
			linkProviderName string
			linkSpeed        int
			bgpProperty      []interface{}
		)

		if prop := e.Properties; prop != nil {
			if prop.IPAddress != nil {
				ipAddress = *prop.IPAddress
			}

			if prop.Fqdn != nil {
				fqdn = *prop.Fqdn
			}

			if linkProp := prop.LinkProperties; linkProp != nil {
				if linkProp.LinkProviderName != nil {
					linkProviderName = *linkProp.LinkProviderName
				}
				if linkProp.LinkSpeedInMbps != nil {
					linkSpeed = int(*linkProp.LinkSpeedInMbps)
				}
			}

			bgpProperty = flattenVpnSiteVpnSiteBgpSettings(prop.BgpProperties)
		}

		link := map[string]interface{}{
			"name":          name,
			"id":            id,
			"provider_name": linkProviderName,
			"speed_in_mbps": linkSpeed,
			"ip_address":    ipAddress,
			"fqdn":          fqdn,
			"bgp":           bgpProperty,
		}

		output = append(output, link)
	}

	return output
}

func expandVpnSiteVpnLinkBgpSettings(input []interface{}) *virtualwans.VpnLinkBgpSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &virtualwans.VpnLinkBgpSettings{
		Asn:               utils.Int64(int64(v["asn"].(int))),
		BgpPeeringAddress: utils.String(v["peering_address"].(string)),
	}
}

func flattenVpnSiteVpnSiteBgpSettings(input *virtualwans.VpnLinkBgpSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var asn int
	if input.Asn != nil {
		asn = int(*input.Asn)
	}

	var peerAddress string
	if input.BgpPeeringAddress != nil {
		peerAddress = *input.BgpPeeringAddress
	}

	return []interface{}{
		map[string]interface{}{
			"asn":             asn,
			"peering_address": peerAddress,
		},
	}
}

func expandVpnSiteO365Policy(input []interface{}) *virtualwans.O365PolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	o365Policy := input[0].(map[string]interface{})

	return &virtualwans.O365PolicyProperties{
		BreakOutCategories: expandVpnSiteO365TrafficCategoryPolicy(o365Policy["traffic_category"].([]interface{})),
	}
}

func expandVpnSiteO365TrafficCategoryPolicy(input []interface{}) *virtualwans.O365BreakOutCategoryPolicies {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	trafficCategory := input[0].(map[string]interface{})

	return &virtualwans.O365BreakOutCategoryPolicies{
		Allow:    utils.Bool(trafficCategory["allow_endpoint_enabled"].(bool)),
		Default:  utils.Bool(trafficCategory["default_endpoint_enabled"].(bool)),
		Optimize: utils.Bool(trafficCategory["optimize_endpoint_enabled"].(bool)),
	}
}

func flattenVpnSiteO365Policy(input *virtualwans.O365PolicyProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	trafficCategory := make([]interface{}, 0)
	if input.BreakOutCategories != nil {
		trafficCategory = flattenVpnSiteO365TrafficCategoryPolicy(input.BreakOutCategories)
	}

	return []interface{}{
		map[string]interface{}{
			"traffic_category": trafficCategory,
		},
	}
}

func flattenVpnSiteO365TrafficCategoryPolicy(input *virtualwans.O365BreakOutCategoryPolicies) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	isAllowed := false
	if input.Allow != nil {
		isAllowed = *input.Allow
	}

	isDefault := false
	if input.Default != nil {
		isDefault = *input.Default
	}

	isOptimized := false
	if input.Optimize != nil {
		isOptimized = *input.Optimize
	}

	return []interface{}{
		map[string]interface{}{
			"allow_endpoint_enabled":    isAllowed,
			"default_endpoint_enabled":  isDefault,
			"optimize_endpoint_enabled": isOptimized,
		},
	}
}
