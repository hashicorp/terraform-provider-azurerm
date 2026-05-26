// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extendedlocation

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/extended-location"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Extended Location",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ExtendedLocationCustomLocationDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	if !features.FivePointOh() {
		return []sdk.Resource{
			CustomLocationResource{}, // this resource is renamed and should be removed in 5.0
			ExtendedLocationCustomLocationResource{},
		}
	}

	return []sdk.Resource{
		ExtendedLocationCustomLocationResource{},
	}
}

func (r Registration) Name() string {
	return "ExtendedLocation"
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
