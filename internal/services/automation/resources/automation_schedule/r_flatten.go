// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_schedule

import (
	// import time/tzdata to embed timezone information in the program
	// add this to resolve https://github.com/hashicorp/terraform-provider-azurerm/issues/20690
	_ "time/tzdata"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

func flattenArmAutomationScheduleAdvancedWeekDays(s *schedule.AdvancedSchedule) *pluginsdk.Set {
	flattenedWeekDays := pluginsdk.NewSet(set.HashStringIgnoreCase, []any{})
	if weekDays := s.WeekDays; weekDays != nil {
		for _, v := range *weekDays {
			flattenedWeekDays.Add(v)
		}
	}
	return flattenedWeekDays
}

func flattenArmAutomationScheduleAdvancedMonthDays(s *schedule.AdvancedSchedule) *pluginsdk.Set {
	flattenedMonthDays := pluginsdk.NewSet(set.HashInt, []any{})
	if monthDays := s.MonthDays; monthDays != nil {
		for _, v := range *monthDays {
			flattenedMonthDays.Add(int(v))
		}
	}
	return flattenedMonthDays
}

func flattenArmAutomationScheduleAdvancedMonthlyOccurrences(s *schedule.AdvancedSchedule) []map[string]any {
	flattenedMonthlyOccurrences := make([]map[string]any, 0)
	if monthlyOccurrences := s.MonthlyOccurrences; monthlyOccurrences != nil {
		for _, v := range *monthlyOccurrences {
			f := make(map[string]any)
			f["day"] = v.Day
			f["occurrence"] = int(*v.Occurrence)
			flattenedMonthlyOccurrences = append(flattenedMonthlyOccurrences, f)
		}
	}
	return flattenedMonthlyOccurrences
}
