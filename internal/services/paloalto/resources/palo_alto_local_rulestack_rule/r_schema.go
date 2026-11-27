// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulesresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Palo Alto Local Rulestack Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackRuleName, // TODO - Check this
		},

		"rulestack_id": {
			Description:  "The ID of the Local Rulestack in which to create this Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: localrulesresources.ValidateLocalRulestackID,
		},

		"priority": {
			Description:  "The Priority of this rule. Rules are executed in numerical order.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 1000000),
		},

		"action": {
			Description:  "The action to take on the rule being triggered.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(localrulesresources.PossibleValuesForActionEnum(), false),
		},

		// Optional

		"applications": {
			Description: "A list of Applications.",
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

		"category": schema.CategorySchema(),

		"decryption_rule_type": {
			Description:  "The type of Decryption to perform on the rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      localrulesresources.DecryptionRuleTypeEnumNone,
			ValidateFunc: validation.StringInSlice(localrulesresources.PossibleValuesForDecryptionRuleTypeEnum(), false),
		},

		"description": {
			Description: "The description for the rule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"destination": schema.DestinationSchema(),

		"logging_enabled": {
			Description: "Should Logging be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"inspection_certificate_id": {
			Description:  "The ID of the certificate for inbound inspection. Only valid when `decryption_rule_type` is set to `SSLInboundInspection`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: certificateobjectlocalrulestackresources.ValidateLocalRulestackCertificateID,
		},

		"negate_destination": {
			Description: "Should the inverse of the Destination configuration be used?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"negate_source": {
			Description: "Should the inverse of the Source configuration be used?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"enabled": {
			Description: "Should this Rule be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"protocol": {
			Description:  "The Protocol and port to use in the form `[protocol]:[port_number]` e.g. `TCP:8080` or `UDP:53`, a port range in the form `[protocol]:[start_port]-[end_port]` e.g. `TCP:1024-1206`, or one of `any` and `application-default`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.Protocol,
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
		},

		"protocol_ports": {
			Description:  "A list of Protocol:Port entries. E.g. `[ \"TCP:80\", \"UDP:5431\" ]`. Each entry may also be a port range in the form `[protocol]:[start_port]-[end_port]` e.g. `TCP:1024-1206`.",
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.ProtocolWithPort,
			},
		},

		"source": schema.SourceSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
