// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_dynatrace_application_performance_monitoring

type SpringCloudDynatraceApplicationPerformanceMonitoringModel struct {
	Name                 string `tfschema:"name"`
	SpringCloudServiceId string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled      bool   `tfschema:"globally_enabled"`
	ApiUrl               string `tfschema:"api_url"`
	ApiToken             string `tfschema:"api_token"`
	ConnectionPoint      string `tfschema:"connection_point"`
	EnvironmentId        string `tfschema:"environment_id"`
	Tenant               string `tfschema:"tenant"`
	TenantToken          string `tfschema:"tenant_token"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudDynatraceApplicationPerformanceMonitoringModel{}
}
