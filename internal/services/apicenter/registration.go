// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package apicenter

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_api_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_api_version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_environment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_metadata_schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/resources/api_center_service"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/api-center"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "API Center"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"API Center",
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		api_center_api.Resource{},
		api_center_api_definition.Resource{},
		api_center_api_version.Resource{},
		api_center_deployment.Resource{},
		api_center_environment.Resource{},
		api_center_metadata_schema.Resource{},
		api_center_service.Resource{},
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
		api_center_api.ListResource{},
		api_center_api_definition.ListResource{},
		api_center_api_version.ListResource{},
		api_center_deployment.ListResource{},
		api_center_environment.ListResource{},
		api_center_metadata_schema.ListResource{},
		api_center_service.ListResource{},
	}
}
