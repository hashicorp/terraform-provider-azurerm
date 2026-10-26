// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_nat_rule_association

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceNatRuleAssociationCreate,
		Read:   resourceNetworkInterfaceNatRuleAssociationRead,
		Delete: resourceNetworkInterfaceNatRuleAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.InboundNatRuleId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkInterfaceNatRuleAssociationSchema(),
	}
}
