// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_attached_network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_catalog"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_dev_box_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_environment_type"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_gallery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_network_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_project"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_project_environment_type"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/resources/dev_center_project_pool"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

type Registration struct{}

func (r Registration) Name() string {
	return "Dev Center"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dev-center"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Dev Center",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		dev_center_attached_network.DataSource{},
		dev_center_catalog.DataSource{},
		dev_center.DataSource{},
		dev_center_dev_box_definition.DataSource{},
		dev_center_environment_type.DataSource{},
		dev_center_gallery.DataSource{},
		dev_center_network_connection.DataSource{},
		dev_center_project.DataSource{},
		dev_center_project_environment_type.DataSource{},
		dev_center_project_pool.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		dev_center_attached_network.Resource{},
		dev_center_catalog.Resource{},
		dev_center_dev_box_definition.Resource{},
		dev_center_environment_type.Resource{},
		dev_center_gallery.Resource{},
		dev_center_network_connection.Resource{},
		dev_center_project_environment_type.Resource{},
		dev_center_project_pool.Resource{},
		dev_center_project.Resource{},
		dev_center.Resource{},
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
