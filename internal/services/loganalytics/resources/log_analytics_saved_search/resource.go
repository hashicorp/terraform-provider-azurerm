// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_saved_search/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsSavedSearchCreate,
		Read:   resourceLogAnalyticsSavedSearchRead,
		Delete: resourceLogAnalyticsSavedSearchDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := savedsearches.ParseSavedSearchID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SavedSearchV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: logAnalyticsSavedSearchSchema(),
	}
}
