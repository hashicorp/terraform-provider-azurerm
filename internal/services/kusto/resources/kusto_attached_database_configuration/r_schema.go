// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_attached_database_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/attacheddatabaseconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoAttachedDatabaseConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Kusto Attached Database Configuration to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"cluster_name": {
			Description:  "The name of the Kusto Cluster for which the configuration will be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"database_name": {
			Description:  "The name of the database which you would like to attach, use * if you want to follow all current and future databases.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(validate.DatabaseName, validation.StringInSlice([]string{"*"}, false)),
		},

		"cluster_id": {
			Description:  "The resource ID of the cluster where the databases you would like to attach reside.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoClusterID,
		},

		"database_name_override": {
			Description:   "The database name to use for the attached database instead of using the original database name. Relevant only when attaching to a specific database.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validate.DatabaseName,
			ConflictsWith: []string{"database_name_prefix"},
		},

		"database_name_prefix": {
			Description:   "Adds a prefix to the attached databases name. When following an entire cluster, that prefix would be added to all of the databases original names from leader cluster.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validate.DatabaseName,
			ConflictsWith: []string{"database_name_override"},
		},

		"attached_database_names": {
			Description: "The list of databases from the `cluster_id` which are currently attached to the cluster.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"default_principal_modification_kind": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      attacheddatabaseconfigurations.DefaultPrincipalsModificationKindNone,
			ValidateFunc: validation.StringInSlice(attacheddatabaseconfigurations.PossibleValuesForDefaultPrincipalsModificationKind(), false),
		},

		"sharing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"external_tables_to_exclude": {
						Description: "List of external tables exclude from the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"external_tables_to_include": {
						Description: "List of external tables to include in the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"functions_to_exclude": {
						Description: "List of functions to exclude from the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"functions_to_include": {
						Description: "List of functions to include in the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"materialized_views_to_exclude": {
						Description: "List of materialized views exclude from the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"materialized_views_to_include": {
						Description: "List of materialized views to include in the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"tables_to_exclude": {
						Description: "List of tables to exclude from the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"tables_to_include": {
						Description: "List of tables to include in the follower database.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}
