// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sreagents

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration     = Registration{}
	_ sdk.FrameworkServiceRegistration = Registration{}
)

func (Registration) Name() string {
	return "SRE Agents"
}

func (Registration) AssociatedGitHubLabel() string {
	return "service/sre-agents"
}

func (Registration) WebsiteCategories() []string {
	return []string{"SRE Agents"}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{SREAgentResource{}}
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{SREAgentListResource{}}
}

func (Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}
