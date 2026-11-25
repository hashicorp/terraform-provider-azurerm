// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_sql_container

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_container/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a SQL Container within a Cosmos DB Account.",
		Create:      resourceCosmosDbSQLContainerCreate,
		Read:        resourceCosmosDbSQLContainerRead,
		Update:      resourceCosmosDbSQLContainerUpdate,
		Delete:      resourceCosmosDbSQLContainerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseContainerID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SqlContainerV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: cosmosdbSqlContainerSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// The analytical_storage_ttl cannot be changed back once enabled on an existing container. -> we need ForceNew
			pluginsdk.ForceNewIfChange("analytical_storage_ttl", func(ctx context.Context, old, new, _ any) bool {
				return (old.(int) == -1 || old.(int) > 0) && new.(int) == 0
			}),

			pluginsdk.ForceNewIfChange("partition_key_version", func(ctx context.Context, old, new, _ any) bool {
				// The behaviour of the Azure API is that `partition_key_version` can be updated to `1` when it is not set at creation time, but it can not be updated to `2`.
				return old.(int) != 0 || new.(int) != 1
			}),
		),
	}
}
