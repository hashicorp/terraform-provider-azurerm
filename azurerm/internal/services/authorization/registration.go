package authorization

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Authorization"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_builtin_role_definition": dataSourceArmBuiltInRoleDefinition(),
		"azurerm_client_config":           dataSourceArmClientConfig(),
		"azurerm_role_definition":         dataSourceArmRoleDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_role_assignment": resourceArmRoleAssignment(),
		"azurerm_role_definition": resourceArmRoleDefinition(),
	}
}
