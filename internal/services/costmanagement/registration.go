// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/billing_account_cost_management_export"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/cost_anomaly_alert"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/cost_management_scheduled_action"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/resource_group_cost_management_export"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/resource_group_cost_management_view"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/subscription_cost_management_export"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/resources/subscription_cost_management_view"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/cost-management"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		cost_anomaly_alert.Resource{},
		billing_account_cost_management_export.Resource{},
		cost_management_scheduled_action.Resource{},
		resource_group_cost_management_export.Resource{},
		resource_group_cost_management_view.Resource{},
		subscription_cost_management_export.Resource{},
		subscription_cost_management_view.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Cost Management"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Cost Management",
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
