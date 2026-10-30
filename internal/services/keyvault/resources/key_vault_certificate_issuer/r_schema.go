// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultCertificateIssuerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateIssuerName,
		},

		"provider_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"DigiCert",
				"GlobalSign",
				"OneCertV2-PrivateCA",
				"OneCertV2-PublicCA",
				"SslAdminV2",
			}, false),
		},

		"org_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"admin": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"first_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"last_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"phone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"password": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
}
