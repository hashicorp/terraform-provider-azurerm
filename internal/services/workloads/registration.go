// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/resources/workloads_sap_discovery_virtual_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/resources/workloads_sap_single_node_virtual_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/resources/workloads_sap_three_tier_virtual_instance"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/workloads"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Workloads"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Workloads",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		workloads_sap_discovery_virtual_instance.Resource{},
		workloads_sap_single_node_virtual_instance.Resource{},
		workloads_sap_three_tier_virtual_instance.Resource{},
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
		workloads_sap_discovery_virtual_instance.ListResource{},
		workloads_sap_single_node_virtual_instance.ListResource{},
		workloads_sap_three_tier_virtual_instance.ListResource{},
	}
}
