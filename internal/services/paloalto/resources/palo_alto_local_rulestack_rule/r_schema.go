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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackRuleName, // TODO - Check this
		},

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: localrulesresources.ValidateLocalRulestackID,
		},

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 1000000),
		},

		"action": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(localrulesresources.PossibleValuesForActionEnum(), false),
		},

		// Optional

		"applications": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"audit_comment": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"category": schema.CategorySchema(),

		"decryption_rule_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      localrulesresources.DecryptionRuleTypeEnumNone,
			ValidateFunc: validation.StringInSlice(localrulesresources.PossibleValuesForDecryptionRuleTypeEnum(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"destination": schema.DestinationSchema(),

		"logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"inspection_certificate_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: certificateobjectlocalrulestackresources.ValidateLocalRulestackCertificateID,
		},

		"negate_destination": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"negate_source": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.Any(
				validate.ProtocolWithPort,
				validation.StringInSlice([]string{protocolApplicationDefault}, false),
			),
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
		},

		"protocol_ports": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.ProtocolWithPort,
			},
			ExactlyOneOf: []string{"protocol", "protocol_ports"},
		},

		"source": schema.SourceSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
