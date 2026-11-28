// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Private DNS Resolver Forwarding Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ForwardingRuleName(),
		},

		"dns_forwarding_ruleset_id": {
			Description:  "The ID of the Private DNS Resolver Forwarding Ruleset.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dnsforwardingrulesets.ValidateDnsForwardingRulesetID,
		},

		"domain_name": {
			Description:  "The domain name for the Private DNS Resolver Forwarding Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_dns_servers": {
			Description: "Can be specified multiple times to define multiple target DNS servers.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Description:  "DNS server IP address.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"port": {
						Description: "DNS server port.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},
				},
			},
		},

		"enabled": {
			Description: "The state of the Private DNS Resolver Forwarding Rule.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"metadata": {
			Description: "Metadata attached to the Private DNS Resolver Forwarding Rule.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
