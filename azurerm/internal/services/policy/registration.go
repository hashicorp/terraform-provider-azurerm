package policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_policy_definition":     dataSourceArmPolicyDefinition(),
		"azurerm_policy_set_definition": dataSourceArmPolicySetDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_policy_assignment":                               resourceArmPolicyAssignment(),
		"azurerm_policy_definition":                               resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                           resourceArmPolicySetDefinition(),
		"azurerm_policy_remediation":                              resourceArmPolicyRemediation(),
		"azurerm_virtual_machine_configuration_policy_assignment": resourceVirtualMachineConfigurationPolicyAssignment(),
	}
}
