// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_key

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseWorkspaceKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
		},

		"customer_managed_key_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"customer_managed_key_versionless_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
		},

		"active": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
	}
}
