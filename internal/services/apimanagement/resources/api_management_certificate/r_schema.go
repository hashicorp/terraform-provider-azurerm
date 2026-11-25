// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName("The name of the API Management Certificate."),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"data": {
			Description:   "The base-64 encoded certificate data, which must be a PFX file.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsBase64,
			AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
			ConflictsWith: []string{"key_vault_secret_id", "key_vault_identity_client_id"},
		},

		"password": {
			Description:  "The password used for this certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"data"},
		},

		"key_vault_secret_id": {
			Description:   "The ID of the Key Vault Secret containing the SSL Certificate, which must be of the type `application/x-pkcs12`.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
			ConflictsWith: []string{"data", "password"},
		},

		"key_vault_identity_client_id": {
			Description:  "The Client ID of the User Assigned Managed Identity to use for retrieving certificate.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"key_vault_secret_id"},
		},

		"expiration": {
			Description: "The Expiration Date of this Certificate, formatted as an RFC3339 string.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subject": {
			Description: "The Subject of this Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"thumbprint": {
			Description: "The Thumbprint of this Certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
