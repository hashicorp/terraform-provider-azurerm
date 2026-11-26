// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment_dapr_component"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment_managed_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_environment_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/resources/container_app_job"
)

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/container-apps"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Container Apps",
	}
}

func (r Registration) Name() string {
	return "Container Apps"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		container_app.DataSource{},
		container_app_environment_certificate.DataSource{},
		container_app_environment.DataSource{},
		container_app_environment_storage.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		container_app_custom_domain.Resource{},
		container_app_environment_certificate.Resource{},
		container_app_environment_custom_domain.Resource{},
		container_app_environment_dapr_component.Resource{},
		container_app_environment_managed_certificate.Resource{},
		container_app_environment.Resource{},
		container_app_environment_storage.Resource{},
		container_app_job.Resource{},
		container_app.Resource{},
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
