// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_account/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var CosmosDbAccountResourceName = "azurerm_cosmosdb_account"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a CosmosDB (formally DocumentDB) Account.",
		Create:      resourceCosmosDbAccountCreate,
		Read:        resourceCosmosDbAccountRead,
		Update:      resourceCosmosDbAccountUpdate,
		Delete:      resourceCosmosDbAccountDelete,
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("backup.0.type", func(ctx context.Context, old, new, _ any) bool {
				// backup type can only change from Periodic to Continuous
				return old.(string) == string(cosmosdb.BackupPolicyTypeContinuous) && new.(string) == string(cosmosdb.BackupPolicyTypePeriodic)
			}),

			pluginsdk.ForceNewIfChange("analytical_storage_enabled", func(ctx context.Context, old, new, _ any) bool {
				// analytical_storage_enabled can not be changed after being set to true
				return old.(bool) && !new.(bool)
			}),

			pluginsdk.ForceNewIf("capabilities", func(ctx context.Context, d *schema.ResourceDiff, meta any) bool {
				kind := d.Get("kind").(string)
				old, new := d.GetChange("capabilities")

				return !checkCapabilitiesCanBeUpdated(kind, prepareCapabilities(old), prepareCapabilities(new))
			}),

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				caps := diff.Get("capabilities")
				mongo34found := false
				enableMongo := false
				isMongo := strings.EqualFold(diff.Get("kind").(string), string(cosmosdb.DatabaseAccountKindMongoDB))

				for _, cap := range caps.(*pluginsdk.Set).List() {
					m := cap.(map[string]any)
					if v, ok := m["name"].(string); ok {
						switch v {
						case "MongoDBv3.4":
							mongo34found = true
						case "EnableMongo":
							enableMongo = true
						}
					}
				}

				if isMongo && (mongo34found && !enableMongo) {
					return fmt.Errorf("capability EnableMongo must be enabled if MongoDBv3.4 is also enabled")
				}
				return nil
			}),
		),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseDatabaseAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(180 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(300 * time.Minute),
		},

		SchemaVersion: 1,

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CosmosDBAccountV0toV1{},
		}),

		Schema: cosmosdbAccountSchema(),
	}
}
