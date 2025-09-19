// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

var _ sdk.FrameworkServiceRegistration = Registration{}

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
		"azurerm_marketplace_agreement":     dataSourceMarketplaceAgreement(),
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
		"azurerm_capacity_reservation":                   resourceCapacityReservation(),
		"azurerm_capacity_reservation_group":             resourceCapacityReservationGroup(),
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
		"azurerm_orchestrated_virtual_machine_scale_set": resourceOrchestratedVirtualMachineScaleSet(),
		"azurerm_linux_virtual_machine":                  resourceLinuxVirtualMachine(),
		"azurerm_linux_virtual_machine_scale_set":        resourceLinuxVirtualMachineScaleSet(),
		"azurerm_virtual_machine_scale_set_extension":    resourceVirtualMachineScaleSetExtension(),
		"azurerm_windows_virtual_machine":                resourceWindowsVirtualMachine(),
		"azurerm_windows_virtual_machine_scale_set":      resourceWindowsVirtualMachineScaleSet(),
		"azurerm_ssh_public_key":                         resourceSshPublicKey(),
		"azurerm_managed_disk_sas_token":                 resourceManagedDiskSasToken(),
	}

	return resources
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ManagedDisksDataSource{},
		OrchestratedVirtualMachineScaleSetDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		VirtualMachineImplicitDataDiskFromSourceResource{},
		VirtualMachineRunCommandResource{},
		GalleryApplicationResource{},
		GalleryApplicationVersionResource{},
		RestorePointCollectionResource{},
		VirtualMachineRestorePointCollectionResource{},
		VirtualMachineRestorePointResource{},
		VirtualMachineGalleryApplicationAssignmentResource{},
		VirtualMachineScaleSetStandbyPoolResource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		newVirtualMachinePowerAction,
	}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []func() list.ListResource {
	return []func() list.ListResource{}
}
