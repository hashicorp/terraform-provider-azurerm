// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "load_balancer_name:loadbalancer_id,resource_group_name:loadbalancer_id,subscription_id:loadbalancer_id" --test-name "basicSkuBasic"

const backendAddressPoolResourceName = "azurerm_lb_backend_address_pool"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Update: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&loadbalancers.LoadBalancerBackendAddressPoolId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&loadbalancers.LoadBalancerBackendAddressPoolId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: lbBackendAddressPoolSchema(),
	}
}
