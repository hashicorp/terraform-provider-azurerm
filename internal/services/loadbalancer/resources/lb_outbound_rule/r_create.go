// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

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

func resourceArmLoadBalancerOutboundRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
						if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
							return tf.ImportAsExistsError("azurerm_lb_outbound_rule", *existingOutboundRule.Id)
						}
					}

					// this outbound rule is being updated/reapplied remove old copy from the slice
					outboundRules = append(outboundRules[:existingOutboundRuleIndex], outboundRules[existingOutboundRuleIndex+1:]...)
				}
			}

			outboundRules = append(outboundRules, *newOutboundRule)

			props.OutboundRules = &outboundRules

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

	return resourceArmLoadBalancerOutboundRuleRead(d, meta)
}

func expandAzureRmLoadBalancerOutboundRule(d *pluginsdk.ResourceData, lb *loadbalancers.LoadBalancer) (*loadbalancers.OutboundRule, error) {
	properties := loadbalancers.OutboundRulePropertiesFormat{
		Protocol:               loadbalancers.LoadBalancerOutboundRuleProtocol(d.Get("protocol").(string)),
		AllocatedOutboundPorts: pointer.To(int64(d.Get("allocated_outbound_ports").(int))),
		EnableTcpReset:         pointer.To(d.Get("tcp_reset_enabled").(bool)),
	}

	feConfigs := d.Get("frontend_ip_configuration").([]any)
	feConfigSubResources := make([]loadbalancers.SubResource, 0)

	for _, raw := range feConfigs {
		v := raw.(map[string]any)
		rule, exists := helpers.FindLoadBalancerFrontEndIpConfigurationByName(lb, v["name"].(string))
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

	return &loadbalancers.OutboundRule{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}, nil
}
