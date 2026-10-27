// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
)

func expandAgentUpdateSchedule(input []any) *[]hostpool.MaintenanceWindowProperties {
	if len(input) == 0 {
		return nil
	}

	results := make([]hostpool.MaintenanceWindowProperties, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)

		hourOfDay := int64(v["hour_of_day"].(int))

		results = append(results, hostpool.MaintenanceWindowProperties{
			DayOfWeek: pointer.ToEnum[hostpool.DayOfWeek](v["day_of_week"].(string)),
			Hour:      pointer.To(hourOfDay),
		})
	}

	return &results
}

func expandAgentUpdatePatch(input []any) *hostpool.AgentUpdatePatchProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	props := hostpool.AgentUpdatePatchProperties{}

	useSessionHostLocalTime := *pointer.To(raw["use_session_host_timezone"].(bool))
	props.MaintenanceWindowTimeZone = pointer.To(raw["timezone"].(string))
	props.UseSessionHostLocalTime = &useSessionHostLocalTime

	if raw["enabled"].(bool) {
		props.Type = pointer.To(hostpool.SessionHostComponentUpdateTypeScheduled)
		props.MaintenanceWindows = expandAgentUpdateSchedulePatch(raw["schedule"].([]any))
	} else {
		props.Type = pointer.To(hostpool.SessionHostComponentUpdateTypeDefault)
		props.MaintenanceWindows = &[]hostpool.MaintenanceWindowPatchProperties{}
	}

	return &props
}

func expandAgentUpdateSchedulePatch(input []any) *[]hostpool.MaintenanceWindowPatchProperties {
	if len(input) == 0 {
		return nil
	}

	results := make([]hostpool.MaintenanceWindowPatchProperties, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]any)

		hourOfDay := int64(v["hour_of_day"].(int))

		results = append(results, hostpool.MaintenanceWindowPatchProperties{
			DayOfWeek: pointer.ToEnum[hostpool.DayOfWeek](v["day_of_week"].(string)),
			Hour:      pointer.To(hourOfDay),
		})
	}

	return &results
}
