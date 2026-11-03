// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_script

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_script/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabaseScriptCreateUpdate,
		Read:   resourceKustoDatabaseScriptRead,
		Update: resourceKustoDatabaseScriptCreateUpdate,
		Delete: resourceKustoDatabaseScriptDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseScriptV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scripts.ParseScriptID(id)
			return err
		}),

		Schema: kustoScriptSchema(),
	}
}
