// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_sql_aad_admin

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseWorkspaceSqlAadAdminSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
		},

		"login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}
