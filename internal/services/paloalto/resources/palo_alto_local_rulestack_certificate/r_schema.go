// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Palo Alto Networks Rulestack Certificate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackCertificateName,
		},

		"rulestack_id": {
			Description:  "The ID of the Palo Alto Local Rulestack.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificateobjectlocalrulestackresources.ValidateLocalRulestackID,
		},

		"audit_comment": {
			Description: "The comment for Audit purposes.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"description": {
			Description: "The description for the Certificate.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"key_vault_certificate_id": {
			Description:  "The `versionless_id` of the Key Vault Certificate to use.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeCertificate),
			ExactlyOneOf: []string{"self_signed", "key_vault_certificate_id"},
		},

		"self_signed": {
			Description:  "Should a Self Signed Certificate be used?",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			ExactlyOneOf: []string{"key_vault_certificate_id", "self_signed"},
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
