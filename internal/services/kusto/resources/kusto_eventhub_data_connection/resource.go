// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventhub_data_connection

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_eventhub_data_connection/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoEventHubDataConnectionCreate,
		Read:   resourceKustoEventHubDataConnectionRead,
		Update: resourceKustoEventHubDataConnectionUpdate,
		Delete: resourceKustoEventHubDataConnectionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseDataConnectionEventHubV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := dataconnections.ParseDataConnectionID(id)
			return err
		}, helpers.ImportDataConnection("EventHub")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: kustoEventhubDataConnectionSchema(),
	}
}
