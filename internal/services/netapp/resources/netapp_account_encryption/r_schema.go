// Copyright IBM Corp. 2014, 2025
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
			Description:  "The ID of the NetApp Account where encryption will be set.",
			ValidateFunc: validate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			Description:   "The resource ID of the User Assigned Identity to use for encryption.",
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			Description:   "The Principal ID of the System Assigned Identity to use for encryption.",
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
			Description:  "The versionless encryption key url.",
		},

		"federated_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"cross_tenant_key_vault_resource_id"},
			Description:  "The Client ID of the multi-tenant Entra ID application used to access cross-tenant key vaults.",
		},

		"cross_tenant_key_vault_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			RequiredWith: []string{"federated_client_id"},
			Description:  "The full resource ID of the cross-tenant key vault. Required when using federated_client_id for cross-tenant scenarios.",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
