// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_api_key

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_api_key/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAPIKeyCreate,
		Read:   resourceApplicationInsightsAPIKeyRead,
		Delete: resourceApplicationInsightsAPIKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := componentapikeysapis.ParseApiKeyID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiKeyUpgradeV0ToV1{},
			1: migration.ApiKeyUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: applicationInsightsApiKeySchema(),
	}
}
