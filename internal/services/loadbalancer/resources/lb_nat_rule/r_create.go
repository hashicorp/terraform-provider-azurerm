// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerNatRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
						if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
							return tf.ImportAsExistsError("azurerm_lb_nat_rule", *existingNatRule.Id)
						}
					}

					// this nat rule is being updated/reapplied remove old copy from the slice
					natRules = append(natRules[:existingNatRuleIndex], natRules[existingNatRuleIndex+1:]...)
				}
			}

			props.InboundNatRules = &natRules

			if d.IsNewResource() {
				if err := client.CreateOrUpdateCallbackThenPoll(ctx, plbId, *model, sdk.SetIDCallback(meta, &id, d)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
				d.SetId(id.ID())
			} else {
				if err := client.CreateOrUpdateThenPoll(ctx, plbId, *model); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
		}
	}

	return resourceArmLoadBalancerNatRuleRead(d, meta)
}

func expandAzureRmLoadBalancerNatRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer, loadBalancerId loadbalancers.LoadBalancerId) (*loadbalancers.InboundNatRule, error) {
	properties := loadbalancers.InboundNatRulePropertiesFormat{
		Protocol:         pointer.ToEnum[loadbalancers.TransportProtocol](d.Get("protocol").(string)),
		BackendPort:      pointer.To(int64(d.Get("backend_port").(int))),
		EnableFloatingIP: pointer.To(d.Get("floating_ip_enabled").(bool)),
		EnableTcpReset:   pointer.To(d.Get("tcp_reset_enabled").(bool)),
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

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		if _, exists := helpers.FindLoadBalancerFrontEndIpConfigurationByName(lb, v); !exists {
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
