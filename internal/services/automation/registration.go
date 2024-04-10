// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AutomationVariablesDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AutomationConnectionTypeResource{},
		HybridRunbookWorkerGroupResource{},
		HybridRunbookWorkerResource{},
		SoftwareUpdateConfigurationResource{},
		SourceControlResource{},
		WatcherResource{},
		Python3PackageResource{},
		PowerShell72ModuleResource{},
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
		"azurerm_automation_account":           dataSourceAutomationAccount(),
		"azurerm_automation_variable_bool":     dataSourceAutomationVariableBool(),
		"azurerm_automation_variable_datetime": dataSourceAutomationVariableDateTime(),
		"azurerm_automation_variable_int":      dataSourceAutomationVariableInt(),
		"azurerm_automation_variable_object":   dataSourceAutomationVariableObject(),
		"azurerm_automation_variable_string":   dataSourceAutomationVariableString(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_automation_account":                        resourceAutomationAccount(),
		"azurerm_automation_certificate":                    resourceAutomationCertificate(),
		"azurerm_automation_connection":                     resourceAutomationConnection(),
		"azurerm_automation_connection_certificate":         resourceAutomationConnectionCertificate(),
		"azurerm_automation_connection_classic_certificate": resourceAutomationConnectionClassicCertificate(),
		"azurerm_automation_connection_service_principal":   resourceAutomationConnectionServicePrincipal(),
		"azurerm_automation_credential":                     resourceAutomationCredential(),
		"azurerm_automation_dsc_configuration":              resourceAutomationDscConfiguration(),
		"azurerm_automation_dsc_nodeconfiguration":          resourceAutomationDscNodeConfiguration(),
		"azurerm_automation_job_schedule":                   resourceAutomationJobSchedule(),
		"azurerm_automation_module":                         resourceAutomationModule(),
		"azurerm_automation_runbook":                        resourceAutomationRunbook(),
		"azurerm_automation_schedule":                       resourceAutomationSchedule(),
		"azurerm_automation_variable_bool":                  resourceAutomationVariableBool(),
		"azurerm_automation_variable_datetime":              resourceAutomationVariableDateTime(),
		"azurerm_automation_variable_int":                   resourceAutomationVariableInt(),
		"azurerm_automation_variable_object":                resourceAutomationVariableObject(),
		"azurerm_automation_variable_string":                resourceAutomationVariableString(),
		"azurerm_automation_webhook":                        resourceAutomationWebhook(),
	}
}
