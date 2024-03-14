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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmLoadBalancerNatPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerNatPoolCreateUpdate,
		Read:   resourceArmLoadBalancerNatPoolRead,
		Update: resourceArmLoadBalancerNatPoolCreateUpdate,
		Delete: resourceArmLoadBalancerNatPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := parse.LoadBalancerInboundNatPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
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

			"frontend_port_start": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_port_end": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"backend_port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"frontend_ip_configuration_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"floating_ip_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tcp_reset_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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
		},
	}
}

func resourceArmLoadBalancerNatPoolCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	id := parse.NewLoadBalancerInboundNatPoolID(subscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, d.Get("name").(string))

	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
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
		newNatPool, err := expandAzureRmLoadBalancerNatPool(d, model)
		if err != nil {
			return fmt.Errorf("expanding NAT Pool: %+v", err)
		}

		natPools := make([]loadbalancers.InboundNatPool, 0)
		if props := model.Properties; props != nil {
			natPools = append(*props.InboundNatPools, *newNatPool)
		}

		existingNatPool, existingNatPoolIndex, exists := FindLoadBalancerNatPoolByName(model, id.InboundNatPoolName)
		if exists {
			if id.InboundNatPoolName == *existingNatPool.Name {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_lb_nat_pool", *existingNatPool.Id)
				}

				// this pool is being updated/reapplied remove old copy from the slice
				natPools = append(natPools[:existingNatPoolIndex], natPools[existingNatPoolIndex+1:]...)
			}
		}

		model.Properties.InboundNatPools = &natPools

		err = client.CreateOrUpdateThenPoll(ctx, plbId, *model)
		if err != nil {
			return fmt.Errorf("creating/updating %s : %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerNatPoolRead(d, meta)
}

func resourceArmLoadBalancerNatPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatPoolID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
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
		config, _, exists := FindLoadBalancerNatPoolByName(model, id.InboundNatPoolName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Nat Pool %q not found. Removing from state", id.InboundNatPoolName)
			return nil
		}

		d.Set("name", config.Name)
		d.Set("resource_group_name", id.ResourceGroup)

		if props := config.Properties; props != nil {
			d.Set("backend_port", props.BackendPort)
			d.Set("floating_ip_enabled", pointer.From(props.EnableFloatingIP))
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))

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
			d.Set("frontend_ip_configuration_id", frontendIPConfigID)
			d.Set("frontend_ip_configuration_name", frontendIPConfigName)
			d.Set("frontend_port_end", props.FrontendPortRangeEnd)
			d.Set("frontend_port_start", props.FrontendPortRangeStart)
			d.Set("idle_timeout_in_minutes", int(*props.IdleTimeoutInMinutes))
			d.Set("protocol", string(props.Protocol))
		}
	}

	return nil
}

func resourceArmLoadBalancerNatPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerInboundNatPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", loadBalancerId, err)
	}

	if model := loadBalancer.Model; model != nil {
		_, index, exists := FindLoadBalancerNatPoolByName(model, id.InboundNatPoolName)
		if !exists {
			return nil
		}

		if props := model.Properties; props != nil {
			natPools := *props.InboundNatPools
			natPools = append(natPools[:index], natPools[index+1:]...)
			props.InboundNatPools = &natPools

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("updating Load Balancer %q (Resource Group %q) for Nat Pool %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.InboundNatPoolName, err)
			}
		}
	}
	return nil
}

func expandAzureRmLoadBalancerNatPool(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.InboundNatPool, error) {
	properties := loadbalancers.InboundNatPoolPropertiesFormat{
		Protocol:               loadbalancers.TransportProtocol(d.Get("protocol").(string)),
		FrontendPortRangeStart: int64(d.Get("frontend_port_start").(int)),
		FrontendPortRangeEnd:   int64(d.Get("frontend_port_end").(int)),
		BackendPort:            int64(d.Get("backend_port").(int)),
	}

	if v, ok := d.GetOk("floating_ip_enabled"); ok {
		properties.EnableFloatingIP = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("tcp_reset_enabled"); ok {
		properties.EnableTcpReset = utils.Bool(v.(bool))
	}

	properties.IdleTimeoutInMinutes = pointer.To(int64(d.Get("idle_timeout_in_minutes").(int)))

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
		if !exists {
			return nil, fmt.Errorf("[ERROR] Cannot find FrontEnd IP Configuration with the name %s", v)
		}

		properties.FrontendIPConfiguration = &loadbalancers.SubResource{
			Id: rule.Id,
		}
	}

	return &loadbalancers.InboundNatPool{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}, nil
}
