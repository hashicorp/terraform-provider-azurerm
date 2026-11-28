// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_fqdn_list

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistlocalrulestackresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Palo Alto Local Rulestack FQDN List.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackFQDNListName, // TODO - Check this
		},

		"rulestack_id": {
			Description:  "The ID of the Palo Alto Local Rulestack.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: fqdnlistlocalrulestackresources.ValidateLocalRulestackID,
		},

		"fully_qualified_domain_names": {
			Description: "A list of Fully Qualified Domain Names.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"audit_comment": {
			Description: "The comment for Audit purposes.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"description": {
			Description: "The description for the FQDN List.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
