// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_action_custom"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_action_http"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_agreement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_assembly"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_batch_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_map"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_partner"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_integration_account_session"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_standard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_trigger_custom"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_trigger_http_request"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_trigger_recurrence"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/resources/logic_app_workflow"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		logic_app_standard.Resource{},
	}
}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/logic"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Logic"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Logic App",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_logic_app_integration_account": logic_app_integration_account.RegisterDataSource(),
		"azurerm_logic_app_standard":            logic_app_standard.RegisterDataSource(),
		"azurerm_logic_app_workflow":            logic_app_workflow.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_logic_app_action_custom":                           logic_app_action_custom.RegisterResource(),
		"azurerm_logic_app_action_http":                             logic_app_action_http.RegisterResource(),
		"azurerm_logic_app_integration_account":                     logic_app_integration_account.RegisterResource(),
		"azurerm_logic_app_integration_account_agreement":           logic_app_integration_account_agreement.RegisterResource(),
		"azurerm_logic_app_integration_account_assembly":            logic_app_integration_account_assembly.RegisterResource(),
		"azurerm_logic_app_integration_account_batch_configuration": logic_app_integration_account_batch_configuration.RegisterResource(),
		"azurerm_logic_app_integration_account_certificate":         logic_app_integration_account_certificate.RegisterResource(),
		"azurerm_logic_app_integration_account_map":                 logic_app_integration_account_map.RegisterResource(),
		"azurerm_logic_app_integration_account_partner":             logic_app_integration_account_partner.RegisterResource(),
		"azurerm_logic_app_integration_account_schema":              logic_app_integration_account_schema.RegisterResource(),
		"azurerm_logic_app_integration_account_session":             logic_app_integration_account_session.RegisterResource(),
		"azurerm_logic_app_trigger_custom":                          logic_app_trigger_custom.RegisterResource(),
		"azurerm_logic_app_trigger_http_request":                    logic_app_trigger_http_request.RegisterResource(),
		"azurerm_logic_app_trigger_recurrence":                      logic_app_trigger_recurrence.RegisterResource(),
		"azurerm_logic_app_workflow":                                logic_app_workflow.RegisterResource(),
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
