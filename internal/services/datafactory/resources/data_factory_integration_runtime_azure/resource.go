// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_integration_runtime_azure/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureCreate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataFactoryIntegrationRuntimeAzureV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationruntimes.ParseIntegrationRuntimeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryIntegrationRuntimeAzureSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
			if d.Get("interactive_authoring_time_to_live_in_minutes").(int) > 0 {
				if !d.Get("virtual_network_enabled").(bool) {
					return errors.New("when `interactive_authoring_time_to_live_in_minutes` is set, `virtual_network_enabled` must be set to `true`")
				}
			}

			return nil
		},
	}
}
