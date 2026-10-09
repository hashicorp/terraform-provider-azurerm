// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package networkmanager

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
		ManagerConnectivityConfigurationDataSource{},
		ManagerDataSource{},
		ManagerIpamPoolDataSource{},
		ManagerNetworkGroupDataSource{},
		NetworkSecurityPerimeterDataSource{},
		NetworkSecurityPerimeterProfileDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagerAdminRuleCollectionResource{},
		ManagerAdminRuleResource{},
		ManagerConnectivityConfigurationResource{},
		ManagerDeploymentResource{},
		ManagerIpamPoolResource{},
		ManagerIpamPoolStaticCidrResource{},
		ManagerManagementGroupConnectionResource{},
		ManagerNetworkGroupResource{},
		ManagerResource{},
		ManagerRoutingConfigurationResource{},
		ManagerRoutingRuleCollectionResource{},
		ManagerRoutingRuleResource{},
		ManagerScopeConnectionResource{},
		ManagerSecurityAdminConfigurationResource{},
		ManagerStaticMemberResource{},
		ManagerSubscriptionConnectionResource{},
		ManagerVerifierWorkspaceReachabilityAnalysisIntentResource{},
		ManagerVerifierWorkspaceResource{},
		NetworkSecurityPerimeterAccessRuleResource{},
		NetworkSecurityPerimeterAssociationResource{},
		NetworkSecurityPerimeterProfileResource{},
		NetworkSecurityPerimeterResource{},
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
