// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/client_config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/marketplace_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/pim_active_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/pim_eligible_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/role_assignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/role_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/resources/role_management_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/authorization"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Authorization"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Authorization",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_client_config": client_config.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_role_assignment": role_assignment.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		role_assignments.DataSource{},
		role_definition.DataSource{},
		role_management_policy.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		pim_active_role_assignment.Resource{},
		pim_eligible_role_assignment.Resource{},
		marketplace_role_assignment.Resource{},
		role_definition.Resource{},
		role_management_policy.Resource{},
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
