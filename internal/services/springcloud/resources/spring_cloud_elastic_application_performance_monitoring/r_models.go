// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_elastic_application_performance_monitoring

type SpringCloudElasticApplicationPerformanceMonitoringModel struct {
	Name                 string   `tfschema:"name"`
	SpringCloudServiceId string   `tfschema:"spring_cloud_service_id"`
	GloballyEnabled      bool     `tfschema:"globally_enabled"`
	ServiceName          string   `tfschema:"service_name"`
	ApplicationPackages  []string `tfschema:"application_packages"`
	ServerUrl            string   `tfschema:"server_url"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudElasticApplicationPerformanceMonitoringModel{}
}
