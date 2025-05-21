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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerOutboundRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Read:   resourceArmLoadBalancerOutboundRuleRead,
		Update: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Delete: resourceArmLoadBalancerOutboundRuleDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := loadbalancers.ParseOutboundRuleID(input)
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

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"frontend_ip_configuration": {
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

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backend_address_pool_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(loadbalancers.TransportProtocolAll),
					string(loadbalancers.TransportProtocolTcp),
					string(loadbalancers.TransportProtocolUdp),
				}, false),
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_tcp_reset": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allocated_outbound_ports": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1024,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  4,
			},
		},
	}
}

func resourceArmLoadBalancerOutboundRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	loadBalancerIDRaw := loadBalancerId.ID()
	id := loadbalancers.NewOutboundRuleID(subscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, d.Get("name").(string))
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

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

	if model := loadBalancer.Model; model != nil {
		newOutboundRule, err := expandAzureRmLoadBalancerOutboundRule(d, model)
		if err != nil {
			return fmt.Errorf("expanding Load Balancer Outbound Rule: %+v", err)
		}

		outboundRules := make([]loadbalancers.OutboundRule, 0)

		if props := model.Properties; props != nil {
			if props.OutboundRules != nil {
				outboundRules = pointer.From(props.OutboundRules)
			}

			existingOutboundRule, existingOutboundRuleIndex, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
			if exists {
				if id.OutboundRuleName == *existingOutboundRule.Name {
					if d.IsNewResource() {
						return tf.ImportAsExistsError("azurerm_lb_outbound_rule", *existingOutboundRule.Id)
					}

					// this outbound rule is being updated/reapplied remove old copy from the slice
					outboundRules = append(outboundRules[:existingOutboundRuleIndex], outboundRules[existingOutboundRuleIndex+1:]...)
				}
			}

			outboundRules = append(outboundRules, *newOutboundRule)

			props.OutboundRules = &outboundRules

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}
		}
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerOutboundRuleRead(d, meta)
}

func resourceArmLoadBalancerOutboundRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseOutboundRuleID(d.Id())
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
		config, _, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Outbound Rule %q not found. Removing from state", id.OutboundRuleName)
			return nil
		}

		d.Set("name", config.Name)

		if props := config.Properties; props != nil {
			allocatedOutboundPorts := 0
			if props.AllocatedOutboundPorts != nil {
				allocatedOutboundPorts = int(*props.AllocatedOutboundPorts)
			}
			d.Set("allocated_outbound_ports", allocatedOutboundPorts)

			backendAddressPoolId := ""
			if props.BackendAddressPool.Id != nil {
				bapid, err := loadbalancers.ParseLoadBalancerBackendAddressPoolIDInsensitively(*props.BackendAddressPool.Id)
				if err != nil {
					return err
				}

				backendAddressPoolId = bapid.ID()
			}
			d.Set("backend_address_pool_id", backendAddressPoolId)
			d.Set("enable_tcp_reset", props.EnableTcpReset)

			frontendIpConfigurations := make([]interface{}, 0)
			if configs := props.FrontendIPConfigurations; configs != nil {
				for _, feConfig := range configs {
					if feConfig.Id == nil {
						continue
					}
					feid, err := loadbalancers.ParseFrontendIPConfigurationIDInsensitively(*feConfig.Id)
					if err != nil {
						return err
					}

					frontendIpConfigurations = append(frontendIpConfigurations, map[string]interface{}{
						"id":   feid.ID(),
						"name": feid.FrontendIPConfigurationName,
					})
				}
			}
			d.Set("frontend_ip_configuration", frontendIpConfigurations)
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("protocol", string(props.Protocol))
		}
	}

	return nil
}

func resourceArmLoadBalancerOutboundRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseOutboundRuleID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

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
		_, index, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
		if !exists {
			return nil
		}

		if props := model.Properties; props != nil {
			outboundRules := *props.OutboundRules
			outboundRules = append(outboundRules[:index], outboundRules[index+1:]...)
			props.OutboundRules = &outboundRules

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
		}
	}
	return nil
}

func expandAzureRmLoadBalancerOutboundRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.OutboundRule, error) {
	properties := loadbalancers.OutboundRulePropertiesFormat{
		Protocol:               loadbalancers.LoadBalancerOutboundRuleProtocol(d.Get("protocol").(string)),
		AllocatedOutboundPorts: pointer.To(int64(d.Get("allocated_outbound_ports").(int))),
	}

	feConfigs := d.Get("frontend_ip_configuration").([]interface{})
	feConfigSubResources := make([]loadbalancers.SubResource, 0)

	for _, raw := range feConfigs {
		v := raw.(map[string]interface{})
		rule, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v["name"].(string))
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v["name"])
		}

		feConfigSubResource := loadbalancers.SubResource{
			Id: rule.Id,
		}

		feConfigSubResources = append(feConfigSubResources, feConfigSubResource)
	}

	properties.FrontendIPConfigurations = feConfigSubResources

	if v := d.Get("backend_address_pool_id").(string); v != "" {
		properties.BackendAddressPool = loadbalancers.SubResource{
			Id: &v,
		}
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("enable_tcp_reset"); ok {
		properties.EnableTcpReset = pointer.To(v.(bool))
	}

	return &loadbalancers.OutboundRule{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}, nil
}
