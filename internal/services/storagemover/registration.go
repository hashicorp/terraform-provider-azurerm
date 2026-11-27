// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover_agent"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover_job_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover_project"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover_source_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/resources/storage_mover_target_endpoint"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
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
		storage_mover_agent.Resource{},
		storage_mover_job_definition.Resource{},
		storage_mover_project.Resource{},
		storage_mover.Resource{},
		storage_mover_source_endpoint.Resource{},
		storage_mover_target_endpoint.Resource{},
	}
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
	return []sdk.FrameworkListWrappedResource{
		storage_mover_agent.ListResource{},
		storage_mover_job_definition.ListResource{},
		storage_mover.ListResource{},
		storage_mover_project.ListResource{},
		storage_mover_source_endpoint.ListResource{},
		storage_mover_target_endpoint.ListResource{},
	}
}
