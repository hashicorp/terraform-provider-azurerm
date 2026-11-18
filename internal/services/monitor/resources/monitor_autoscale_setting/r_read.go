// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_autoscale_setting

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAutoScaleSettingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := autoscalesettings.ParseAutoScaleSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.AutoScaleSettingName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("location", location.Normalize(model.Location))
		d.Set("enabled", props.Enabled)
		d.Set("target_resource_id", props.TargetResourceUri)

		profile, err := flattenAzureRmMonitorAutoScaleSettingProfile(props.Profiles)
		if err != nil {
			return fmt.Errorf("flattening `profile` of %s: %+v", *id, err)
		}
		if err = d.Set("profile", profile); err != nil {
			return fmt.Errorf("setting `profile` of %s: %+v", *id, err)
		}

		if err = d.Set("predictive", flattenAzureRmMonitorAutoScaleSettingPredictive(props.PredictiveAutoscalePolicy)); err != nil {
			return fmt.Errorf("setting `predictive_scale_mode` of %s: %+v", *id, err)
		}

		if err = d.Set("notification", flattenAzureRmMonitorAutoScaleSettingNotification(props.Notifications)); err != nil {
			return fmt.Errorf("setting `notification` of %s: %+v", *id, err)
		}

		// Return a new tag map filtered by the specified tag names.
		tagMap := tags.Filter(model.Tags, "$type")

		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(tagMap)); err != nil {
			return err
		}
	}
	return nil
}

func validateAutoScaleSettingsTimeZone() pluginsdk.SchemaValidateFunc {
	// from https://docs.microsoft.com/rest/api/monitor/autoscalesettings/createorupdate#timewindow
	timeZones := []string{
		"Dateline Standard Time",
		"UTC-11",
		"Hawaiian Standard Time",
		"Alaskan Standard Time",
		"Pacific Standard Time (Mexico)",
		"Pacific Standard Time",
		"US Mountain Standard Time",
		"Mountain Standard Time (Mexico)",
		"Mountain Standard Time",
		"Central America Standard Time",
		"Central Standard Time",
		"Central Standard Time (Mexico)",
		"Canada Central Standard Time",
		"SA Pacific Standard Time",
		"Eastern Standard Time",
		"US Eastern Standard Time",
		"Venezuela Standard Time",
		"Paraguay Standard Time",
		"Atlantic Standard Time",
		"Central Brazilian Standard Time",
		"SA Western Standard Time",
		"Pacific SA Standard Time",
		"Newfoundland Standard Time",
		"E. South America Standard Time",
		"Argentina Standard Time",
		"SA Eastern Standard Time",
		"Greenland Standard Time",
		"Montevideo Standard Time",
		"Bahia Standard Time",
		"UTC-02",
		"Mid-Atlantic Standard Time",
		"Azores Standard Time",
		"Cape Verde Standard Time",
		"Morocco Standard Time",
		"UTC",
		"GMT Standard Time",
		"Greenwich Standard Time",
		"W. Europe Standard Time",
		"Central Europe Standard Time",
		"Romance Standard Time",
		"Central European Standard Time",
		"W. Central Africa Standard Time",
		"Namibia Standard Time",
		"Jordan Standard Time",
		"GTB Standard Time",
		"Middle East Standard Time",
		"Egypt Standard Time",
		"Syria Standard Time",
		"E. Europe Standard Time",
		"South Africa Standard Time",
		"FLE Standard Time",
		"Turkey Standard Time",
		"Israel Standard Time",
		"Kaliningrad Standard Time",
		"Libya Standard Time",
		"Arabic Standard Time",
		"Arab Standard Time",
		"Belarus Standard Time",
		"Russian Standard Time",
		"E. Africa Standard Time",
		"Iran Standard Time",
		"Arabian Standard Time",
		"Azerbaijan Standard Time",
		"Russia Time Zone 3",
		"Mauritius Standard Time",
		"Georgian Standard Time",
		"Caucasus Standard Time",
		"Afghanistan Standard Time",
		"West Asia Standard Time",
		"Ekaterinburg Standard Time",
		"Pakistan Standard Time",
		"India Standard Time",
		"Sri Lanka Standard Time",
		"Nepal Standard Time",
		"Central Asia Standard Time",
		"Bangladesh Standard Time",
		"N. Central Asia Standard Time",
		"Myanmar Standard Time",
		"SE Asia Standard Time",
		"North Asia Standard Time",
		"China Standard Time",
		"North Asia East Standard Time",
		"Singapore Standard Time",
		"W. Australia Standard Time",
		"Taipei Standard Time",
		"Ulaanbaatar Standard Time",
		"Tokyo Standard Time",
		"Korea Standard Time",
		"Yakutsk Standard Time",
		"Cen. Australia Standard Time",
		"AUS Central Standard Time",
		"E. Australia Standard Time",
		"AUS Eastern Standard Time",
		"West Pacific Standard Time",
		"Tasmania Standard Time",
		"Magadan Standard Time",
		"Vladivostok Standard Time",
		"Russia Time Zone 10",
		"Central Pacific Standard Time",
		"Russia Time Zone 11",
		"New Zealand Standard Time",
		"UTC+12",
		"Fiji Standard Time",
		"Kamchatka Standard Time",
		"Tonga Standard Time",
		"Samoa Standard Time",
		"Line Islands Standard Time",
	}
	return validation.StringInSlice(timeZones, false)
}
