// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultAccessPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"tenant_id": {
			Description:  "The Azure Active Directory tenant ID that should be used for authenticating requests to the Key Vault.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"object_id": {
			Description:  "The object ID of a user, Service Principal or security group in the Azure Active Directory tenant for the vault. The object ID of a Service Principal can be fetched from `azuread_service_principal.object_id`. The object ID must be unique for the list of access policies.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"application_id": {
			Description:  "The object ID of an Application in Azure Active Directory.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"certificate_permissions": helpers.SchemaCertificatePermissions(),

		"key_permissions": helpers.SchemaKeyPermissions(),

		"secret_permissions": helpers.SchemaSecretPermissions(),

		"storage_permissions": helpers.SchemaStoragePermissions(),
	}
}
