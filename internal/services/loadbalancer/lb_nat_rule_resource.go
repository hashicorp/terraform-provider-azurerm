// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerNatRule() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceArmLoadBalancerNatRuleCreateUpdate,
		Read:   resourceArmLoadBalancerNatRuleRead,
		Update: resourceArmLoadBalancerNatRuleCreateUpdate,
		Delete: resourceArmLoadBalancerNatRuleDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := loadbalancers.ParseInboundNatRuleID(input)
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
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

			"frontend_port": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ValidateFunc:  validate.PortNumberOrZero,
				ConflictsWith: []string{"frontend_port_start", "frontend_port_end", "backend_address_pool_id"},
			},

			"backend_port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumberOrZero,
			},

			"frontend_ip_configuration_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_floating_ip": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_tcp_reset": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"backend_address_pool_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
				ConflictsWith: []string{"frontend_port"},
				RequiredWith:  []string{"frontend_port_start", "frontend_port_end"},
			},

			"frontend_port_start": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ValidateFunc:  validate.PortNumber,
				RequiredWith:  []string{"backend_address_pool_id", "frontend_port_end"},
				ConflictsWith: []string{"frontend_port"},
			},

			"frontend_port_end": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ValidateFunc:  validate.PortNumber,
				RequiredWith:  []string{"backend_address_pool_id", "frontend_port_start"},
				ConflictsWith: []string{"frontend_port"},
			},

			"idle_timeout_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      4,
				ValidateFunc: validation.IntBetween(4, 30),
			},

			"frontend_ip_configuration_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"backend_ip_configuration_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	return resource
}

func resourceArmLoadBalancerNatRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer Name and Group: %+v", err)
	}
	id := loadbalancers.NewInboundNatRuleID(subscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, d.Get("name").(string))

	loadBalancerIdRaw := loadBalancerId.ID()
	locks.ByID(loadBalancerIdRaw)
	defer locks.UnlockByID(loadBalancerIdRaw)

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
		newNatRule, err := expandAzureRmLoadBalancerNatRule(d, model, *loadBalancerId)
		if err != nil {
			return fmt.Errorf("expanding NAT Rule: %+v", err)
		}

		natRules := make([]loadbalancers.InboundNatRule, 0)
		if props := model.Properties; props != nil {
			natRules = append(*props.InboundNatRules, *newNatRule)

			existingNatRule, existingNatRuleIndex, exists := FindLoadBalancerNatRuleByName(model, id.InboundNatRuleName)
			if exists {
				if id.InboundNatRuleName == *existingNatRule.Name {
					if d.IsNewResource() {
						return tf.ImportAsExistsError("azurerm_lb_nat_rule", *existingNatRule.Id)
					}

					// this nat rule is being updated/reapplied remove old copy from the slice
					natRules = append(natRules[:existingNatRuleIndex], natRules[existingNatRuleIndex+1:]...)
				}
			}

			props.InboundNatRules = &natRules

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
	}
	d.SetId(id.ID())

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func resourceArmLoadBalancerNatRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseInboundNatRuleID(d.Id())
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
		config, _, exists := FindLoadBalancerNatRuleByName(model, id.InboundNatRuleName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Nat Rule %q not found. Removing from state", id.InboundNatRuleName)
			return nil
		}

		d.Set("name", config.Name)
		d.Set("resource_group_name", id.ResourceGroupName)

		if props := config.Properties; props != nil {
			backendIPConfigId := ""
			if props.BackendIPConfiguration != nil && props.BackendIPConfiguration.Id != nil {
				backendIPConfigId = *props.BackendIPConfiguration.Id
			}
			d.Set("backend_ip_configuration_id", backendIPConfigId)
			d.Set("backend_port", pointer.From(props.BackendPort))
			d.Set("enable_floating_ip", pointer.From(props.EnableFloatingIP))
			d.Set("enable_tcp_reset", pointer.From(props.EnableTcpReset))

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

			if props.BackendAddressPool != nil {
				d.Set("backend_address_pool_id", pointer.From(props.BackendAddressPool.Id))
			}
			d.Set("frontend_port", pointer.From(props.FrontendPort))
			d.Set("frontend_port_start", int(pointer.From(props.FrontendPortRangeStart)))
			d.Set("frontend_port_end", int(pointer.From(props.FrontendPortRangeEnd)))
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("protocol", string(pointer.From(props.Protocol)))
		}
	}
	return nil
}

func resourceArmLoadBalancerNatRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseInboundNatRuleID(d.Id())
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
		_, index, exists := FindLoadBalancerNatRuleByName(model, id.InboundNatRuleName)
		if !exists {
			return nil
		}

		if props := model.Properties; props != nil {
			natRules := *props.InboundNatRules
			natRules = append(natRules[:index], natRules[index+1:]...)
			props.InboundNatRules = &natRules

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("Creating/Updating %s: %+v", *id, err)
			}
		}
	}
	return nil
}

func expandAzureRmLoadBalancerNatRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer, loadBalancerId loadbalancers.LoadBalancerId) (*loadbalancers.InboundNatRule, error) {
	properties := loadbalancers.InboundNatRulePropertiesFormat{
		Protocol:       pointer.To(loadbalancers.TransportProtocol(d.Get("protocol").(string))),
		BackendPort:    pointer.To(int64(d.Get("backend_port").(int))),
		EnableTcpReset: pointer.To(d.Get("enable_tcp_reset").(bool)),
	}

	backendAddressPoolSet, frontendPort := false, false
	if port := d.Get("frontend_port"); port != "" {
		frontendPort = true
	}
	if _, ok := d.GetOk("backend_address_pool_id"); ok {
		backendAddressPoolSet = true
	}

	if backendAddressPoolSet {
		properties.FrontendPortRangeStart = pointer.To(int64(d.Get("frontend_port_start").(int)))
		properties.FrontendPortRangeEnd = pointer.To(int64(d.Get("frontend_port_end").(int)))
		properties.BackendAddressPool = &loadbalancers.SubResource{
			Id: pointer.To(d.Get("backend_address_pool_id").(string)),
		}
	} else {
		if frontendPort {
			properties.FrontendPort = pointer.To(int64(d.Get("frontend_port").(int)))
		} else {
			properties.FrontendPortRangeStart = pointer.To(int64(d.Get("frontend_port_start").(int)))
			properties.FrontendPortRangeEnd = pointer.To(int64(d.Get("frontend_port_end").(int)))
		}
	}

	if v, ok := d.GetOk("enable_floating_ip"); ok {
		properties.EnableFloatingIP = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		if _, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v); !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		id := loadbalancers.NewFrontendIPConfigurationID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, v).ID()
		properties.FrontendIPConfiguration = &loadbalancers.SubResource{
			Id: pointer.To(id),
		}
	}

	natRule := loadbalancers.InboundNatRule{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}

	return &natRule, nil
}
