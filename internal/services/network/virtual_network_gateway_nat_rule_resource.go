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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayNatRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayNatRuleCreate,
		Read:   resourceVirtualNetworkGatewayNatRuleRead,
		Update: resourceVirtualNetworkGatewayNatRuleUpdate,
		Delete: resourceVirtualNetworkGatewayNatRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworkgateways.ParseVirtualNetworkGatewayNatRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
			},

			"external_mapping": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"address_space": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},

						"port_range": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"internal_mapping": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"address_space": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},

						"port_range": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(virtualnetworkgateways.VpnNatRuleModeEgressSnat),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualnetworkgateways.VpnNatRuleModeEgressSnat),
					string(virtualnetworkgateways.VpnNatRuleModeIngressSnat),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(virtualnetworkgateways.VpnNatRuleTypeStatic),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualnetworkgateways.VpnNatRuleTypeStatic),
					string(virtualnetworkgateways.VpnNatRuleTypeDynamic),
				}, false),
			},

			"ip_configuration_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.VirtualNetworkGatewayIpConfigurationID,
			},
		},
	}
}

func resourceVirtualNetworkGatewayNatRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetGatewayId, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Get("virtual_network_gateway_id").(string))
	if err != nil {
		return err
	}

	id := virtualnetworkgateways.NewVirtualNetworkGatewayNatRuleID(subscriptionId, d.Get("resource_group_name").(string), vnetGatewayId.VirtualNetworkGatewayName, d.Get("name").(string))

	existing, err := client.VirtualNetworkGatewayNatRulesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_network_gateway_nat_rule", id.ID())
	}

	props := virtualnetworkgateways.VirtualNetworkGatewayNatRule{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualnetworkgateways.VirtualNetworkGatewayNatRuleProperties{
			ExternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("external_mapping").([]interface{})),
			InternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("internal_mapping").([]interface{})),
			Mode:             pointer.To(virtualnetworkgateways.VpnNatRuleMode(d.Get("mode").(string))),
			Type:             pointer.To(virtualnetworkgateways.VpnNatRuleType(d.Get("type").(string))),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.Properties.IPConfigurationId = pointer.To(v.(string))
	}

	if err := client.VirtualNetworkGatewayNatRulesCreateOrUpdateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}

func resourceVirtualNetworkGatewayNatRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualNetworkGatewayNatRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NatRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	vnetGatewayId := virtualnetworkgateways.NewVirtualNetworkGatewayID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName)
	d.Set("virtual_network_gateway_id", vnetGatewayId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("external_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.ExternalMappings)); err != nil {
				return fmt.Errorf("setting `external_mapping`: %+v", err)
			}

			if err := d.Set("internal_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.InternalMappings)); err != nil {
				return fmt.Errorf("setting `internal_mapping`: %+v", err)
			}

			d.Set("ip_configuration_id", props.IPConfigurationId)
			d.Set("mode", string(pointer.From(props.Mode)))
			d.Set("type", string(pointer.From(props.Type)))
		}
	}
	return nil
}

func resourceVirtualNetworkGatewayNatRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	props := virtualnetworkgateways.VirtualNetworkGatewayNatRule{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualnetworkgateways.VirtualNetworkGatewayNatRuleProperties{
			ExternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("external_mapping").([]interface{})),
			InternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("internal_mapping").([]interface{})),
			Mode:             pointer.To(virtualnetworkgateways.VpnNatRuleMode(d.Get("mode").(string))),
			Type:             pointer.To(virtualnetworkgateways.VpnNatRuleType(d.Get("type").(string))),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.Properties.IPConfigurationId = pointer.To(v.(string))
	}

	if err := client.VirtualNetworkGatewayNatRulesCreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}

func resourceVirtualNetworkGatewayNatRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VirtualNetworkGatewayNatRulesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandVirtualNetworkGatewayNatRuleMappings(input []interface{}) *[]virtualnetworkgateways.VpnNatRuleMapping {
	results := make([]virtualnetworkgateways.VpnNatRuleMapping, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := virtualnetworkgateways.VpnNatRuleMapping{
			AddressSpace: pointer.To(v["address_space"].(string)),
		}

		if portRange := v["port_range"].(string); portRange != "" {
			result.PortRange = pointer.To(portRange)
		}

		results = append(results, result)
	}

	return &results
}

func flattenVirtualNetworkGatewayNatRuleMappings(input *[]virtualnetworkgateways.VpnNatRuleMapping) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var addressSpace string
		if item.AddressSpace != nil {
			addressSpace = *item.AddressSpace
		}

		var portRange string
		if item.PortRange != nil {
			portRange = *item.PortRange
		}

		results = append(results, map[string]interface{}{
			"address_space": addressSpace,
			"port_range":    portRange,
		})
	}

	return results
}
