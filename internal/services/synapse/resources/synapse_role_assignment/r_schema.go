// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_role_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_role_assignment/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseRoleAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"synapse_workspace_id", "synapse_spark_pool_id"},
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
		},

		"synapse_spark_pool_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"synapse_workspace_id", "synapse_spark_pool_id"},
			ValidateFunc: validate.SparkPoolID,
		},

		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"principal_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"User",
				"Group",
				"ServicePrincipal",
			}, false),
		},

		"role_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
				return migration.MigrateToNewRole(old) == migration.MigrateToNewRole(new)
			},
			ValidateFunc: validation.StringInSlice([]string{
				"Apache Spark Administrator",
				"Synapse Administrator",
				"Synapse Artifact Publisher",
				"Synapse Artifact User",
				"Synapse Compute Operator",
				"Synapse Contributor",
				"Synapse Credential User",
				"Synapse Linked Data Manager",
				"Synapse Monitoring Operator",
				"Synapse SQL Administrator",
				"Synapse User",
			}, false),
		},
	}
}
