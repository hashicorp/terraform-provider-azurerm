// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/resources/monitor_activity_log_alert/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorActivityLogAlertCreateUpdate,
		Read:   resourceMonitorActivityLogAlertRead,
		Update: resourceMonitorActivityLogAlertCreateUpdate,
		Delete: resourceMonitorActivityLogAlertDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&activitylogalertsapis.ActivityLogAlertId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&activitylogalertsapis.ActivityLogAlertId{}),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ActivityLogAlertUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: monitorActivityLogAlertSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
				// Validate location constraints for Activity Log Alert resources
				loc := diff.Get("location").(string)
				normalizedLocation := location.Normalize(loc)

				supportedLocations := []string{"global", "westeurope", "northeurope", "eastus2euap"}
				supported := false
				for _, supportedLocation := range supportedLocations {
					if normalizedLocation == strings.ToLower(supportedLocation) {
						supported = true
						break
					}
				}

				if !supported {
					return fmt.Errorf("`azurerm_monitor_activity_log_alert` resources are only supported in the following regions: [%s], got `%s`", strings.Join(supportedLocations, ", "), loc)
				}

				return nil
			}),
		),
	}
}
