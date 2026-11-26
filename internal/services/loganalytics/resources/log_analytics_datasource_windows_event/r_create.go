// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_event

func expandLogAnalyticsDataSourceWindowsEventEventType(input []any) []dataSourceWindowsEventEventType {
	output := []dataSourceWindowsEventEventType{}
	for _, eventType := range input {
		output = append(output, dataSourceWindowsEventEventType{eventType.(string)})
	}
	return output
}
