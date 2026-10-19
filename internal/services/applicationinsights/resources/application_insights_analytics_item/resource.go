// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_analytics_item

import (
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_analytics_item/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAnalyticsItemCreate,
		Read:   resourceApplicationInsightsAnalyticsItemRead,
		Update: resourceApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceApplicationInsightsAnalyticsItemDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			if strings.Contains(id, userScopePath) || strings.Contains(id, sharedScopePath) {
				if _, err := analyticsitemsapis.ParseProviderComponentID(id); err != nil {
					return err
				}
			}
			return nil
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AnalyticsItemUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: applicationInsightsAnalyticsItemSchema(),
	}
}
