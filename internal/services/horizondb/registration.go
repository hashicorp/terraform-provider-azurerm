// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (Registration) AssociatedGitHubLabel() string {
	return "service/horizondb"
}

func (Registration) Name() string {
	return "HorizonDB"
}

func (Registration) WebsiteCategories() []string {
	return []string{"HorizonDB"}
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		HorizondbClusterResource{},
		HorizondbFirewallRuleResource{},
	}
}

func (Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		HorizondbClusterListResource{},
		HorizondbFirewallRuleListResource{},
	}
}
