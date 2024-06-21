// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hybridcompute

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/hybrid-compute"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Hybrid Compute"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Hybrid Compute",
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

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	dataSources := []sdk.DataSource{
		ArcMachineDataSource{},
	}

	if !features.FourPointOhBeta() {
		dataSources = append(dataSources, HybridComputeMachineDataSource{})
	}

	return dataSources
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ArcMachineExtensionResource{},
		ArcPrivateLinkScopeResource{},
	}
}
