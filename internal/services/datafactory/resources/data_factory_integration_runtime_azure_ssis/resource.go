// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_integration_runtime_azure_ssis/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureSsisCreate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureSsisRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureSsisUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureSsisDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataFactoryIntegrationRuntimeAzureSsisV0ToV1{},
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

		Schema: dataFactoryIntegrationRuntimeAzureSsisSchema(),
	}
}
