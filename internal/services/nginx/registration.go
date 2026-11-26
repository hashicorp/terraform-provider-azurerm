// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/resources/nginx_api_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/resources/nginx_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/resources/nginx_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/resources/nginx_deployment"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/nginx"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Nginx"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"NGINX",
	}
}

// DataSources ...
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		nginx_api_key.DataSource{},
		nginx_certificate.DataSource{},
		nginx_configuration.DataSource{},
		nginx_deployment.DataSource{},
	}
}

// Resources ...
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		nginx_api_key.Resource{},
		nginx_certificate.Resource{},
		nginx_configuration.Resource{},
		nginx_deployment.Resource{},
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
