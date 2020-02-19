package policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Policy"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_policy_definition": dataSourceArmPolicyDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_policy_assignment":     resourceArmPolicyAssignment(),
		"azurerm_policy_definition":     resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition": resourceArmPolicySetDefinition(),
	}
}
