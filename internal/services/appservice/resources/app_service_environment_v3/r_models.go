// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

type ClusterSettingModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type AppServiceEnvironmentV3Model struct {
	Name                               string                            `tfschema:"name"`
	ResourceGroup                      string                            `tfschema:"resource_group_name"`
	SubnetId                           string                            `tfschema:"subnet_id"`
	AllowNewPrivateEndpointConnections bool                              `tfschema:"allow_new_private_endpoint_connections"`
	ClusterSetting                     []ClusterSettingModel             `tfschema:"cluster_setting"`
	DedicatedHostCount                 int64                             `tfschema:"dedicated_host_count"`
	InternalLoadBalancingMode          string                            `tfschema:"internal_load_balancing_mode"`
	RemoteDebuggingEnabled             bool                              `tfschema:"remote_debugging_enabled"`
	ZoneRedundant                      bool                              `tfschema:"zone_redundant"`
	Tags                               map[string]string                 `tfschema:"tags"`
	DnsSuffix                          string                            `tfschema:"dns_suffix"`
	ExternalInboundIPAddresses         []string                          `tfschema:"external_inbound_ip_addresses"`
	InboundNetworkDependencies         []AppServiceV3InboundDependencies `tfschema:"inbound_network_dependencies"`
	InternalInboundIPAddresses         []string                          `tfschema:"internal_inbound_ip_addresses"`
	IpSSLAddressCount                  int64                             `tfschema:"ip_ssl_address_count"`
	LinuxOutboundIPAddresses           []string                          `tfschema:"linux_outbound_ip_addresses"`
	Location                           string                            `tfschema:"location"`
	PricingTier                        string                            `tfschema:"pricing_tier"`
	WindowsOutboundIPAddresses         []string                          `tfschema:"windows_outbound_ip_addresses"`
}

func (r Resource) ModelObject() any {
	return &AppServiceEnvironmentV3Model{}
}
