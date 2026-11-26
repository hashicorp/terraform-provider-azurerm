// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/resources/load_test"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/resources/playwright_workspace"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/load-test"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Load Test",
	}
}

func (r Registration) Name() string {
	return "LoadTestService"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		load_test.DataSource{},
		playwright_workspace.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		load_test.Resource{},
		playwright_workspace.Resource{},
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
		playwright_workspace.ListResource{},
	}
}
