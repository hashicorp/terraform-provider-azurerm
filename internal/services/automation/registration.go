// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_connection_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_connection_classic_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_connection_service_principal"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_connection_type"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_credential"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_dsc_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_dsc_nodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_hybrid_runbook_worker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_hybrid_runbook_worker_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_job_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_powershell72_module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_python3_package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_runbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_runtime_environment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_runtime_environment_package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_source_control"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variable_bool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variable_datetime"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variable_int"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variable_object"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variable_string"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_variables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/resources/automation_webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		automation_runbook.DataSource{},
		automation_variables.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		automation_connection_type.Resource{},
		automation_runtime_environment_package.Resource{},
		automation_runtime_environment.Resource{},
		automation_hybrid_runbook_worker_group.Resource{},
		automation_hybrid_runbook_worker.Resource{},
		automation_powershell72_module.Resource{},
		automation_python3_package.Resource{},
		automation_source_control.Resource{},
		automation_watcher.Resource{},
	}
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/automation"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Automation"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Automation",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_automation_account":           automation_account.RegisterDataSource(),
		"azurerm_automation_variable_bool":     automation_variable_bool.RegisterDataSource(),
		"azurerm_automation_variable_datetime": automation_variable_datetime.RegisterDataSource(),
		"azurerm_automation_variable_int":      automation_variable_int.RegisterDataSource(),
		"azurerm_automation_variable_object":   automation_variable_object.RegisterDataSource(),
		"azurerm_automation_variable_string":   automation_variable_string.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_automation_account":                        automation_account.RegisterResource(),
		"azurerm_automation_certificate":                    automation_certificate.RegisterResource(),
		"azurerm_automation_connection":                     automation_connection.RegisterResource(),
		"azurerm_automation_connection_certificate":         automation_connection_certificate.RegisterResource(),
		"azurerm_automation_connection_classic_certificate": automation_connection_classic_certificate.RegisterResource(),
		"azurerm_automation_connection_service_principal":   automation_connection_service_principal.RegisterResource(),
		"azurerm_automation_credential":                     automation_credential.RegisterResource(),
		"azurerm_automation_dsc_configuration":              automation_dsc_configuration.RegisterResource(),
		"azurerm_automation_dsc_nodeconfiguration":          automation_dsc_nodeconfiguration.RegisterResource(),
		"azurerm_automation_job_schedule":                   automation_job_schedule.RegisterResource(),
		"azurerm_automation_module":                         automation_module.RegisterResource(),
		"azurerm_automation_runbook":                        automation_runbook.RegisterResource(),
		"azurerm_automation_schedule":                       automation_schedule.RegisterResource(),
		"azurerm_automation_variable_bool":                  automation_variable_bool.RegisterResource(),
		"azurerm_automation_variable_datetime":              automation_variable_datetime.RegisterResource(),
		"azurerm_automation_variable_int":                   automation_variable_int.RegisterResource(),
		"azurerm_automation_variable_object":                automation_variable_object.RegisterResource(),
		"azurerm_automation_variable_string":                automation_variable_string.RegisterResource(),
		"azurerm_automation_webhook":                        automation_webhook.RegisterResource(),
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
		automation_account.ListResource{},
	}
}
