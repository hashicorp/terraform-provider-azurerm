// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_maintenance_configuration

type ContainerAppEnvironmentMaintenanceConfigurationModel struct {
	ContainerAppEnvironmentId string              `tfschema:"container_app_environment_id"`
	MaintenanceWindows        []MaintenanceWindow `tfschema:"maintenance_window"`
}

type MaintenanceWindow struct {
	DayOfWeek      string `tfschema:"day_of_week"`
	StartHourInUtc int64  `tfschema:"start_hour_in_utc"`
	DurationHours  int64  `tfschema:"duration_hours"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentMaintenanceConfigurationModel{}
}
