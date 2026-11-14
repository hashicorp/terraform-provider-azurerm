// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dynatrace

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dynatrace/resources/dynatrace_monitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dynatrace/resources/dynatrace_tag_rules"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dynatrace"
}

func (r Registration) Name() string {
	return "Dynatrace"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		dynatrace_monitor.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		dynatrace_monitor.Resource{},
		dynatrace_tag_rules.Resource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Dynatrace",
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
