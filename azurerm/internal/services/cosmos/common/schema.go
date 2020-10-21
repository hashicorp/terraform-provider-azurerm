package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func DatabaseAutoscaleSettingsSchema() *schema.Schema {
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

func ContainerAutoscaleSettingsSchema() *schema.Schema {
	autoscaleSettingsDatabaseSchema := DatabaseAutoscaleSettingsSchema()
	autoscaleSettingsDatabaseSchema.RequiredWith = []string{"partition_key_path"}

	return autoscaleSettingsDatabaseSchema
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
				"indexing_mode": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          documentdb.Consistent,
					DiffSuppressFunc: suppress.CaseDifference, // Open issue https://github.com/Azure/azure-sdk-for-go/issues/6603
					ValidateFunc: validation.StringInSlice([]string{
						string(documentdb.Consistent),
						string(documentdb.None),
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
										"order": {
											Type:     schema.TypeString,
											Required: true,
											// Workaround for Azure/azure-rest-api-specs#11222
											DiffSuppressFunc: suppress.CaseDifference,
											ValidateFunc: validation.StringInSlice(
												[]string{
													string(documentdb.Ascending),
													string(documentdb.Descending),
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
