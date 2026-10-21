// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func databricksWorkspaceRootDbfsCustomerManagedKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},
	}
}
