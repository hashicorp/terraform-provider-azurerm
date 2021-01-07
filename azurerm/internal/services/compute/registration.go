package compute

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Compute"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Compute",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_availability_set":          dataSourceAvailabilitySet(),
		"azurerm_dedicated_host":            dataSourceDedicatedHost(),
		"azurerm_dedicated_host_group":      dataSourceDedicatedHostGroup(),
		"azurerm_disk_encryption_set":       dataSourceDiskEncryptionSet(),
		"azurerm_managed_disk":              dataSourceArmManagedDisk(),
		"azurerm_image":                     dataSourceArmImage(),
		"azurerm_images":                    dataSourceArmImages(),
		"azurerm_platform_image":            dataSourceArmPlatformImage(),
		"azurerm_proximity_placement_group": dataSourceArmProximityPlacementGroup(),
		"azurerm_shared_image_gallery":      dataSourceArmSharedImageGallery(),
		"azurerm_shared_image_version":      dataSourceArmSharedImageVersion(),
		"azurerm_shared_image_versions":     dataSourceArmSharedImageVersions(),
		"azurerm_shared_image":              dataSourceArmSharedImage(),
		"azurerm_snapshot":                  dataSourceArmSnapshot(),
		"azurerm_virtual_machine":           dataSourceArmVirtualMachine(),
		"azurerm_virtual_machine_scale_set": dataSourceArmVirtualMachineScaleSet(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	resources := map[string]*schema.Resource{
		"azurerm_availability_set":                       resourceAvailabilitySet(),
		"azurerm_dedicated_host":                         resourceDedicatedHost(),
		"azurerm_dedicated_host_group":                   resourceDedicatedHostGroup(),
		"azurerm_disk_encryption_set":                    resourceDiskEncryptionSet(),
		"azurerm_image":                                  resourceArmImage(),
		"azurerm_managed_disk":                           resourceArmManagedDisk(),
		"azurerm_marketplace_agreement":                  resourceArmMarketplaceAgreement(),
		"azurerm_proximity_placement_group":              resourceArmProximityPlacementGroup(),
		"azurerm_shared_image_gallery":                   resourceArmSharedImageGallery(),
		"azurerm_shared_image_version":                   resourceArmSharedImageVersion(),
		"azurerm_shared_image":                           resourceArmSharedImage(),
		"azurerm_snapshot":                               resourceArmSnapshot(),
		"azurerm_virtual_machine_data_disk_attachment":   resourceArmVirtualMachineDataDiskAttachment(),
		"azurerm_virtual_machine_extension":              resourceArmVirtualMachineExtension(),
		"azurerm_virtual_machine_scale_set":              resourceArmVirtualMachineScaleSet(),
		"azurerm_orchestrated_virtual_machine_scale_set": resourceArmOrchestratedVirtualMachineScaleSet(),
		"azurerm_virtual_machine":                        resourceArmVirtualMachine(),
		"azurerm_linux_virtual_machine":                  resourceLinuxVirtualMachine(),
		"azurerm_linux_virtual_machine_scale_set":        resourceArmLinuxVirtualMachineScaleSet(),
		"azurerm_virtual_machine_scale_set_extension":    resourceArmVirtualMachineScaleSetExtension(),
		"azurerm_windows_virtual_machine":                resourceWindowsVirtualMachine(),
		"azurerm_windows_virtual_machine_scale_set":      resourceArmWindowsVirtualMachineScaleSet(),
	}

	return resources
}
