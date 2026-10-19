// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func marketplaceAgreementSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"offer": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"publisher": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"license_text_link": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"privacy_policy_link": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
