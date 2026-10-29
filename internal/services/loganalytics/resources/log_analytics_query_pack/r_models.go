// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack

type LogAnalyticsQueryPackModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &LogAnalyticsQueryPackModel{}
}
