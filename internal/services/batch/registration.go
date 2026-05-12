// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/batch"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Batch"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Batch",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	datasources := map[string]*pluginsdk.Resource{
		"azurerm_batch_account":     dataSourceBatchAccount(),
		"azurerm_batch_application": dataSourceBatchApplication(),
		"azurerm_batch_pool":        dataSourceBatchPool(),
	}

	// The batch certificate feature was retired by Azure on 2024-02-29
	if !features.FivePointOh() {
		datasources["azurerm_batch_certificate"] = dataSourceBatchCertificate()
	}

	return datasources
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_batch_account":     resourceBatchAccount(),
		"azurerm_batch_application": resourceBatchApplication(),
		"azurerm_batch_pool":        resourceBatchPool(),
	}

	// The batch certificate feature was retired by Azure on 2024-02-29
	if !features.FivePointOh() {
		resources["azurerm_batch_certificate"] = resourceBatchCertificate()
	}

	return resources
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		BatchJobResource{},
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
	return []sdk.FrameworkListWrappedResource{}
}
