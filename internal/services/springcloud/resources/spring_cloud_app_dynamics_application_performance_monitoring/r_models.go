// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_dynamics_application_performance_monitoring

type SpringCloudAppDynamicsApplicationPerformanceMonitoringModel struct {
	Name                  string `tfschema:"name"`
	SpringCloudServiceId  string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled       bool   `tfschema:"globally_enabled"`
	AgentApplicationName  string `tfschema:"agent_application_name"`
	AgentTierName         string `tfschema:"agent_tier_name"`
	AgentNodeName         string `tfschema:"agent_node_name"`
	AgentUniqueHostId     string `tfschema:"agent_unique_host_id"`
	ControllerHostName    string `tfschema:"controller_host_name"`
	ControllerSslEnabled  bool   `tfschema:"controller_ssl_enabled"`
	ControllerPort        int64  `tfschema:"controller_port"`
	AgentAccountName      string `tfschema:"agent_account_name"`
	AgentAccountAccessKey string `tfschema:"agent_account_access_key"`
}

func (s Resource) ModelObject() any {
	return &SpringCloudAppDynamicsApplicationPerformanceMonitoringModel{}
}
