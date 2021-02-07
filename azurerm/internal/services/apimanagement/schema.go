package apimanagement

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

func apiManagementResourceHostnameSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_name": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"key_vault_id": {
			// TODO: should this become `key_vault_key_id` since that's what this is?
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
		},

		"certificate": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"certificate_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"negotiate_client_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func apiManagementResourceHostnameProxySchema() map[string]*schema.Schema {
	hostnameSchema := apiManagementResourceHostnameSchema()

	hostnameSchema["default_ssl_binding"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true, // Azure has certain logic to set this, which we cannot predict
	}

	return hostnameSchema
}
