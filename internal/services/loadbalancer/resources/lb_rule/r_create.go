// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_rule

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

func resourceArmLoadBalancerRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
				if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return tf.ImportAsExistsError("azurerm_lb_rule", *existingRule.Id)
				}
			}

			// this rule is being updated/reapplied remove old copy from the slice
			lbRules = append(lbRules[:existingRuleIndex], lbRules[existingRuleIndex+1:]...)
		}
	}

	loadBalancer.Model.Properties.LoadBalancingRules = &lbRules

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, plbId, *loadBalancer.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, plbId, *loadBalancer.Model); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceArmLoadBalancerRuleRead(d, meta)
}

func expandAzureRmLoadBalancerRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.LoadBalancingRule, error) {
	properties := loadbalancers.LoadBalancingRulePropertiesFormat{
		Protocol:            loadbalancers.TransportProtocol(d.Get("protocol").(string)),
		FrontendPort:        int64(d.Get("frontend_port").(int)),
		BackendPort:         pointer.To(int64(d.Get("backend_port").(int))),
		EnableFloatingIP:    pointer.To(d.Get("floating_ip_enabled").(bool)),
		EnableTcpReset:      pointer.To(d.Get("tcp_reset_enabled").(bool)),
		DisableOutboundSnat: pointer.To(d.Get("disable_outbound_snat").(bool)),
	}

	if v, ok := d.GetOk("idle_timeout_in_minutes"); ok {
		properties.IdleTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	if v := d.Get("load_distribution").(string); v != "" {
		properties.LoadDistribution = pointer.ToEnum[loadbalancers.LoadDistribution](v)
	}

	// TODO: ensure these ID's are consistent
	if v := d.Get("frontend_ip_configuration_name").(string); v != "" {
		rule, exists := helpers.FindLoadBalancerFrontEndIpConfigurationByName(lb, v)
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

	if l := d.Get("backend_address_pool_ids").([]any); len(l) != 0 {
		if isGateway {
			var baps []loadbalancers.SubResource
			for _, p := range l {
				baps = append(baps, loadbalancers.SubResource{
					Id: pointer.To(p.(string)),
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
