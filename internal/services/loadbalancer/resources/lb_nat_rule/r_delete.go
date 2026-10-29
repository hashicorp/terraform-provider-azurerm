// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerNatRuleDelete(d *pluginsdk.ResourceData, meta any) error {
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

			if err := client.CreateOrUpdateThenPoll(ctx, plbId, *model); err != nil {
				return fmt.Errorf("Creating/Updating %s: %+v", *id, err)
			}
		}
	}
	return nil
}
