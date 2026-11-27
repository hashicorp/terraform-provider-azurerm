// Copyright IBM Corp. 2014, 2026
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
		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"name": {
			Description:  "The name which should be used for this Key Vault Certificate Issuer.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateIssuerName,
		},

		"provider_name": {
			Description: "The name of the third-party Certificate Issuer.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"DigiCert",
				"GlobalSign",
				"OneCertV2-PrivateCA",
				"OneCertV2-PublicCA",
				"SslAdminV2",
			}, false),
		},

		"org_id": {
			Description:  "The ID of the organization as provided to the issuer.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"account_id": {
			Description: "The account number with the third-party Certificate Issuer.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"admin": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email_address": {
						Description: "E-mail address of the admin.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"first_name": {
						Description: "First name of the admin.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"last_name": {
						Description: "Last name of the admin.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"phone": {
						Description: "Phone number of the admin.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},

		"password": {
			Description: "The password associated with the account and organization ID at the third-party Certificate Issuer. If not specified, will not overwrite any previous value.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Sensitive:   true,
		},
	}
}
