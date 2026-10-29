// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ContainerAppEnvironmentModel struct {
	Name                                    string                                     `tfschema:"name"`
	ResourceGroup                           string                                     `tfschema:"resource_group_name"`
	Location                                string                                     `tfschema:"location"`
	DaprApplicationInsightsConnectionString string                                     `tfschema:"dapr_application_insights_connection_string"`
	LogAnalyticsWorkspaceId                 string                                     `tfschema:"log_analytics_workspace_id"`
	LogsDestination                         string                                     `tfschema:"logs_destination"`
	InfrastructureSubnetId                  string                                     `tfschema:"infrastructure_subnet_id"`
	InternalLoadBalancerEnabled             bool                                       `tfschema:"internal_load_balancer_enabled"`
	Identity                                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess                     string                                     `tfschema:"public_network_access"`
	ZoneRedundant                           bool                                       `tfschema:"zone_redundancy_enabled"`
	Tags                                    map[string]any                             `tfschema:"tags"`
	WorkloadProfiles                        []WorkloadProfileModel                     `tfschema:"workload_profile"`
	InfrastructureResourceGroup             string                                     `tfschema:"infrastructure_resource_group_name"`
	Mtls                                    bool                                       `tfschema:"mutual_tls_enabled"`

	CustomDomainVerificationId string `tfschema:"custom_domain_verification_id"`

	DefaultDomain         string `tfschema:"default_domain"`
	DockerBridgeCidr      string `tfschema:"docker_bridge_cidr"`
	PlatformReservedCidr  string `tfschema:"platform_reserved_cidr"`
	PlatformReservedDnsIP string `tfschema:"platform_reserved_dns_ip_address"`
	StaticIP              string `tfschema:"static_ip_address"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentModel{}
}
