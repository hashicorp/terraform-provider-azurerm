// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/actions/virtual_machine_power"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/availability_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/capacity_reservation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/capacity_reservation_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/dedicated_host"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/dedicated_host_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/disk_access"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/disk_encryption_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/gallery_application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/gallery_application_version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/image"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/linux_virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/linux_virtual_machine_scale_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/managed_disk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/managed_disk_sas_token"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/managed_disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/marketplace_agreement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/orchestrated_virtual_machine_scale_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/platform_image"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/proximity_placement_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/shared_image"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/shared_image_gallery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/shared_image_version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/shared_image_versions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/snapshot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/ssh_public_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_data_disk_attachment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_extension"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_gallery_application_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_implicit_data_disk_from_source"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_restore_point"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_restore_point_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_run_command"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_scale_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_scale_set_extension"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/virtual_machine_scale_set_standby_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/windows_virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/windows_virtual_machine_scale_set"
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
		"azurerm_availability_set":          availability_set.RegisterDataSource(),
		"azurerm_dedicated_host":            dedicated_host.RegisterDataSource(),
		"azurerm_dedicated_host_group":      dedicated_host_group.RegisterDataSource(),
		"azurerm_disk_access":               disk_access.RegisterDataSource(),
		"azurerm_disk_encryption_set":       disk_encryption_set.RegisterDataSource(),
		"azurerm_image":                     image.RegisterDataSource(),
		"azurerm_images":                    images.RegisterDataSource(),
		"azurerm_managed_disk":              managed_disk.RegisterDataSource(),
		"azurerm_marketplace_agreement":     marketplace_agreement.RegisterDataSource(),
		"azurerm_platform_image":            platform_image.RegisterDataSource(),
		"azurerm_proximity_placement_group": proximity_placement_group.RegisterDataSource(),
		"azurerm_shared_image":              shared_image.RegisterDataSource(),
		"azurerm_shared_image_gallery":      shared_image_gallery.RegisterDataSource(),
		"azurerm_shared_image_version":      shared_image_version.RegisterDataSource(),
		"azurerm_shared_image_versions":     shared_image_versions.RegisterDataSource(),
		"azurerm_snapshot":                  snapshot.RegisterDataSource(),
		"azurerm_ssh_public_key":            ssh_public_key.RegisterDataSource(),
		"azurerm_virtual_machine":           virtual_machine.RegisterDataSource(),
		"azurerm_virtual_machine_scale_set": virtual_machine_scale_set.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_availability_set":                       availability_set.RegisterResource(),
		"azurerm_capacity_reservation":                   capacity_reservation.RegisterResource(),
		"azurerm_capacity_reservation_group":             capacity_reservation_group.RegisterResource(),
		"azurerm_dedicated_host":                         dedicated_host.RegisterResource(),
		"azurerm_dedicated_host_group":                   dedicated_host_group.RegisterResource(),
		"azurerm_disk_access":                            disk_access.RegisterResource(),
		"azurerm_disk_encryption_set":                    disk_encryption_set.RegisterResource(),
		"azurerm_image":                                  image.RegisterResource(),
		"azurerm_linux_virtual_machine":                  linux_virtual_machine.RegisterResource(),
		"azurerm_linux_virtual_machine_scale_set":        linux_virtual_machine_scale_set.RegisterResource(),
		"azurerm_managed_disk":                           managed_disk.RegisterResource(),
		"azurerm_managed_disk_sas_token":                 managed_disk_sas_token.RegisterResource(),
		"azurerm_marketplace_agreement":                  marketplace_agreement.RegisterResource(),
		"azurerm_orchestrated_virtual_machine_scale_set": orchestrated_virtual_machine_scale_set.RegisterResource(),
		"azurerm_proximity_placement_group":              proximity_placement_group.RegisterResource(),
		"azurerm_shared_image":                           shared_image.RegisterResource(),
		"azurerm_shared_image_gallery":                   shared_image_gallery.RegisterResource(),
		"azurerm_shared_image_version":                   shared_image_version.RegisterResource(),
		"azurerm_snapshot":                               snapshot.RegisterResource(),
		"azurerm_ssh_public_key":                         ssh_public_key.RegisterResource(),
		"azurerm_virtual_machine_data_disk_attachment":   virtual_machine_data_disk_attachment.RegisterResource(),
		"azurerm_virtual_machine_extension":              virtual_machine_extension.RegisterResource(),
		"azurerm_virtual_machine_scale_set_extension":    virtual_machine_scale_set_extension.RegisterResource(),
		"azurerm_windows_virtual_machine":                windows_virtual_machine.RegisterResource(),
		"azurerm_windows_virtual_machine_scale_set":      windows_virtual_machine_scale_set.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		managed_disks.DataSource{},
		orchestrated_virtual_machine_scale_set.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		gallery_application.Resource{},
		gallery_application_version.Resource{},
		virtual_machine_gallery_application_assignment.Resource{},
		virtual_machine_implicit_data_disk_from_source.Resource{},
		virtual_machine_restore_point_collection.Resource{},
		virtual_machine_restore_point.Resource{},
		virtual_machine_run_command.Resource{},
		virtual_machine_scale_set_standby_pool.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		virtual_machine_power.RegisterAction,
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

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		availability_set.ListResource{},
		capacity_reservation_group.ListResource{},
		dedicated_host_group.ListResource{},
		linux_virtual_machine.ListResource{},
		windows_virtual_machine.ListResource{},
	}
}
