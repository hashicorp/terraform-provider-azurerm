// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbGremlinGraphSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cosmos DB Gremlin Graph.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "The name of the CosmosDB Account to create the Gremlin Graph within.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosAccountName,
		},

		"database_name": {
			Description:  "The name of the Cosmos DB Graph Database in which the Cosmos DB Gremlin Graph is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"analytical_storage_ttl": {
			Description: "The time to live of Analytical Storage for this Cosmos DB Gremlin Graph. If present and the value is set to `-1`, it means never expire.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			ValidateFunc: validation.All(
				validation.IntBetween(-1, math.MaxInt32),
				validation.IntNotInSlice([]int{0}),
			),
		},

		"default_ttl": {
			Description: "If the value is missing or set to `-1`, items don't expire.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
		},

		"throughput": {
			Description:  "The throughput of the Gremlin graph (RU/s). The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.CosmosThroughput,
		},

		"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

		"partition_key_path": {
			Description:  "Define a partition key.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partition_key_version": {
			Description:  "Define a partition key version. This should be set to `2` to use large partition keys.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		},

		"index_policy": {
			Description: "The configuration of the indexing policy.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic": {
						Description: "Whether the indexing policy is automatic.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					// case change in 2021-01-15, issue https://github.com/Azure/azure-rest-api-specs/issues/14051
					"indexing_mode": {
						Description:  "The indexing mode.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForIndexingMode(), false),
					},

					"included_paths": {
						Description: "List of paths to include in the indexing.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Set:         pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"excluded_paths": {
						Description: "List of paths to exclude from indexing.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Set:         pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"composite_index": common.CosmosDbIndexingPolicyCompositeIndexSchema(),

					"spatial_index": common.CosmosDbIndexingPolicySpatialIndexSchema(),
				},
			},
		},

		"conflict_resolution_policy": common.ConflictResolutionPolicy(),

		"unique_key": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"paths": {
						Description: "A list of paths to use for this unique key.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}
