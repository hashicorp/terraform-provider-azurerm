// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package manageddevopspools

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/resources/managed_devops_pool"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		managed_devops_pool.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		managed_devops_pool.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Managed DevOps Pools"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Managed DevOps Pools",
	}
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/managed-devops-pools"
}
