// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
			_, err := parse.VirtualNetworkGatewayNatRuleID(id)
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
				ValidateFunc: validate.VirtualNetworkGatewayID,
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
				Default:  string(network.VpnNatRuleModeEgressSnat),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VpnNatRuleModeEgressSnat),
					string(network.VpnNatRuleModeIngressSnat),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.VpnNatRuleTypeStatic),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VpnNatRuleTypeStatic),
					string(network.VpnNatRuleTypeDynamic),
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
	client := meta.(*clients.Client).Network.VnetGatewayNatRuleClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetGatewayId, err := parse.VirtualNetworkGatewayID(d.Get("virtual_network_gateway_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewVirtualNetworkGatewayNatRuleID(subscriptionId, d.Get("resource_group_name").(string), vnetGatewayId.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_virtual_network_gateway_nat_rule", id.ID())
	}

	props := network.VirtualNetworkGatewayNatRule{
		Name: utils.String(d.Get("name").(string)),
		VirtualNetworkGatewayNatRuleProperties: &network.VirtualNetworkGatewayNatRuleProperties{
			ExternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("external_mapping").([]interface{})),
			InternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("internal_mapping").([]interface{})),
			Mode:             network.VpnNatRuleMode(d.Get("mode").(string)),
			Type:             network.VpnNatRuleType(d.Get("type").(string)),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.VirtualNetworkGatewayNatRuleProperties.IPConfigurationID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of the %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}

func resourceVirtualNetworkGatewayNatRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayNatRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NatRuleName)
	d.Set("resource_group_name", id.ResourceGroup)

	vnetGatewayId := parse.NewVirtualNetworkGatewayID(id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkGatewayName)
	d.Set("virtual_network_gateway_id", vnetGatewayId.ID())

	if props := resp.VirtualNetworkGatewayNatRuleProperties; props != nil {
		if err := d.Set("external_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.ExternalMappings)); err != nil {
			return fmt.Errorf("setting `external_mapping`: %+v", err)
		}

		if err := d.Set("internal_mapping", flattenVirtualNetworkGatewayNatRuleMappings(props.InternalMappings)); err != nil {
			return fmt.Errorf("setting `internal_mapping`: %+v", err)
		}

		d.Set("ip_configuration_id", props.IPConfigurationID)
		d.Set("mode", props.Mode)
		d.Set("type", props.Type)
	}

	return nil
}

func resourceVirtualNetworkGatewayNatRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayNatRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	props := network.VirtualNetworkGatewayNatRule{
		Name: utils.String(d.Get("name").(string)),
		VirtualNetworkGatewayNatRuleProperties: &network.VirtualNetworkGatewayNatRuleProperties{
			ExternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("external_mapping").([]interface{})),
			InternalMappings: expandVirtualNetworkGatewayNatRuleMappings(d.Get("internal_mapping").([]interface{})),
			Mode:             network.VpnNatRuleMode(d.Get("mode").(string)),
			Type:             network.VpnNatRuleType(d.Get("type").(string)),
		},
	}

	if v, ok := d.GetOk("ip_configuration_id"); ok {
		props.VirtualNetworkGatewayNatRuleProperties.IPConfigurationID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of the %s: %+v", id, err)
	}

	return resourceVirtualNetworkGatewayNatRuleRead(d, meta)
}

func resourceVirtualNetworkGatewayNatRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayNatRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualNetworkGatewayName, id.NatRuleName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
	}

	return nil
}

func expandVirtualNetworkGatewayNatRuleMappings(input []interface{}) *[]network.VpnNatRuleMapping {
	results := make([]network.VpnNatRuleMapping, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.VpnNatRuleMapping{
			AddressSpace: utils.String(v["address_space"].(string)),
		}

		if portRange := v["port_range"].(string); portRange != "" {
			result.PortRange = utils.String(portRange)
		}

		results = append(results, result)
	}

	return &results
}

func flattenVirtualNetworkGatewayNatRuleMappings(input *[]network.VpnNatRuleMapping) []interface{} {
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
