// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementResourceHostnameSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"host_name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"key_vault_id": {
			// TODO: 4.0 - should this become `key_vault_key_id` since that's what this is?
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
		},

		"certificate": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"negotiate_client_certificate": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"ssl_keyvault_identity_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"certificate_source": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"certificate_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

func apiManagementResourceHostnameProxySchema() map[string]*pluginsdk.Schema {
	hostnameSchema := apiManagementResourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Computed: true, // Azure has certain logic to set this, which we cannot predict
	}

	return hostnameSchema
}
