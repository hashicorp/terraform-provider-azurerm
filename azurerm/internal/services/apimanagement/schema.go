package apimanagement

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateKeyVaultChildIdVersionOptional,
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
