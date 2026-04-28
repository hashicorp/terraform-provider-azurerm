// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package chaosstudio

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct {
	autoRegistration
}

var _ sdk.FrameworkServiceRegistration = Registration{}

// Name is the name of this Service
func (r Registration) Name() string {
	return r.autoRegistration.Name()
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return r.autoRegistration.WebsiteCategories()
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	dataSources := []sdk.DataSource{}
	dataSources = append(dataSources, r.autoRegistration.DataSources()...)
	return dataSources
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		ChaosStudioCapabilityResource{},
		ChaosStudioExperimentResource{},
	}
	resources = append(resources, r.autoRegistration.Resources()...)
	return resources
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
