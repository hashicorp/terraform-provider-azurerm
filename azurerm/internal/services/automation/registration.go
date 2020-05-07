package automation

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_automation_account":           dataSourceArmAutomationAccount(),
		"azurerm_automation_variable_bool":     dataSourceArmAutomationVariableBool(),
		"azurerm_automation_variable_datetime": dataSourceArmAutomationVariableDateTime(),
		"azurerm_automation_variable_int":      dataSourceArmAutomationVariableInt(),
		"azurerm_automation_variable_string":   dataSourceArmAutomationVariableString(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_automation_account":               resourceArmAutomationAccount(),
		"azurerm_automation_certificate":           resourceArmAutomationCertificate(),
		"azurerm_automation_credential":            resourceArmAutomationCredential(),
		"azurerm_automation_dsc_configuration":     resourceArmAutomationDscConfiguration(),
		"azurerm_automation_dsc_nodeconfiguration": resourceArmAutomationDscNodeConfiguration(),
		"azurerm_automation_job_schedule":          resourceArmAutomationJobSchedule(),
		"azurerm_automation_module":                resourceArmAutomationModule(),
		"azurerm_automation_runbook":               resourceArmAutomationRunbook(),
		"azurerm_automation_schedule":              resourceArmAutomationSchedule(),
		"azurerm_automation_variable_bool":         resourceArmAutomationVariableBool(),
		"azurerm_automation_variable_datetime":     resourceArmAutomationVariableDateTime(),
		"azurerm_automation_variable_int":          resourceArmAutomationVariableInt(),
		"azurerm_automation_variable_string":       resourceArmAutomationVariableString(),
	}
}
