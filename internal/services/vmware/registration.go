// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/resources/vmware_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/resources/vmware_express_route_authorization"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/resources/vmware_netapp_volume_attachment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/resources/vmware_private_cloud"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/vmware"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "VMware"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Azure VMware Solution",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_vmware_private_cloud": vmware_private_cloud.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_vmware_cluster":                     vmware_cluster.RegisterResource(),
		"azurerm_vmware_express_route_authorization": vmware_express_route_authorization.RegisterResource(),
		"azurerm_vmware_private_cloud":               vmware_private_cloud.RegisterResource(),
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		vmware_netapp_volume_attachment.Resource{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
