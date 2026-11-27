// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_security_admin_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityadminconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Network Manager Security Admin Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Description:  "The ID of the Network Manager Security Admin Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: securityadminconfigurations.ValidateNetworkManagerID,
		},

		"apply_on_network_intent_policy_based_services": {
			Description: "A list of network intent policy based services. Exactly one value should be set. The `All` option requires `Microsoft.Network/AllowAdminRulesOnNipBasedServices` feature registration to Subscription. See [this document](https://learn.microsoft.com/azure/virtual-network-manager/concept-security-admins#network-intent-policies-and-security-admin-rules) for more information.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(securityadminconfigurations.PossibleValuesForNetworkIntentPolicyBasedService(), false),
			},
		},

		"description": {
			Description:  "A description of the Security Admin Configuration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
