package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

func CassandraTableSchemaPropertySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"column": {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Required:     true,
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Required:     true,
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"partition_key": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Required:     true,
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"cluster_key": {
					Optional: true,
					Type:     schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"order_by": {
								Type:     schema.TypeString,
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

func DatabaseAutoscaleSettingsSchema() *schema.Schema {
	//lintignore:XS003
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"max_throughput": {
					Type:          schema.TypeInt,
					Optional:      true,
					Computed:      true,
					ConflictsWith: []string{"throughput"},
					ValidateFunc:  validate.CosmosMaxThroughput,
				},
			},
		},
	}
}

func MongoCollectionAutoscaleSettingsSchema() *schema.Schema {
	autoscaleSettingsDatabaseSchema := DatabaseAutoscaleSettingsSchema()
	autoscaleSettingsDatabaseSchema.RequiredWith = []string{"shard_key"}

	return autoscaleSettingsDatabaseSchema
}

func CosmosDbIndexingPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// `automatic` is excluded as it is deprecated; see https://stackoverflow.com/a/58721386
				// `indexing_mode` case changes from 2020-04-01 to 2021-01-15 issue https://github.com/Azure/azure-rest-api-specs/issues/14051
				// todo: change to SDK constants and remove translation code in 3.0
				"indexing_mode": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          documentdb.Consistent,
					DiffSuppressFunc: suppress.CaseDifference, // Open issue https://github.com/Azure/azure-sdk-for-go/issues/6603
					ValidateFunc: validation.StringInSlice([]string{
						"Consistent",
						"None",
					}, false),
				},

				"included_path": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"path": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"excluded_path": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"path": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"composite_index": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": {
								Type:     schema.TypeList,
								MinItems: 1,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"path": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										// `order` case changes from 2020-04-01 to 2021-01-15, issue opened:https://github.com/Azure/azure-rest-api-specs/issues/14051
										// todo: change to SDK constants and remove translation code in 3.0
										"order": {
											Type:     schema.TypeString,
											Required: true,
											// Workaround for Azure/azure-rest-api-specs#11222
											DiffSuppressFunc: suppress.CaseDifference,
											ValidateFunc: validation.StringInSlice(
												[]string{
													"Ascending",
													"Descending",
												}, false),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ConflictResolutionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mode": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(documentdb.LastWriterWins),
						string(documentdb.Custom),
					}, false),
				},

				"conflict_resolution_path": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"conflict_resolution_procedure": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
