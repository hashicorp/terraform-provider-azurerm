// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cosmosdbSqlContainerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cosmos DB SQL Container.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Description:  "The name of the Cosmos DB Account to create the container within.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosAccountName,
		},

		"database_name": {
			Description:  "The name of the Cosmos DB SQL Database to create the container within.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CosmosEntityName,
		},

		// lintignore:S013
		"partition_key_paths": {
			Description: "A list of partition key paths.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"partition_key_kind": {
			Description: "Define a partition key kind.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(cosmosdb.PartitionKindHash),
			ValidateFunc: validation.StringInSlice([]string{
				string(cosmosdb.PartitionKindHash),
				string(cosmosdb.PartitionKindMultiHash),
			}, false),
		},

		"partition_key_version": {
			Description:  "Define a partition key version. This should be set to `2` to use large partition keys.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		},

		"conflict_resolution_policy": common.ConflictResolutionPolicy(),

		"throughput": {
			Description:  "The throughput of SQL container (RU/s). The minimum value is `400`. This must be set upon container creation otherwise it cannot be updated without a manual terraform destroy-apply.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.CosmosThroughput,
		},

		"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),

		"analytical_storage_ttl": {
			Description:  "If present and the value is set to `-1`, it is equal to infinity, and items don't expire by default. If present and the value is set to some number `n` - items will expire `n` seconds after their last modified time.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"default_ttl": {
			Description:  "If missing, items are not expired automatically. If present and the value is set to `-1`, it is equal to infinity, and items don't expire by default. If present and the value is set to some number `n` - items will expire `n` seconds after their last modified time.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

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
		"indexing_policy": common.CosmosDbIndexingPolicySchema(),
	}
}
