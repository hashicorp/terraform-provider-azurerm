// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_autoscale_setting

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAutoScaleSettingCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := autoscalesettings.NewAutoScaleSettingID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_monitor_autoscale_setting", id.ID())
			}
		}
	}

	location := location.Normalize(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)
	targetResourceId := d.Get("target_resource_id").(string)

	notificationsRaw := d.Get("notification").([]any)
	notifications := expandAzureRmMonitorAutoScaleSettingNotifications(notificationsRaw)

	profilesRaw := d.Get("profile").([]any)
	profiles, err := expandAzureRmMonitorAutoScaleSettingProfile(profilesRaw)
	if err != nil {
		return fmt.Errorf("expanding `profile`: %+v", err)
	}

	t := d.Get("tags").(map[string]any)

	parameters := autoscalesettings.AutoscaleSettingResource{
		Location: location,
		Properties: autoscalesettings.AutoscaleSetting{
			Enabled:                   &enabled,
			Profiles:                  profiles,
			PredictiveAutoscalePolicy: expandAzureRmMonitorAutoScaleSettingPredictive(d.Get("predictive").([]any)),
			Notifications:             notifications,
			TargetResourceUri:         &targetResourceId,
		},
		Tags: pluginsdk.ExpandPtrMapStringString(t),
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorAutoScaleSettingRead(d, meta)
}
