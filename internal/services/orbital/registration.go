// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orbital

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/orbital"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Orbital",
	}
}

func (r Registration) Name() string {
	return "Orbital"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	if !features.FivePointOh() {
		return []sdk.Resource{
			SpacecraftResource{},
			ContactProfileResource{},
			ContactResource{},
		}
	}

	return nil
}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

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
