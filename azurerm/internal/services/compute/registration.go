package compute

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_availability_set":          dataSourceAvailabilitySet(),
		"azurerm_dedicated_host":            dataSourceDedicatedHost(),
		"azurerm_dedicated_host_group":      dataSourceDedicatedHostGroup(),
		"azurerm_disk_encryption_set":       dataSourceDiskEncryptionSet(),
		"azurerm_managed_disk":              dataSourceManagedDisk(),
		"azurerm_image":                     dataSourceImage(),
		"azurerm_images":                    dataSourceImages(),
		"azurerm_disk_access":               dataSourceDiskAccess(),
		"azurerm_platform_image":            dataSourcePlatformImage(),
		"azurerm_proximity_placement_group": dataSourceProximityPlacementGroup(),
		"azurerm_shared_image_gallery":      dataSourceSharedImageGallery(),
		"azurerm_shared_image_version":      dataSourceSharedImageVersion(),
		"azurerm_shared_image_versions":     dataSourceSharedImageVersions(),
		"azurerm_shared_image":              dataSourceSharedImage(),
		"azurerm_snapshot":                  dataSourceSnapshot(),
		"azurerm_virtual_machine":           dataSourceVirtualMachine(),
		"azurerm_virtual_machine_scale_set": dataSourceVirtualMachineScaleSet(),
		"azurerm_ssh_public_key":            dataSourceSshPublicKey(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_availability_set":                       resourceAvailabilitySet(),
		"azurerm_dedicated_host":                         resourceDedicatedHost(),
		"azurerm_dedicated_host_group":                   resourceDedicatedHostGroup(),
		"azurerm_disk_encryption_set":                    resourceDiskEncryptionSet(),
		"azurerm_image":                                  resourceImage(),
		"azurerm_managed_disk":                           resourceManagedDisk(),
		"azurerm_disk_access":                            resourceDiskAccess(),
		"azurerm_marketplace_agreement":                  resourceMarketplaceAgreement(),
		"azurerm_proximity_placement_group":              resourceProximityPlacementGroup(),
		"azurerm_shared_image_gallery":                   resourceSharedImageGallery(),
		"azurerm_shared_image_version":                   resourceSharedImageVersion(),
		"azurerm_shared_image":                           resourceSharedImage(),
		"azurerm_snapshot":                               resourceSnapshot(),
		"azurerm_virtual_machine_data_disk_attachment":   resourceVirtualMachineDataDiskAttachment(),
		"azurerm_virtual_machine_extension":              resourceVirtualMachineExtension(),
		"azurerm_virtual_machine_scale_set":              resourceVirtualMachineScaleSet(),
		"azurerm_orchestrated_virtual_machine_scale_set": resourceOrchestratedVirtualMachineScaleSet(),
		"azurerm_virtual_machine":                        resourceVirtualMachine(),
		"azurerm_linux_virtual_machine":                  resourceLinuxVirtualMachine(),
		"azurerm_linux_virtual_machine_scale_set":        resourceLinuxVirtualMachineScaleSet(),
		"azurerm_virtual_machine_scale_set_extension":    resourceVirtualMachineScaleSetExtension(),
		"azurerm_windows_virtual_machine":                resourceWindowsVirtualMachine(),
		"azurerm_windows_virtual_machine_scale_set":      resourceWindowsVirtualMachineScaleSet(),
		"azurerm_ssh_public_key":                         resourceSshPublicKey(),
	}

	return resources
}
