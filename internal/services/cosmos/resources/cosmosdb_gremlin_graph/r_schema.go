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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosAccountName,
		},

		"database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"analytical_storage_ttl": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ValidateFunc: validation.All(
				validation.IntBetween(-1, math.MaxInt32),
				validation.IntNotInSlice([]int{0}),
			),
		},

		"default_ttl": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"throughput": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.CosmosThroughput,
		},

		"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

		"partition_key_path": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partition_key_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		},

		"index_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					// case change in 2021-01-15, issue https://github.com/Azure/azure-rest-api-specs/issues/14051
					"indexing_mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForIndexingMode(), false),
					},

					"included_paths": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
					},

					"excluded_paths": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
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
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
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
