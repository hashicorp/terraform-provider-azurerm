// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package networkmanager

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_admin_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_admin_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_connectivity_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_ipam_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_ipam_pool_static_cidr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_management_group_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_network_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_routing_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_routing_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_routing_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_scope_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_security_admin_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_static_member"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_subscription_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_verifier_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_manager_verifier_workspace_reachability_analysis_intent"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_security_perimeter"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_security_perimeter_access_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_security_perimeter_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/resources/network_security_perimeter_profile"
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
	return "Network Manager"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/network-manager"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		network_manager_connectivity_configuration.DataSource{},
		network_manager.DataSource{},
		network_manager_ipam_pool.DataSource{},
		network_manager_network_group.DataSource{},
		network_security_perimeter.DataSource{},
		network_security_perimeter_profile.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		network_manager_admin_rule_collection.Resource{},
		network_manager_admin_rule.Resource{},
		network_manager_connectivity_configuration.Resource{},
		network_manager_deployment.Resource{},
		network_manager_ipam_pool.Resource{},
		network_manager_ipam_pool_static_cidr.Resource{},
		network_manager_management_group_connection.Resource{},
		network_manager_network_group.Resource{},
		network_manager.Resource{},
		network_manager_routing_configuration.Resource{},
		network_manager_routing_rule_collection.Resource{},
		network_manager_routing_rule.Resource{},
		network_manager_scope_connection.Resource{},
		network_manager_security_admin_configuration.Resource{},
		network_manager_static_member.Resource{},
		network_manager_subscription_connection.Resource{},
		network_manager_verifier_workspace_reachability_analysis_intent.Resource{},
		network_manager_verifier_workspace.Resource{},
		network_security_perimeter_access_rule.Resource{},
		network_security_perimeter_association.Resource{},
		network_security_perimeter_profile.Resource{},
		network_security_perimeter.Resource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
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
