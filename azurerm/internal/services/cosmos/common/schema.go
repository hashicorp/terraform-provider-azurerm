package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
