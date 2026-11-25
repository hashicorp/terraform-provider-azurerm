// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_gremlin_graph

import (
	"context"
	"math"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_gremlin_graph/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Gremlin Graph within a Cosmos DB Account.",
		Create:      resourceCosmosDbGremlinGraphCreate,
		Read:        resourceCosmosDbGremlinGraphRead,
		Update:      resourceCosmosDbGremlinGraphUpdate,
		Delete:      resourceCosmosDbGremlinGraphDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseGraphID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.GremlinGraphV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: cosmosdbGremlinGraphSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// `analytical_storage_ttl` can't be disabled once it's enabled
			pluginsdk.ForceNewIfChange("analytical_storage_ttl", func(ctx context.Context, old, new, _ any) bool {
				return (old.(int) == -1 || (old.(int) >= 1 && old.(int) <= math.MaxInt32)) && new.(int) == 0
			}),
		),
	}
}
