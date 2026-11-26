// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_transparent_data_encryption

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlManagedInstanceTransparentDataEncryptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_id": {
			Description:  "The name of the MS SQL Managed Instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlManagedInstanceIDInsensitively),
		},

		"key_vault_key_id": {
			Description:      "To use customer managed keys from Azure Key Vault, provide the AKV Key ID. To use service managed keys, omit this field.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
			DiffSuppressFunc: diffSuppressKeyVaultVersionedKey,
		},

		"auto_rotation_enabled": {
			Description: "When enabled, the SQL Managed Instance will continuously check the Key Vault for any new versions of the key being used as the TDE protector. If a new version of the key is detected, the TDE protector on the SQL Managed Instance will be automatically rotated to the latest key version within 60 minutes.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}
