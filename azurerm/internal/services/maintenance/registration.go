package maintenance

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_maintenance_configuration": dataSourceMaintenanceConfiguration(),
	}
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_maintenance_assignment_dedicated_host":  resourceArmMaintenanceAssignmentDedicatedHost(),
		"azurerm_maintenance_assignment_virtual_machine": resourceArmMaintenanceAssignmentVirtualMachine(),
		"azurerm_maintenance_configuration":              resourceArmMaintenanceConfiguration(),
	}
}
