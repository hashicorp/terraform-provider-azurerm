// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/storagemover"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Storage Mover"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Storage Mover",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		StorageMoverResource{},
		StorageMoverAgentResource{},
		StorageMoverSourceEndpointResource{},
		StorageMoverTargetEndpointResource{},
		StorageMoverProjectResource{},
		StorageMoverJobDefinitionResource{},
	}
}
