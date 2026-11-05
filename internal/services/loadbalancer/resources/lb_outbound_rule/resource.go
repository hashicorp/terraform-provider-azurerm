// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Read:   resourceArmLoadBalancerOutboundRuleRead,
		Update: resourceArmLoadBalancerOutboundRuleCreateUpdate,
		Delete: resourceArmLoadBalancerOutboundRuleDelete,

		Importer: helpers.LoadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
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

		Schema: lbOutboundRuleSchema(),
	}
}
