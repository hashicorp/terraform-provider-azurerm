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
			Description:  "The ID of the NetApp Account where encryption will be set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Description:   "The resource ID of the User Assigned Identity to use for encryption.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Description:   "The Principal ID of the System Assigned Identity to use for encryption.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption_key": {
			Description:  "The versionless encryption key URL.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"federated_client_id": {
			Description:  "The Client ID of the multi-tenant Entra ID application used to access cross-tenant Key Vaults.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"cross_tenant_key_vault_resource_id"},
		},

		"cross_tenant_key_vault_resource_id": {
			Description:  "The full resource ID of the cross-tenant Key Vault. Required when using federated_client_id for cross-tenant scenarios.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			RequiredWith: []string{"federated_client_id"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
