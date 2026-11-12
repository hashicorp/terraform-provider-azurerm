// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerNatPoolCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_lb_nat_pool", *existingNatPool.Id)
					}
				}

				// this pool is being updated/reapplied remove old copy from the slice
				natPools = append(natPools[:existingNatPoolIndex], natPools[existingNatPoolIndex+1:]...)
			}
		}

		model.Properties.InboundNatPools = &natPools

		// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
		if err = client.CreateOrUpdateThenPoll(ctx, plbId, *model); err != nil {
			return fmt.Errorf("creating/updating %s : %+v", id, err)
		}
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceArmLoadBalancerNatPoolRead(d, meta)
}

func expandAzureRmLoadBalancerNatPool(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.InboundNatPool, error) {
	properties := loadbalancers.InboundNatPoolPropertiesFormat{
		Protocol:               loadbalancers.TransportProtocol(d.Get("protocol").(string)),
		FrontendPortRangeStart: int64(d.Get("frontend_port_start").(int)),
		FrontendPortRangeEnd:   int64(d.Get("frontend_port_end").(int)),
		BackendPort:            int64(d.Get("backend_port").(int)),
	}

	if v, ok := d.GetOk("floating_ip_enabled"); ok {
		properties.EnableFloatingIP = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("tcp_reset_enabled"); ok {
		properties.EnableTcpReset = pointer.To(v.(bool))
	}

	properties.IdleTimeoutInMinutes = pointer.To(int64(d.Get("idle_timeout_in_minutes").(int)))

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := helpers.FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
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
