package maintenance

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) Name() string {
	return "Maintenance"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Maintenance",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_maintenance_configuration": dataSourceMaintenanceConfiguration(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_maintenance_assignment_dedicated_host":  resourceArmMaintenanceAssignmentDedicatedHost(),
		"azurerm_maintenance_assignment_virtual_machine": resourceArmMaintenanceAssignmentVirtualMachine(),
		"azurerm_maintenance_configuration":              resourceArmMaintenanceConfiguration(),
	}
}
