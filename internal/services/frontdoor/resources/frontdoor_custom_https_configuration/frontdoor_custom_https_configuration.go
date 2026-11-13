// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func schemaCustomHttpsConfiguration() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"certificate_source": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(frontdoors.FrontDoorCertificateSourceFrontDoor),
			ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorCertificateSource(), false),
		},
		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"provisioning_substate": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		// NOTE: None of these attributes are valid if
		//       certificate_source is set to FrontDoor
		"azure_key_vault_certificate_secret_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"azure_key_vault_certificate_secret_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"azure_key_vault_certificate_vault_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
