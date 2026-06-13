package synapse

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sqlPoolMaintenanceWindowResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"day_of_week": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(synapse.DayOfWeekSunday),
					string(synapse.DayOfWeekTuesday),
					string(synapse.DayOfWeekWednesday),
					string(synapse.DayOfWeekThursday),
					string(synapse.DayOfWeekSaturday),
				}, false),
			},
			"start_time_utc": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					v := i.(string)
					_, err := time.Parse(time.TimeOnly, v)
					if err != nil {
						return nil, []error{fmt.Errorf("expected `start_time_utc` to be in the format HH:MM:SS - got: %q", v)}
					}
					return nil, nil
				},
			},
			"duration_in_hours": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(3, 8),
			},
		},
	}
}

func validateMaintenanceSchedule(d *pluginsdk.ResourceDiff) error {
	maintenanceScheduleRaw := d.Get("maintenance_schedule").([]interface{})
	if len(maintenanceScheduleRaw) == 0 || maintenanceScheduleRaw[0] == nil {
		return nil
	}

	maintenanceSchedule := maintenanceScheduleRaw[0].(map[string]interface{})

	primaryWindowRaw := maintenanceSchedule["primary_maintenance_window"].([]interface{})
	secondaryWindowRaw := maintenanceSchedule["secondary_maintenance_window"].([]interface{})

	primaryDay := primaryWindowRaw[0].(map[string]interface{})["day_of_week"].(string)
	secondaryDay := secondaryWindowRaw[0].(map[string]interface{})["day_of_week"].(string)

	isTueThuRange := func(day string) bool {
		return day == string(synapse.DayOfWeekTuesday) || day == string(synapse.DayOfWeekWednesday) || day == string(synapse.DayOfWeekThursday)
	}

	isSatSunRange := func(day string) bool {
		return day == string(synapse.DayOfWeekSaturday) || day == string(synapse.DayOfWeekSunday)
	}

	primaryInTueThu, secondaryInTueThu := isTueThuRange(primaryDay), isTueThuRange(secondaryDay)
	primaryInSatSun, secondaryInSatSun := isSatSunRange(primaryDay), isSatSunRange(secondaryDay)

	// Validate they are in different ranges
	if (primaryInTueThu && secondaryInTueThu) || (primaryInSatSun && secondaryInSatSun) {
		return fmt.Errorf("`primary_maintenance_window` and `secondary_maintenance_window` must be in different date ranges: Saturday-Sunday and Tuesday-Thursday")
	}

	return nil
}

func expandSqlPoolMaintenanceSchedule(input []interface{}) *synapse.MaintenanceWindows {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	windows := make([]synapse.MaintenanceWindowTimeRange, 2)

	primaryWindow := v["primary_maintenance_window"].([]interface{})[0]
	windows[0] = expandSqlPoolMaintenanceWindow(primaryWindow)

	secondaryWindow := v["secondary_maintenance_window"].([]interface{})[0]
	windows[1] = expandSqlPoolMaintenanceWindow(secondaryWindow)

	maintenanceWindow := &synapse.MaintenanceWindows{
		MaintenanceWindowsProperties: &synapse.MaintenanceWindowsProperties{
			TimeRanges: &windows,
		},
	}
	return maintenanceWindow
}

func flattenSqlPoolMaintenanceSchedule(maintenanceWindows synapse.MaintenanceWindows) ([]interface{}, error) {
	if maintenanceWindows.MaintenanceWindowsProperties == nil || maintenanceWindows.MaintenanceWindowsProperties.TimeRanges == nil {
		return nil, fmt.Errorf("expected `maintenance_windows_properties.time_ranges` to be populated")
	}

	windows := *maintenanceWindows.MaintenanceWindowsProperties.TimeRanges
	flattenedPrimary, err := flattenSqlPoolMaintenanceWindow(windows[0])
	if err != nil {
		return nil, err
	}
	flattenedSecondary, err := flattenSqlPoolMaintenanceWindow(windows[1])
	if err != nil {
		return nil, err
	}

	result := []interface{}{
		map[string]interface{}{
			"primary_maintenance_window":   flattenedPrimary,
			"secondary_maintenance_window": flattenedSecondary,
		},
	}
	return result, nil
}

func expandSqlPoolMaintenanceWindow(input interface{}) synapse.MaintenanceWindowTimeRange {
	window := input.(map[string]interface{})
	durationHours := window["duration_in_hours"].(int)
	durationISO := fmt.Sprintf("PT%dM", durationHours*60) // API only allows minutes

	windowTimeRange := synapse.MaintenanceWindowTimeRange{
		DayOfWeek: synapse.DayOfWeek(window["day_of_week"].(string)),
		StartTime: pointer.To(window["start_time_utc"].(string)),
		Duration:  pointer.To(durationISO),
	}

	return windowTimeRange
}

func flattenSqlPoolMaintenanceWindow(timeRanges synapse.MaintenanceWindowTimeRange) ([]interface{}, error) {
	result := make([]interface{}, 1)

	var durationHours int
	if timeRanges.Duration != nil {
		hours, err := convertISO8601MinuteDurationToHours(pointer.From(timeRanges.Duration))
		if err != nil {
			return nil, err
		}
		durationHours = hours
	}

	result[0] = map[string]interface{}{
		"day_of_week":       timeRanges.DayOfWeek,
		"start_time_utc":    pointer.From(timeRanges.StartTime),
		"duration_in_hours": durationHours,
	}
	return result, nil
}

func convertISO8601MinuteDurationToHours(iso8601Duration string) (int, error) {
	if !regexp.MustCompile(`^PT\d+M$`).MatchString(iso8601Duration) {
		return 0, fmt.Errorf("invalid ISO8601 minute format %q, expected PT###M", iso8601Duration)
	}

	var minutes int
	if _, err := fmt.Sscanf(iso8601Duration, "PT%dM", &minutes); err != nil {
		return 0, fmt.Errorf("parsing minutes from %q: %+v", iso8601Duration, err)
	}

	hours := minutes / 60
	return hours, nil
}
