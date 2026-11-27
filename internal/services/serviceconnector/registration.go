// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/resources/app_service_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/resources/function_app_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/resources/spring_cloud_connection"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/service-connector"
}

func (r Registration) WebsiteCategories() []string {
	return []string{}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		app_service_connection.Resource{},
		function_app_connection.Resource{},
		spring_cloud_connection.Resource{},
	}
}

func (r Registration) Name() string {
	return "ServiceConnector"
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
