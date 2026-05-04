// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
		DevCenterAttachedNetworkDataSource{},
		DevCenterCatalogDataSource{},
		DevCenterDataSource{},
		DevCenterDevBoxDefinitionDataSource{},
		DevCenterEnvironmentTypeDataSource{},
		DevCenterGalleryDataSource{},
		DevCenterNetworkConnectionDataSource{},
		DevCenterProjectDataSource{},
		DevCenterProjectEnvironmentTypeDataSource{},
		DevCenterProjectPoolDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DevCenterAttachedNetworkResource{},
		DevCenterGalleryResource{},
		DevCenterCatalogsResource{},
		DevCenterDevBoxDefinitionResource{},
		DevCenterEnvironmentTypeResource{},
		DevCenterNetworkConnectionResource{},
		DevCenterProjectPoolResource{},
		DevCenterProjectResource{},
		DevCenterProjectEnvironmentTypeResource{},
		DevCenterResource{},
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
