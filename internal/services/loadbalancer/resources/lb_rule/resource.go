// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azurermLbRuleResourceName = "azurerm_lb_rule"

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "load_balancer_name:loadbalancer_id,resource_group_name:loadbalancer_id,subscription_id:loadbalancer_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerRuleCreateUpdate,
		Read:   resourceArmLoadBalancerRuleRead,
		Update: resourceArmLoadBalancerRuleCreateUpdate,
		Delete: resourceArmLoadBalancerRuleDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&loadbalancers.LoadBalancingRuleId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&loadbalancers.LoadBalancingRuleId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: loadBalancerRuleSchema(),
	}
}
