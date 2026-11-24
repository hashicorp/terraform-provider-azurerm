// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_account_encryption

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"netapp_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the NetApp account where volume under it will have customer managed keys-based encryption enabled.",
			ValidateFunc: validate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			Description:   "The ID of the User Assigned Managed Identity.",
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			Description:   "The ID of the System Assigned Managed Identity.",
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
			Description:  "The versionless ID of the encryption key.",
		},

		"federated_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"cross_tenant_key_vault_resource_id"},
			Description:  "The Client ID of the multi-tenant Entra ID application used to access cross-tenant Key Vaults.",
		},

		"cross_tenant_key_vault_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			RequiredWith: []string{"federated_client_id"},
			Description:  "The full resource ID of the cross-tenant Key Vault. This is recommended when using `federated_client_id` for cross-tenant scenarios to ensure proper validation by Azure APIs.",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
