// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceArmLoadBalancerRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerRuleCreateUpdate,
		Read:   resourceArmLoadBalancerRuleRead,
		Update: resourceArmLoadBalancerRuleCreateUpdate,
		Delete: resourceArmLoadBalancerRuleDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancingRuleID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
			return &lbId, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceArmLoadBalancerRuleSchema(),
	}
}

func resourceArmLoadBalancerRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLoadBalancingRuleID(subscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, d.Get("name").(string))

	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	loadBalancer, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.Name, err)
	}

	if loadBalancer.LoadBalancerPropertiesFormat == nil {
		return fmt.Errorf("retrieving Load Balancer %s: `properties` was nil", id)
	}

	newLbRule, err := expandAzureRmLoadBalancerRule(d, &loadBalancer)
	if err != nil {
		return fmt.Errorf("expanding Load Balancer Rule: %+v", err)
	}

	lbRules := append(*loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules, *newLbRule)

	existingRule, existingRuleIndex, exists := FindLoadBalancerRuleByName(&loadBalancer, id.Name)
	if exists {
		if id.Name == *existingRule.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_rule", *existingRule.ID)
			}

			// this rule is being updated/reapplied remove old copy from the slice
			lbRules = append(lbRules[:existingRuleIndex], lbRules[existingRuleIndex+1:]...)
		}
	}

	loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = &lbRules

	future, err := client.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, loadBalancer)
	if err != nil {
		return fmt.Errorf("updating Loadbalancer %q (resource group %q) for Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.Name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Load Balancer %q (resource group %q) for Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.Name, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRuleRead(d, meta)
}

func resourceArmLoadBalancerRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancingRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancer, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.Name, err)
	}

	config, _, exists := FindLoadBalancerRuleByName(&loadBalancer, id.Name)
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer Rule %q not found. Removing from state", id.Name)
		return nil
	}

	d.Set("name", config.Name)

	if props := config.LoadBalancingRulePropertiesFormat; props != nil {
		d.Set("disable_outbound_snat", props.DisableOutboundSnat)
		d.Set("enable_floating_ip", props.EnableFloatingIP)
		d.Set("enable_tcp_reset", props.EnableTCPReset)
		d.Set("protocol", string(props.Protocol))

		backendPort := 0
		if props.BackendPort != nil {
			backendPort = int(*props.BackendPort)
		}
		d.Set("backend_port", backendPort)

		// The backendAddressPools is designed for Gateway LB, while the backendAddressPool is designed for other skus.
		// Thought currently the API returns both, but for the sake of stability, we do use different fields here depending on the LB sku.
		var isGateway bool
		if loadBalancer.Sku != nil && loadBalancer.Sku.Name == network.LoadBalancerSkuNameGateway {
			isGateway = true
		}
		var (
			backendAddressPoolId  string
			backendAddressPoolIds []interface{}
		)
		if isGateway {
			// The gateway LB rule can have up to 2 backend address pools.
			// In case there is only one BAP, we set it to both "backendAddressPoolId" and "backendAddressPoolIds".
			// Otherwise, we leave the "backendAddressPoolId" as empty.
			if props.BackendAddressPools != nil {
				for _, p := range *props.BackendAddressPools {
					if p.ID != nil {
						backendAddressPoolIds = append(backendAddressPoolIds, *p.ID)
					}
				}
			}
		} else {
			if props.BackendAddressPool != nil && props.BackendAddressPool.ID != nil {
				backendAddressPoolId = *props.BackendAddressPool.ID
				backendAddressPoolIds = []interface{}{backendAddressPoolId}
			}
		}
		d.Set("backend_address_pool_ids", backendAddressPoolIds)

		frontendIPConfigName := ""
		frontendIPConfigID := ""
		if props.FrontendIPConfiguration != nil && props.FrontendIPConfiguration.ID != nil {
			feid, err := parse.LoadBalancerFrontendIpConfigurationIDInsensitively(*props.FrontendIPConfiguration.ID)
			if err != nil {
				return err
			}

			frontendIPConfigName = feid.FrontendIPConfigurationName
			frontendIPConfigID = feid.ID()
		}
		d.Set("frontend_ip_configuration_name", frontendIPConfigName)
		d.Set("frontend_ip_configuration_id", frontendIPConfigID)

		frontendPort := 0
		if props.FrontendPort != nil {
			frontendPort = int(*props.FrontendPort)
		}
		d.Set("frontend_port", frontendPort)

		idleTimeoutInMinutes := 0
		if props.IdleTimeoutInMinutes != nil {
			idleTimeoutInMinutes = int(*props.IdleTimeoutInMinutes)
		}
		d.Set("idle_timeout_in_minutes", idleTimeoutInMinutes)

		loadDistribution := ""
		if props.LoadDistribution != "" {
			loadDistribution = string(props.LoadDistribution)
		}
		d.Set("load_distribution", loadDistribution)

		probeId := ""
		if props.Probe != nil && props.Probe.ID != nil {
			probeId = *props.Probe.ID
		}
		d.Set("probe_id", probeId)
	}

	return nil
}

func resourceArmLoadBalancerRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancingRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerIDRaw := loadBalancerId.ID()
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

	loadBalancer, err := client.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(loadBalancer.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Rule %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.Name, err)
	}

	if loadBalancer.LoadBalancerPropertiesFormat != nil && loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules != nil {
		_, index, exists := FindLoadBalancerRuleByName(&loadBalancer, d.Get("name").(string))
		if !exists {
			return nil
		}

		lbRules := *loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules
		lbRules = append(lbRules[:index], lbRules[index+1:]...)
		loadBalancer.LoadBalancerPropertiesFormat.LoadBalancingRules = &lbRules

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, loadBalancer)
		if err != nil {
			return fmt.Errorf("Creating/Updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for completion of Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandAzureRmLoadBalancerRule(d *pluginsdk.ResourceData, lb *network.LoadBalancer) (*network.LoadBalancingRule, error) {
	properties := network.LoadBalancingRulePropertiesFormat{
		Protocol:            network.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort:        utils.Int32(int32(d.Get("frontend_port").(int))),
		BackendPort:         utils.Int32(int32(d.Get("backend_port").(int))),
		EnableFloatingIP:    utils.Bool(d.Get("enable_floating_ip").(bool)),
		EnableTCPReset:      utils.Bool(d.Get("enable_tcp_reset").(bool)),
		DisableOutboundSnat: utils.Bool(d.Get("disable_outbound_snat").(bool)),
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = utils.Int32(int32(v.(int)))
	}

	if v := d.Get("load_distribution").(string); v != "" {
		properties.LoadDistribution = network.LoadDistribution(v)
	}

	// TODO: ensure these ID's are consistent
	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &network.SubResource{
			ID: rule.ID,
		}
	}

	var isGateway bool
	if lb != nil && lb.Sku != nil && lb.Sku.Name == network.LoadBalancerSkuNameGateway {
		isGateway = true
	}

	if l := d.Get("backend_address_pool_ids").([]interface{}); len(l) != 0 {
		if isGateway {
			var baps []network.SubResource
			for _, p := range l {
				p := p.(string)
				baps = append(baps, network.SubResource{
					ID: &p,
				})
			}
			properties.BackendAddressPools = &baps
		} else {
			if len(l) > 1 {
				return nil, fmt.Errorf(`only Gateway SKU Load Balancer can have more than one "backend_address_pool_ids"`)
			}
			properties.BackendAddressPool = &network.SubResource{
				ID: utils.String(l[0].(string)),
			}
		}
	}

	if v := d.Get("probe_id").(string); v != "" {
		properties.Probe = &network.SubResource{
			ID: &v,
		}
	}

	return &network.LoadBalancingRule{
		Name:                              utils.String(d.Get("name").(string)),
		LoadBalancingRulePropertiesFormat: &properties,
	}, nil
}

func resourceArmLoadBalancerRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadBalancerValidate.RuleName,
		},

		"loadbalancer_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadBalancerValidate.LoadBalancerID,
		},

		"frontend_ip_configuration_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frontend_ip_configuration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"backend_address_pool_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 2, // Only Gateway SKU LB can have 2 backend address pools
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: loadBalancerValidate.LoadBalancerBackendAddressPoolID,
			},
		},

		"protocol": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.TransportProtocolAll),
				string(network.TransportProtocolTCP),
				string(network.TransportProtocolUDP),
			}, false),
		},

		"frontend_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.PortNumberOrZero,
		},

		"backend_port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validate.PortNumberOrZero,
		},

		"probe_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_floating_ip": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_tcp_reset": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"disable_outbound_snat": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"idle_timeout_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"load_distribution": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
