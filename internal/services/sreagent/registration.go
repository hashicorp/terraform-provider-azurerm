// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration             = Registration{}
)

type Registration struct{}

func (Registration) Name() string {
	return "SRE Agent"
}

func (Registration) AssociatedGitHubLabel() string {
	return "service/sre-agent"
}

func (Registration) WebsiteCategories() []string {
	return []string{"SRE Agent"}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{SreAgentResource{}}
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{SreAgentListResource{}}
}

func (Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}
