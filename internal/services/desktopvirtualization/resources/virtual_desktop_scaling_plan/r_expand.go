// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandScalingPlanSchedule(input []any) *[]scalingplan.ScalingSchedule {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingSchedule, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)
		name := v["name"].(string)
		daysOfWeekRaw := v["days_of_week"].(*pluginsdk.Set).List()
		daysOfWeek := make([]scalingplan.DayOfWeek, 0)
		for _, weekday := range daysOfWeekRaw {
			daysOfWeek = append(daysOfWeek, scalingplan.DayOfWeek(weekday.(string)))
		}

		rampUpStartTime := v["ramp_up_start_time"].(string)
		rampUpMinimumHostsPct := v["ramp_up_minimum_hosts_percent"].(int)
		rampUpCapacityThresholdPct := v["ramp_up_capacity_threshold_percent"].(int)
		peakStartTime := v["peak_start_time"].(string)
		rampDownStartTime := v["ramp_down_start_time"].(string)
		rampDownMinimumHostsPct := v["ramp_down_minimum_hosts_percent"].(int)
		rampDownCapacityThresholdPct := v["ramp_down_capacity_threshold_percent"].(int)
		rampDownForceLogoffUsers := v["ramp_down_force_logoff_users"].(bool)
		rampDownWaitTimeMinutes := v["ramp_down_wait_time_minutes"].(int)
		rampDownNotificationMessage := v["ramp_down_notification_message"].(string)
		offPeakStartTime := v["off_peak_start_time"].(string)

		results = append(results, scalingplan.ScalingSchedule{
			Name:                           pointer.To(name),
			DaysOfWeek:                     &daysOfWeek,
			RampUpStartTime:                expandScalingPlanScheduleTime(rampUpStartTime),
			RampUpLoadBalancingAlgorithm:   pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["ramp_up_load_balancing_algorithm"].(string)),
			RampUpMinimumHostsPct:          pointer.To(int64(rampUpMinimumHostsPct)),
			RampUpCapacityThresholdPct:     pointer.To(int64(rampUpCapacityThresholdPct)),
			PeakStartTime:                  expandScalingPlanScheduleTime(peakStartTime),
			PeakLoadBalancingAlgorithm:     pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["peak_load_balancing_algorithm"].(string)),
			RampDownStartTime:              expandScalingPlanScheduleTime(rampDownStartTime),
			RampDownLoadBalancingAlgorithm: pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["ramp_down_load_balancing_algorithm"].(string)),
			RampDownMinimumHostsPct:        pointer.To(int64(rampDownMinimumHostsPct)),
			RampDownCapacityThresholdPct:   pointer.To(int64(rampDownCapacityThresholdPct)),
			RampDownForceLogoffUsers:       pointer.To(rampDownForceLogoffUsers),
			RampDownStopHostsWhen:          pointer.ToEnum[scalingplan.StopHostsWhen](v["ramp_down_stop_hosts_when"].(string)),
			RampDownWaitTimeMinutes:        pointer.To(int64(rampDownWaitTimeMinutes)),
			RampDownNotificationMessage:    pointer.To(rampDownNotificationMessage),
			OffPeakStartTime:               expandScalingPlanScheduleTime(offPeakStartTime),
			OffPeakLoadBalancingAlgorithm:  pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["off_peak_load_balancing_algorithm"].(string)),
		})
	}

	return &results
}

func expandScalingPlanScheduleTime(input string) *scalingplan.Time {
	if len(input) == 0 {
		return nil
	}

	time := strings.Split(input, ":")
	hour, _ := strconv.Atoi(time[0])
	minute, _ := strconv.Atoi(time[1])

	return &scalingplan.Time{
		Hour:   int64(hour),
		Minute: int64(minute),
	}
}

func expandScalingPlanHostpoolReference(input []any) *[]scalingplan.ScalingHostPoolReference {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingHostPoolReference, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)
		hostPoolArmPath := v["hostpool_id"].(string)
		scalingPlanEnabled := v["scaling_plan_enabled"].(bool)

		results = append(results, scalingplan.ScalingHostPoolReference{
			HostPoolArmPath:    pointer.To(hostPoolArmPath),
			ScalingPlanEnabled: pointer.To(scalingPlanEnabled),
		})
	}
	return &results
}
