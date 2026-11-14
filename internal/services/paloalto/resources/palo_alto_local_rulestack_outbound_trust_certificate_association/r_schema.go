// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_outbound_trust_certificate_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (l Resource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"certificate_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificateobjectlocalrulestackresources.ValidateLocalRulestackCertificateID,
		},
	}
}

func (l Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
