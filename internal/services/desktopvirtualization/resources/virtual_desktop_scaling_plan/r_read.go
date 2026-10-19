// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scalingplan.ParseScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ScalingPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("description", model.Properties.Description)
		d.Set("friendly_name", model.Properties.FriendlyName)
		d.Set("time_zone", model.Properties.TimeZone)
		d.Set("exclusion_tag", model.Properties.ExclusionTag)
		d.Set("schedule", flattenScalingPlanSchedule(model.Properties.Schedules))
		d.Set("host_pool", flattenScalingHostpoolReference(model.Properties.HostPoolReferences))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenScalingPlanSchedule(input *[]scalingplan.ScalingSchedule) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		rampUpStartTime := ""
		if item.RampUpStartTime != nil {
			rampUpStartTime = fmt.Sprintf("%02d:%02d", item.RampUpStartTime.Hour, item.RampUpStartTime.Minute)
		}
		peakStartTime := ""
		if item.PeakStartTime != nil {
			peakStartTime = fmt.Sprintf("%02d:%02d", item.PeakStartTime.Hour, item.PeakStartTime.Minute)
		}
		rampDownStartTime := ""
		if item.RampDownStartTime != nil {
			rampDownStartTime = fmt.Sprintf("%02d:%02d", item.RampDownStartTime.Hour, item.RampDownStartTime.Minute)
		}
		offPeakStartTime := ""
		if item.OffPeakStartTime != nil {
			offPeakStartTime = fmt.Sprintf("%02d:%02d", item.OffPeakStartTime.Hour, item.OffPeakStartTime.Minute)
		}
		daysOfWeek := make([]string, 0)
		if item.DaysOfWeek != nil {
			for _, weekday := range *item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, string(weekday))
			}
		}

		results = append(results, map[string]any{
			"name":                                 pointer.From(item.Name),
			"days_of_week":                         daysOfWeek,
			"ramp_up_start_time":                   rampUpStartTime,
			"ramp_up_load_balancing_algorithm":     item.RampUpLoadBalancingAlgorithm,
			"ramp_up_minimum_hosts_percent":        pointer.From(item.RampUpMinimumHostsPct),
			"ramp_up_capacity_threshold_percent":   pointer.From(item.RampUpCapacityThresholdPct),
			"peak_start_time":                      peakStartTime,
			"peak_load_balancing_algorithm":        item.PeakLoadBalancingAlgorithm,
			"ramp_down_start_time":                 rampDownStartTime,
			"ramp_down_load_balancing_algorithm":   item.RampDownLoadBalancingAlgorithm,
			"ramp_down_minimum_hosts_percent":      pointer.From(item.RampDownMinimumHostsPct),
			"ramp_down_capacity_threshold_percent": pointer.From(item.RampDownCapacityThresholdPct),
			"ramp_down_force_logoff_users":         pointer.From(item.RampDownForceLogoffUsers),
			"ramp_down_stop_hosts_when":            item.RampDownStopHostsWhen,
			"ramp_down_wait_time_minutes":          pointer.From(item.RampDownWaitTimeMinutes),
			"ramp_down_notification_message":       pointer.From(item.RampDownNotificationMessage),
			"off_peak_start_time":                  offPeakStartTime,
			"off_peak_load_balancing_algorithm":    item.OffPeakLoadBalancingAlgorithm,
		})
	}
	return results
}

func flattenScalingHostpoolReference(input *[]scalingplan.ScalingHostPoolReference) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"hostpool_id":          pointer.From(item.HostPoolArmPath),
			"scaling_plan_enabled": pointer.From(item.ScalingPlanEnabled),
		})
	}
	return results
}

var scalingPlanResourceType = "azurerm_virtual_desktop_scaling_plan"

func validateTime() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`), `The time must be in the format HH:MM.`)
}
