// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/application_security_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/bastion_host"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/custom_ip_prefix"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/ip_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/ip_group_cidr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/ip_groups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/nat_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/nat_gateway_public_ip_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/nat_gateway_public_ip_prefix_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_connection_monitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_ddos_protection_plan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_application_gateway_backend_address_pool_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_application_security_group_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_backend_address_pool_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_nat_rule_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_security_group_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_profile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_security_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_security_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_service_tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_watcher_flow_log"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/private_endpoint_application_security_group_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/private_endpoint_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/private_link_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/private_link_service_endpoint_connections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/public_ip"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/public_ip_prefix"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/public_ips"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/route"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/route_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet_nat_gateway_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet_network_security_group_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet_route_table_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/subnet_service_endpoint_storage_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/virtual_machine_packet_capture"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/virtual_machine_scale_set_packet_capture"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/virtual_network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/virtual_network_dns_servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/virtual_network_peering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Network"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/network"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		virtual_network_peering.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		custom_ip_prefix.Resource{},
		private_endpoint_application_security_group_association.Resource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_security_group":                application_security_group.RegisterDataSource(),
		"azurerm_bastion_host":                              bastion_host.RegisterDataSource(),
		"azurerm_ip_group":                                  ip_group.RegisterDataSource(),
		"azurerm_ip_groups":                                 ip_groups.RegisterDataSource(),
		"azurerm_nat_gateway":                               nat_gateway.RegisterDataSource(),
		"azurerm_network_ddos_protection_plan":              network_ddos_protection_plan.RegisterDataSource(),
		"azurerm_network_interface":                         network_interface.RegisterDataSource(),
		"azurerm_network_security_group":                    network_security_group.RegisterDataSource(),
		"azurerm_network_service_tags":                      network_service_tags.RegisterDataSource(),
		"azurerm_network_watcher":                           network_watcher.RegisterDataSource(),
		"azurerm_private_endpoint_connection":               private_endpoint_connection.RegisterDataSource(),
		"azurerm_private_link_service":                      private_link_service.RegisterDataSource(),
		"azurerm_private_link_service_endpoint_connections": private_link_service_endpoint_connections.RegisterDataSource(),
		"azurerm_public_ip":                                 public_ip.RegisterDataSource(),
		"azurerm_public_ip_prefix":                          public_ip_prefix.RegisterDataSource(),
		"azurerm_public_ips":                                public_ips.RegisterDataSource(),
		"azurerm_route_table":                               route_table.RegisterDataSource(),
		"azurerm_subnet":                                    subnet.RegisterDataSource(),
		"azurerm_virtual_network":                           virtual_network.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_security_group":               application_security_group.RegisterResource(),
		"azurerm_bastion_host":                             bastion_host.RegisterResource(),
		"azurerm_ip_group":                                 ip_group.RegisterResource(),
		"azurerm_ip_group_cidr":                            ip_group_cidr.RegisterResource(),
		"azurerm_nat_gateway":                              nat_gateway.RegisterResource(),
		"azurerm_nat_gateway_public_ip_association":        nat_gateway_public_ip_association.RegisterResource(),
		"azurerm_nat_gateway_public_ip_prefix_association": nat_gateway_public_ip_prefix_association.RegisterResource(),
		"azurerm_network_connection_monitor":               network_connection_monitor.RegisterResource(),
		"azurerm_network_ddos_protection_plan":             network_ddos_protection_plan.RegisterResource(),
		"azurerm_network_interface":                        network_interface.RegisterResource(),

		"azurerm_network_interface_application_gateway_backend_address_pool_association": network_interface_application_gateway_backend_address_pool_association.RegisterResource(),
		"azurerm_network_interface_application_security_group_association":               network_interface_application_security_group_association.RegisterResource(),
		"azurerm_network_interface_backend_address_pool_association":                     network_interface_backend_address_pool_association.RegisterResource(),
		"azurerm_network_interface_nat_rule_association":                                 network_interface_nat_rule_association.RegisterResource(),
		"azurerm_network_interface_security_group_association":                           network_interface_security_group_association.RegisterResource(),

		"azurerm_network_profile":                           network_profile.RegisterResource(),
		"azurerm_network_security_group":                    network_security_group.RegisterResource(),
		"azurerm_network_security_rule":                     network_security_rule.RegisterResource(),
		"azurerm_network_watcher":                           network_watcher.RegisterResource(),
		"azurerm_network_watcher_flow_log":                  network_watcher_flow_log.RegisterResource(),
		"azurerm_private_endpoint":                          private_endpoint.RegisterResource(),
		"azurerm_private_link_service":                      private_link_service.RegisterResource(),
		"azurerm_public_ip":                                 public_ip.RegisterResource(),
		"azurerm_public_ip_prefix":                          public_ip_prefix.RegisterResource(),
		"azurerm_route":                                     route.RegisterResource(),
		"azurerm_route_table":                               route_table.RegisterResource(),
		"azurerm_subnet":                                    subnet.RegisterResource(),
		"azurerm_subnet_nat_gateway_association":            subnet_nat_gateway_association.RegisterResource(),
		"azurerm_subnet_network_security_group_association": subnet_network_security_group_association.RegisterResource(),
		"azurerm_subnet_route_table_association":            subnet_route_table_association.RegisterResource(),
		"azurerm_subnet_service_endpoint_storage_policy":    subnet_service_endpoint_storage_policy.RegisterResource(),
		"azurerm_virtual_machine_packet_capture":            virtual_machine_packet_capture.RegisterResource(),
		"azurerm_virtual_machine_scale_set_packet_capture":  virtual_machine_scale_set_packet_capture.RegisterResource(),
		"azurerm_virtual_network":                           virtual_network.RegisterResource(),
		"azurerm_virtual_network_dns_servers":               virtual_network_dns_servers.RegisterResource(),
		"azurerm_virtual_network_peering":                   virtual_network_peering.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		application_security_group.ListResource{},
		ip_group.ListResource{},
		nat_gateway.ListResource{},
		network_ddos_protection_plan.ListResource{},
		network_interface.ListResource{},
		network_profile.ListResource{},
		network_security_group.ListResource{},
		network_security_rule.ListResource{},
		network_watcher_flow_log.ListResource{},
		private_endpoint.ListResource{},
		public_ip.ListResource{},
		route.ListResource{},
		route_table.ListResource{},
		subnet.ListResource{},
		virtual_network.ListResource{},
		virtual_network_peering.ListResource{},
	}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}
