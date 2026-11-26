// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayNatRuleCreate,
		Read:   resourceVirtualNetworkGatewayNatRuleRead,
		Update: resourceVirtualNetworkGatewayNatRuleUpdate,
		Delete: resourceVirtualNetworkGatewayNatRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := networkgateways.ParseVirtualNetworkGatewayNatRuleID(id)
			return err
		}),

		Schema: virtualNetworkGatewayNatRuleSchema(),
	}
}
