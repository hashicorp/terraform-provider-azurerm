// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/mobile-network"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Mobile Network"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Mobile Network",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		DataNetworkDataSource{},
		MobileNetworkDataSource{},
		ServiceDataSource{},
		SiteDataSource{},
		SimGroupDataSource{},
		SliceDataSource{},
		SimPolicyDataSource{},
		PacketCoreControlPlaneDataSource{},
		PacketCoreDataPlaneDataSource{},
		AttachedDataNetworkDataSource{},
		SimDataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AttachedDataNetworkResource{},
		DataNetworkResource{},
		MobileNetworkResource{},
		PacketCoreControlPlaneResource{},
		PacketCoreDataPlaneResource{},
		SiteResource{},
		SliceResource{},
		ServiceResource{},
		SimGroupResource{},
		SimPolicyResource{},
		SimResource{},
	}
}
