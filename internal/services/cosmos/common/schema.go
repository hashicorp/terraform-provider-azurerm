// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
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
								Required:     true,
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Required:     true,
								Type:         pluginsdk.TypeString,
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
								Required:     true,
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"cluster_key": {
					Optional: true,
					Type:     pluginsdk.TypeList,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"order_by": {
								Type:     pluginsdk.TypeString,
								Required: true,
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
					Type:          pluginsdk.TypeInt,
					Optional:      true,
					Computed:      true,
					ConflictsWith: []string{"throughput"},
					ValidateFunc:  validate.CosmosMaxThroughput,
				},
			},
		},
	}
}

func CosmosDbIndexingPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// `automatic` is excluded as it is deprecated; see https://stackoverflow.com/a/58721386
				// `indexing_mode` case changes from 2020-04-01 to 2021-01-15 issue https://github.com/Azure/azure-rest-api-specs/issues/14051
				"indexing_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  documentdb.IndexingModeConsistent,
					ValidateFunc: validation.StringInSlice([]string{
						string(documentdb.IndexingModeConsistent),
						string(documentdb.IndexingModeNone),
					}, false),
				},

				"included_path": {
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
				},
				"excluded_path": {
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
				},
				"composite_index": CosmosDbIndexingPolicyCompositeIndexSchema(),

				"spatial_index": CosmosDbIndexingPolicySpatialIndexSchema(),
			},
		},
	}
}

func ConflictResolutionPolicy() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mode": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(documentdb.ConflictResolutionModeLastWriterWins),
						string(documentdb.ConflictResolutionModeCustom),
					}, false),
				},

				"conflict_resolution_path": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"conflict_resolution_procedure": {
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
					MinItems: 1,
					Required: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							// `order` case changes from 2020-04-01 to 2021-01-15, issue opened:https://github.com/Azure/azure-rest-api-specs/issues/14051
							"order": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// Workaround for Azure/azure-rest-api-specs#11222
								DiffSuppressFunc: suppress.CaseDifference,
								ValidateFunc: validation.StringInSlice([]string{
									string(documentdb.CompositePathSortOrderAscending),
									string(documentdb.CompositePathSortOrderDescending),
								}, true),
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
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"types": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}
