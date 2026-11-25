// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_collection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbMongoCollectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cosmos DB Mongo Collection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "The name of the Cosmos DB Account in which the Cosmos DB Mongo Collection is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosAccountName,
		},

		"database_name": {
			Description:  "The name of the Cosmos DB Mongo Database in which the Cosmos DB Mongo Collection is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		// SDK/api accepts an array.. but only one is allowed
		"shard_key": {
			Description:  "The name of the key to partition on for sharding. There must not be any other unique index keys.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// default TTL is simply an index on _ts with expireAfterOption, given we can't seem to set TTLs on a given index lets expose this to match the portal
		"default_ttl_seconds": {
			Description: "If the value is `-1`, items are not automatically expired.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			ValidateFunc: validation.All(
				validation.IntAtLeast(-1),
				validation.IntNotInSlice([]int{0}),
			),
		},

		"analytical_storage_ttl": {
			Description:  "If present and the value is set to `-1`, it is equal to infinity, and items don't expire by default. If present and the value is set to some number `n` - items will expire `n` seconds after their last modified time.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"throughput": {
			Description:  "The throughput of the MongoDB collection (RU/s). The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.CosmosThroughput,
		},

		"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

		"index": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"keys": {
						Description: "The list of user settable keys for each Cosmos DB Mongo Collection.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"unique": {
						Description: "Whether the index is unique.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"system_indexes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"keys": {
						Description: "The list of system keys which are not settable for each Cosmos DB Mongo Collection.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"unique": {
						Description: "Identifies whether the table contains no duplicate values.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},
	}
}
