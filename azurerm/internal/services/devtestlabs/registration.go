package devtestlabs

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Dev Test"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Dev Test",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_dev_test_lab":             dataSourceDevTestLab(),
		"azurerm_dev_test_virtual_network": dataSourceArmDevTestVirtualNetwork(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_dev_test_global_vm_shutdown_schedule": resourceDevTestGlobalVMShutdownSchedule(),
		"azurerm_dev_test_lab":                         resourceDevTestLab(),
		"azurerm_dev_test_schedule":                    resourceDevTestLabSchedules(),
		"azurerm_dev_test_linux_virtual_machine":       resourceArmDevTestLinuxVirtualMachine(),
		"azurerm_dev_test_policy":                      resourceArmDevTestPolicy(),
		"azurerm_dev_test_virtual_network":             resourceArmDevTestVirtualNetwork(),
		"azurerm_dev_test_windows_virtual_machine":     resourceArmDevTestWindowsVirtualMachine(),
	}
}
