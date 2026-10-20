// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_attached_database_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/attacheddatabaseconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_attached_database_configuration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoAttachedDatabaseConfigurationCreateUpdate,
		Read:   resourceKustoAttachedDatabaseConfigurationRead,
		Update: resourceKustoAttachedDatabaseConfigurationCreateUpdate,
		Delete: resourceKustoAttachedDatabaseConfigurationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoAttachedDatabaseConfigurationV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := attacheddatabaseconfigurations.ParseAttachedDatabaseConfigurationID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			databaseName := diff.Get("database_name").(string)
			databaseNameOverride := diff.Get("database_name_override").(string)

			if databaseName == "*" && databaseNameOverride != "" {
				return fmt.Errorf("cannot set `database_name_override` when `database_name` is set to `*` (all databases)")
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: kustoAttachedDatabaseConfigurationSchema(),
	}
}
