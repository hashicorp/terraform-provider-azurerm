// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	loadBalancerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerRuleCreateUpdate,
		Read:   resourceArmLoadBalancerRuleRead,
		Update: resourceArmLoadBalancerRuleCreateUpdate,
		Delete: resourceArmLoadBalancerRuleDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := loadbalancers.ParseLoadBalancingRuleID(input)
			if err != nil {
				return nil, err
			}

			lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
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

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	id := loadbalancers.NewLoadBalancingRuleID(subscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, d.Get("name").(string))

	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	if loadBalancer.Model == nil {
		return fmt.Errorf("retrieving Load Balancer %s: `model` was nil", id)
	}
	if loadBalancer.Model.Properties == nil {
		return fmt.Errorf("retrieving Load Balancer %s: `properties` was nil", id)
	}

	newLbRule, err := expandAzureRmLoadBalancerRule(d, loadBalancer.Model)
	if err != nil {
		return fmt.Errorf("expanding Load Balancer Rule: %+v", err)
	}

	lbRules := append(*loadBalancer.Model.Properties.LoadBalancingRules, *newLbRule)

	existingRule, existingRuleIndex, exists := FindLoadBalancerRuleByName(loadBalancer.Model, id.LoadBalancingRuleName)
	if exists {
		if id.LoadBalancingRuleName == *existingRule.Name {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_lb_rule", *existingRule.Id)
			}

			// this rule is being updated/reapplied remove old copy from the slice
			lbRules = append(lbRules[:existingRuleIndex], lbRules[existingRuleIndex+1:]...)
		}
	}

	loadBalancer.Model.Properties.LoadBalancingRules = &lbRules

	err = client.CreateOrUpdateThenPoll(ctx, plbId, *loadBalancer.Model)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRuleRead(d, meta)
}

func resourceArmLoadBalancerRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancingRuleID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", plbId, err)
	}

	if model := loadBalancer.Model; model != nil {
		config, _, exists := FindLoadBalancerRuleByName(model, id.LoadBalancingRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Rule %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}

		d.Set("name", config.Name)

		if props := config.Properties; props != nil {
			d.Set("disable_outbound_snat", pointer.From(props.DisableOutboundSnat))
			d.Set("enable_floating_ip", pointer.From(props.EnableFloatingIP))
			d.Set("enable_tcp_reset", pointer.From(props.EnableTcpReset))
			d.Set("protocol", string(props.Protocol))
			d.Set("backend_port", int(pointer.From(props.BackendPort)))

			// The backendAddressPools is designed for Gateway LB, while the backendAddressPool is designed for other skus.
			// Thought currently the API returns both, but for the sake of stability, we do use different fields here depending on the LB sku.
			var isGateway bool
			if model.Sku != nil && pointer.From(model.Sku.Name) == loadbalancers.LoadBalancerSkuNameGateway {
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
						if p.Id != nil {
							backendAddressPoolIds = append(backendAddressPoolIds, *p.Id)
						}
					}
				}
			} else {
				if props.BackendAddressPool != nil && props.BackendAddressPool.Id != nil {
					backendAddressPoolId = *props.BackendAddressPool.Id
					backendAddressPoolIds = []interface{}{backendAddressPoolId}
				}
			}
			d.Set("backend_address_pool_ids", backendAddressPoolIds)

			frontendIPConfigName := ""
			frontendIPConfigID := ""
			if props.FrontendIPConfiguration != nil && props.FrontendIPConfiguration.Id != nil {
				feid, err := loadbalancers.ParseFrontendIPConfigurationIDInsensitively(*props.FrontendIPConfiguration.Id)
				if err != nil {
					return err
				}

				frontendIPConfigName = feid.FrontendIPConfigurationName
				frontendIPConfigID = feid.ID()
			}
			d.Set("frontend_ip_configuration_name", frontendIPConfigName)
			d.Set("frontend_ip_configuration_id", frontendIPConfigID)
			d.Set("frontend_port", int(props.FrontendPort))
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("load_distribution", string(pointer.From(props.LoadDistribution)))

			probeId := ""
			if props.Probe != nil {
				probeId = pointer.From(props.Probe.Id)
			}
			d.Set("probe_id", probeId)
		}
	}
	return nil
}

func resourceArmLoadBalancerRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancingRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	loadBalancerIDRaw := loadBalancerId.ID()
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", loadBalancerId, err)
	}

	if model := loadBalancer.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.LoadBalancingRules != nil {
				_, index, exists := FindLoadBalancerRuleByName(model, d.Get("name").(string))
				if !exists {
					return nil
				}

				lbRules := *props.LoadBalancingRules
				lbRules = append(lbRules[:index], lbRules[index+1:]...)
				props.LoadBalancingRules = &lbRules

				err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
				if err != nil {
					return fmt.Errorf("Creating/Updating %s: %+v", id, err)
				}
			}
		}
	}

	return nil
}

func expandAzureRmLoadBalancerRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.LoadBalancingRule, error) {
	properties := loadbalancers.LoadBalancingRulePropertiesFormat{
		Protocol:            loadbalancers.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort:        int64(d.Get("frontend_port").(int)),
		BackendPort:         pointer.To(int64(d.Get("backend_port").(int))),
		EnableFloatingIP:    pointer.To(d.Get("enable_floating_ip").(bool)),
		EnableTcpReset:      pointer.To(d.Get("enable_tcp_reset").(bool)),
		DisableOutboundSnat: pointer.To(d.Get("disable_outbound_snat").(bool)),
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	if v := d.Get("load_distribution").(string); v != "" {
		properties.LoadDistribution = pointer.To(loadbalancers.LoadDistribution(v))
	}

	// TODO: ensure these ID's are consistent
	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &loadbalancers.SubResource{
			Id: rule.Id,
		}
	}

	var isGateway bool
	if lb != nil && lb.Sku != nil && pointer.From(lb.Sku.Name) == loadbalancers.LoadBalancerSkuNameGateway {
		isGateway = true
	}

	if l := d.Get("backend_address_pool_ids").([]interface{}); len(l) != 0 {
		if isGateway {
			var baps []loadbalancers.SubResource
			for _, p := range l {
				p := p.(string)
				baps = append(baps, loadbalancers.SubResource{
					Id: &p,
				})
			}
			properties.BackendAddressPools = &baps
		} else {
			if len(l) > 1 {
				return nil, fmt.Errorf(`only Gateway SKU Load Balancer can have more than one "backend_address_pool_ids"`)
			}
			properties.BackendAddressPool = &loadbalancers.SubResource{
				Id: pointer.To(l[0].(string)),
			}
		}
	}

	if v := d.Get("probe_id").(string); v != "" {
		properties.Probe = &loadbalancers.SubResource{
			Id: &v,
		}
	}

	return &loadbalancers.LoadBalancingRule{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}, nil
}

func resourceArmLoadBalancerRuleSchema() map[string]*pluginsdk.Schema {
	resource := map[string]*pluginsdk.Schema{
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
			ValidateFunc: loadbalancers.ValidateLoadBalancerID,
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
				ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
			},
		},

		"protocol": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.TransportProtocolAll),
				string(loadbalancers.TransportProtocolTcp),
				string(loadbalancers.TransportProtocolUdp),
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
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 100),
		},

		"load_distribution": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(loadbalancers.LoadDistributionDefault),
		},
	}

	return resource
}
