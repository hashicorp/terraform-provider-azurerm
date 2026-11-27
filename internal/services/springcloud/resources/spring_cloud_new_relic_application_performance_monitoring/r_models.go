// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_new_relic_application_performance_monitoring

type SpringCloudNewRelicApplicationPerformanceMonitoringModel struct {
	Name                         string            `tfschema:"name"`
	SpringCloudServiceId         string            `tfschema:"spring_cloud_service_id"`
	GloballyEnabled              bool              `tfschema:"globally_enabled"`
	AppName                      string            `tfschema:"app_name"`
	AgentEnabled                 bool              `tfschema:"agent_enabled"`
	AppServerPort                int64             `tfschema:"app_server_port"`
	AuditModeEnabled             bool              `tfschema:"audit_mode_enabled"`
	AutoAppNamingEnabled         bool              `tfschema:"auto_app_naming_enabled"`
	AutoTransactionNamingEnabled bool              `tfschema:"auto_transaction_naming_enabled"`
	CustomTracingEnabled         bool              `tfschema:"custom_tracing_enabled"`
	Labels                       map[string]string `tfschema:"labels"`
	LicenseKey                   string            `tfschema:"license_key"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudNewRelicApplicationPerformanceMonitoringModel{}
}
