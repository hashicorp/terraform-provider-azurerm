// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApiManagementResourceHostnameSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"host_name": {
			Description:      "The host name to use.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"key_vault_certificate_id": {
			Description:  "The ID of the Key Vault Secret containing the SSL Certificate, which must be of the type `application/x-pkcs12`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
		},

		"certificate": {
			Description:  "The Base64-encoded certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_password": {
			Description:  "The password associated with the certificate provided above.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"negotiate_client_certificate": {
			Description: "Should Client Certificate Negotiation be enabled for this Hostname?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"ssl_keyvault_identity_client_id": {
			Description:  "The Managed Service Identity Client ID, which has `GET` access to the Key Vault containing the SSL certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"certificate_source": {
			Description: "The source of the certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"certificate_status": {
			Description: "The status of the certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"expiry": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subject": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func ApiManagementResourceHostnameProxySchema() map[string]*pluginsdk.Schema {
	hostnameSchema := ApiManagementResourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		// Note: O+C because Azure has certain logic to set the default SSL binding which cannot be predicted
		Computed: true,
	}

	return hostnameSchema
}
