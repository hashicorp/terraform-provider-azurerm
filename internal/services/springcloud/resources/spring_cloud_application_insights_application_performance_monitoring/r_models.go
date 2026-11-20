// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_application_insights_application_performance_monitoring

type SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel struct {
	Name                      string `tfschema:"name"`
	SpringCloudServiceId      string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled           bool   `tfschema:"globally_enabled"`
	ConnectionString          string `tfschema:"connection_string"`
	RoleName                  string `tfschema:"role_name"`
	RoleInstance              string `tfschema:"role_instance"`
	SamplingRequestsPerSecond int64  `tfschema:"sampling_requests_per_second"`
	SamplingPercentage        int64  `tfschema:"sampling_percentage"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel{}
}
