// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CassandraTableSchemaPropertySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"column": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Description:  "Name of the column to be created.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Description:  "Type of the column to be created.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"partition_key": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Description:  "Name of the column to partition by.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"cluster_key": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Description:  "Name of the cluster key to be created.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"order_by": {
								Description: "Order of the key.",
								Type:        pluginsdk.TypeString,
								Required:    true,
								ValidateFunc: validation.StringInSlice([]string{
									"Asc",
									"Desc",
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func DatabaseAutoscaleSettingsSchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_throughput": {
					Description:   "The maximum throughput (RU/s).",
					Type:          pluginsdk.TypeInt,
					Optional:      true,
					Computed:      true, // azignore:AZS007 - pre-existing violation
					ValidateFunc:  validate.CosmosMaxThroughput,
					ConflictsWith: []string{"throughput"},
				},
			},
		},
	}
}

func CosmosDbIndexingPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true, // azignore:AZS007 - pre-existing violation
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// `automatic` is excluded as it is deprecated; see https://stackoverflow.com/a/58721386
				// `indexing_mode` case changes from 2020-04-01 to 2021-01-15 issue https://github.com/Azure/azure-rest-api-specs/issues/14051
				"indexing_mode": {
					Description: "The indexing mode.",
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Default:     cosmosdb.IndexingModeConsistent,
					ValidateFunc: validation.StringInSlice([]string{
						string(cosmosdb.IndexingModeConsistent),
						string(cosmosdb.IndexingModeNone),
					}, false),
				},

				"included_path": {
					Description: "Either `included_path` or `excluded_path` must contain the `path` `/*`.",
					Type:        pluginsdk.TypeList,
					Optional:    true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Description:  "Path for which the indexing behaviour applies to.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"excluded_path": {
					Description: "Either `included_path` or `excluded_path` must contain the `path` `/*`.",
					Type:        pluginsdk.TypeList,
					Optional:    true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Description:  "Path that is excluded from indexing.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"composite_index": CosmosDbIndexingPolicyCompositeIndexSchema(),

				"spatial_index": CosmosDbIndexingPolicySpatialIndexSchema(),

				"full_text_index": CosmosDbIndexingPolicyFullTextIndexSchema(),
			},
		},
	}
}

func ConflictResolutionPolicy() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true, // azignore:AZS007 - pre-existing violation
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mode": {
					Description:  "The conflict resolution mode.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForConflictResolutionMode(), false),
				},

				"conflict_resolution_path": {
					Description:  "The conflict resolution path in the case of `LastWriterWins` mode.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"conflict_resolution_procedure": {
					Description:  "The procedure to resolve conflicts in the case of `Custom` mode.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func CosmosDbIndexingPolicyCompositeIndexSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"index": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Description:  "Path for which the indexing behaviour applies to.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							// `order` case changes from 2020-04-01 to 2021-01-15, issue opened:https://github.com/Azure/azure-rest-api-specs/issues/14051
							"order": {
								Description:  "Order of the index.",
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(cosmosdb.PossibleValuesForCompositePathSortOrder(), true),
								// Workaround for Azure/azure-rest-api-specs#11222
								DiffSuppressFunc: suppress.CaseDifference,
							},
						},
					},
				},
			},
		},
	}
}

func CosmosDbIndexingPolicySpatialIndexSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"path": {
					Description:  "Path for which the indexing behaviour applies to. According to the service design, all spatial types including `LineString`, `MultiPolygon`, `Point`, and `Polygon` will be applied to the path.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"types": {
					Description: "A set of spatial types of the path.",
					Type:        pluginsdk.TypeSet,
					Computed:    true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func CosmosDbIndexingPolicyFullTextIndexSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func CosmosDbFullTextPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"default_language": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"full_text_path": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"language": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}
