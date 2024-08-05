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
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var VPNGatewayResourceName = "azurerm_vpn_gateway"

func resourceVPNGateway() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceVPNGatewayCreate,
		Read:   resourceVPNGatewayRead,
		Update: resourceVPNGatewayUpdate,
		Delete: resourceVPNGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVpnGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"routing_preference": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Microsoft Network",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Microsoft Network",
					"Internet",
				}, false),
			},

			"bgp_route_translation_for_nat_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"bgp_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"asn": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"peer_weight": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"bgp_peering_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"instance_0_bgp_peering_address": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_ips": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: commonValidate.IPv4Address,
										},
									},

									"ip_configuration_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"default_ips": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"tunnel_ips": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},

						"instance_1_bgp_peering_address": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"custom_ips": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: commonValidate.IPv4Address,
										},
									},

									"ip_configuration_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"default_ips": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"tunnel_ips": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"scale_unit": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["routing_preference"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Microsoft Network",
				"Internet",
			}, false),
		}
	}

	return resource
}

func resourceVPNGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVpnGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.VpnGatewaysGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vpn_gateway", id.ID())
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]interface{})
	bgpSettings := expandVPNGatewayBGPSettings(bgpSettingsRaw)
	payload := virtualwans.VpnGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VpnGatewayProperties{
			EnableBgpRouteTranslationForNat: pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool)),
			BgpSettings:                     bgpSettings,
			VirtualHub: &virtualwans.SubResource{
				Id: utils.String(d.Get("virtual_hub_id").(string)),
			},
			VpnGatewayScaleUnit:         pointer.To(int64(d.Get("scale_unit").(int))),
			IsRoutingPreferenceInternet: pointer.To(d.Get("routing_preference").(string) == "Internet"),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// `vpnGatewayParameters.Properties.bgpSettings.bgpPeeringAddress` customer cannot provide this field during create. This will be set with default value once gateway is created.
	// it could only be updated
	if len(bgpSettingsRaw) > 0 {
		resp, err := client.VpnGatewaysGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
		}
		props := resp.Model.Properties

		if props.BgpSettings != nil && props.BgpSettings.BgpPeeringAddresses != nil {

			val := bgpSettingsRaw[0].(map[string]interface{})
			input0 := val["instance_0_bgp_peering_address"].([]interface{})
			input1 := val["instance_1_bgp_peering_address"].([]interface{})

			if len(input0) > 0 || len(input1) > 0 {
				if len(input0) > 0 && input0[0] != nil {
					val := input0[0].(map[string]interface{})
					(*props.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
				}
				if len(input1) > 0 && input1[0] != nil {
					val := input1[0].(map[string]interface{})
					(*props.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
				}

				resp.Model.Properties = props

				if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, id, *resp.Model); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}
		}
	}

	return resourceVPNGatewayRead(d, meta)
}

func resourceVPNGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(id.VpnGatewayName, VPNGatewayResourceName)

	existing, err := client.VpnGatewaysGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	model := *existing.Model
	if d.HasChange("scale_unit") {
		model.Properties.VpnGatewayScaleUnit = pointer.To(int64(d.Get("scale_unit").(int)))
	}
	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}
	if d.HasChange("bgp_route_translation_for_nat_enabled") {
		model.Properties.EnableBgpRouteTranslationForNat = utils.Bool(d.Get("bgp_route_translation_for_nat_enabled").(bool))
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]interface{})
	if len(bgpSettingsRaw) > 0 {
		val := bgpSettingsRaw[0].(map[string]interface{})

		if d.HasChange("bgp_settings.0.instance_0_bgp_peering_address") {
			if input := val["instance_0_bgp_peering_address"].([]interface{}); len(input) > 0 {
				val := input[0].(map[string]interface{})
				(*model.Properties.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
			}
		}
		if d.HasChange("bgp_settings.0.instance_1_bgp_peering_address") {
			if input := val["instance_1_bgp_peering_address"].([]interface{}); len(input) > 0 {
				val := input[0].(map[string]interface{})
				(*model.Properties.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*pluginsdk.Set).List())
			}
		}
	}

	if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVPNGatewayRead(d, meta)
}

func resourceVPNGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VpnGatewaysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VpnGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("bgp_settings", flattenVPNGatewayBGPSettings(props.BgpSettings)); err != nil {
				return fmt.Errorf("setting `bgp_settings`: %+v", err)
			}

			bgpRouteTranslationForNatEnabled := false
			if props.EnableBgpRouteTranslationForNat != nil {
				bgpRouteTranslationForNatEnabled = *props.EnableBgpRouteTranslationForNat
			}
			d.Set("bgp_route_translation_for_nat_enabled", bgpRouteTranslationForNatEnabled)

			scaleUnit := 0
			if props.VpnGatewayScaleUnit != nil {
				scaleUnit = int(*props.VpnGatewayScaleUnit)
			}
			d.Set("scale_unit", scaleUnit)

			virtualHubId := ""
			if props.VirtualHub != nil && props.VirtualHub.Id != nil {
				virtualHubId = *props.VirtualHub.Id
			}
			d.Set("virtual_hub_id", virtualHubId)

			isRoutingPreferenceInternet := "Microsoft Network"
			if props.IsRoutingPreferenceInternet != nil && *props.IsRoutingPreferenceInternet {
				isRoutingPreferenceInternet = "Internet"
			}
			d.Set("routing_preference", isRoutingPreferenceInternet)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceVPNGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVpnGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VpnGatewaysDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVPNGatewayBGPSettings(input []interface{}) *virtualwans.BgpSettings {
	if len(input) == 0 {
		return nil
	}

	val := input[0].(map[string]interface{})
	return &virtualwans.BgpSettings{
		Asn:        pointer.To(int64(val["asn"].(int))),
		PeerWeight: pointer.To(int64(val["peer_weight"].(int))),
	}
}

func flattenVPNGatewayBGPSettings(input *virtualwans.BgpSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	asn := 0
	if input.Asn != nil {
		asn = int(*input.Asn)
	}

	bgpPeeringAddress := ""
	if input.BgpPeeringAddress != nil {
		bgpPeeringAddress = *input.BgpPeeringAddress
	}

	peerWeight := 0
	if input.PeerWeight != nil {
		peerWeight = int(*input.PeerWeight)
	}

	var instance0BgpPeeringAddress, instance1BgpPeeringAddress []interface{}
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 0 {
		instance0BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[0])
	}
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 1 {
		instance1BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[1])
	}

	return []interface{}{
		map[string]interface{}{
			"asn":                            asn,
			"bgp_peering_address":            bgpPeeringAddress,
			"instance_0_bgp_peering_address": instance0BgpPeeringAddress,
			"instance_1_bgp_peering_address": instance1BgpPeeringAddress,
			"peer_weight":                    peerWeight,
		},
	}
}

func flattenVPNGatewayIPConfigurationBgpPeeringAddress(input virtualwans.IPConfigurationBgpPeeringAddress) []interface{} {
	ipConfigurationID := ""
	if input.IPconfigurationId != nil {
		ipConfigurationID = *input.IPconfigurationId
	}

	return []interface{}{
		map[string]interface{}{
			"ip_configuration_id": ipConfigurationID,
			"custom_ips":          utils.FlattenStringSlice(input.CustomBgpIPAddresses),
			"default_ips":         utils.FlattenStringSlice(input.DefaultBgpIPAddresses),
			"tunnel_ips":          utils.FlattenStringSlice(input.TunnelIPAddresses),
		},
	}
}
