// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_application_insights

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceLoggerApplicationInsightsModel{}
}

type ApiManagementWorkspaceLoggerApplicationInsightsModel struct {
	Name                     string                     `tfschema:"name"`
	ApiManagementWorkspaceId string                     `tfschema:"api_management_workspace_id"`
	ApplicationInsights      []ApplicationInsightsModel `tfschema:"application_insights"`
	BufferingEnabled         bool                       `tfschema:"buffering_enabled"`
	Description              string                     `tfschema:"description"`
	ResourceId               string                     `tfschema:"resource_id"`
}

type ApplicationInsightsModel struct {
	InstrumentationKey string `tfschema:"instrumentation_key"`
	ConnectionString   string `tfschema:"connection_string"`
}
