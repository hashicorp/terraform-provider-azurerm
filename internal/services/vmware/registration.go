// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vmware

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
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
		"VMware (AVS)",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_vmware_private_cloud": dataSourceVmwarePrivateCloud(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_vmware_private_cloud":               resourceVmwarePrivateCloud(),
		"azurerm_vmware_cluster":                     resourceVmwareCluster(),
		"azurerm_vmware_express_route_authorization": resourceVmwareExpressRouteAuthorization(),
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		NetappFileVolumeAttachmentResource{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}
