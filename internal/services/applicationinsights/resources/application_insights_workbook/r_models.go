// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

type ApplicationInsightsWorkbookModel struct {
	Name               string            `tfschema:"name"`
	ResourceGroupName  string            `tfschema:"resource_group_name"`
	Category           string            `tfschema:"category"`
	Description        string            `tfschema:"description"`
	DisplayName        string            `tfschema:"display_name"`
	Location           string            `tfschema:"location"`
	DataJson           string            `tfschema:"data_json"`
	SourceId           string            `tfschema:"source_id"`
	StorageContainerId string            `tfschema:"storage_container_id"`
	Tags               map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ApplicationInsightsWorkbookModel{}
}
