package policy

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Policy"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Policy",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_policy_definition":     dataSourceArmPolicyDefinition(),
		"azurerm_policy_set_definition": dataSourceArmPolicySetDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_policy_assignment":                               resourceArmPolicyAssignment(),
		"azurerm_policy_definition":                               resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                           resourceArmPolicySetDefinition(),
		"azurerm_policy_remediation":                              resourceArmPolicyRemediation(),
		"azurerm_virtual_machine_configuration_policy_assignment": resourceVirtualMachineConfigurationPolicyAssignment(),
	}
}
