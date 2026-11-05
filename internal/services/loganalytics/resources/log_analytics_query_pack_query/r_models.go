// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack_query

type LogAnalyticsQueryPackQueryModel struct {
	Name                   string            `tfschema:"name"`
	QueryPackId            string            `tfschema:"query_pack_id"`
	Body                   string            `tfschema:"body"`
	DisplayName            string            `tfschema:"display_name"`
	Categories             []string          `tfschema:"categories"`
	Description            string            `tfschema:"description"`
	AdditionalSettingsJson string            `tfschema:"additional_settings_json"`
	ResourceTypes          []string          `tfschema:"resource_types"`
	Solutions              []string          `tfschema:"solutions"`
	Tags                   map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &LogAnalyticsQueryPackQueryModel{}
}
