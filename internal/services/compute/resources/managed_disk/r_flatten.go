// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceManagedDiskFlatten(d *pluginsdk.ResourceData, id *commonids.ManagedDiskId, model *disks.Disk, includeResource bool) error {
	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		d.Set("edge_zone", helpers.FlattenManagedDiskEdgeZone(model.ExtendedLocation))

		zone := ""
		if model.Zones != nil && len(*model.Zones) > 0 {
			z := *model.Zones
			zone = z[0]
		}
		d.Set("zone", zone)

		if sku := model.Sku; sku != nil {
			d.Set("storage_account_type", pointer.FromEnum(sku.Name))
		}

		if includeResource {
			if props := model.Properties; props != nil {
				creationData := props.CreationData
				d.Set("create_option", string(creationData.CreateOption))
				if creationData.LogicalSectorSize != nil {
					d.Set("logical_sector_size", creationData.LogicalSectorSize)
				}

				// imageReference is returned as well when galleryImageReference is used, only check imageReference when galleryImageReference is not returned
				galleryImageReferenceId := ""
				imageReferenceId := ""
				if creationData.GalleryImageReference != nil && creationData.GalleryImageReference.Id != nil {
					galleryImageReferenceId = pointer.From(creationData.GalleryImageReference.Id)
				} else if creationData.ImageReference != nil {
					imageReferenceId = pointer.From(creationData.ImageReference.Id)
				}
				d.Set("gallery_image_reference_id", galleryImageReferenceId)
				d.Set("image_reference_id", imageReferenceId)

				d.Set("performance_plus_enabled", creationData.PerformancePlus)
				d.Set("source_resource_id", creationData.SourceResourceId)
				d.Set("source_uri", creationData.SourceUri)
				d.Set("storage_account_id", creationData.StorageAccountId)
				d.Set("upload_size_bytes", creationData.UploadSizeBytes)

				d.Set("disk_size_gb", props.DiskSizeGB)
				d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
				d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
				d.Set("disk_iops_read_only", props.DiskIOPSReadOnly)
				d.Set("disk_mbps_read_only", props.DiskMBpsReadOnly)
				d.Set("optimized_frequent_attach_enabled", props.OptimizedForFrequentAttach)
				d.Set("os_type", pointer.FromEnum(props.OsType))
				d.Set("tier", props.Tier)
				d.Set("max_shares", props.MaxShares)
				d.Set("hyper_v_generation", pointer.FromEnum(props.HyperVGeneration))
				d.Set("network_access_policy", pointer.FromEnum(props.NetworkAccessPolicy))
				d.Set("disk_access_id", props.DiskAccessId)
				d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == disks.PublicNetworkAccessEnabled)

				diskEncryptionSetId := ""
				if props.Encryption != nil {
					diskEncryptionSetId = pointer.From(props.Encryption.DiskEncryptionSetId)
				}
				d.Set("disk_encryption_set_id", diskEncryptionSetId)

				if err := d.Set("encryption_settings", helpers.FlattenManagedDiskEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
					return fmt.Errorf("setting `encryption_settings`: %+v", err)
				}

				trustedLaunchEnabled := false
				securityType := ""
				secureVMDiskEncryptionSetId := ""
				if securityProfile := props.SecurityProfile; securityProfile != nil {
					if pointer.From(securityProfile.SecurityType) == disks.DiskSecurityTypesTrustedLaunch {
						trustedLaunchEnabled = true
					} else {
						securityType = pointer.FromEnum(securityProfile.SecurityType)
					}

					secureVMDiskEncryptionSetId = pointer.From(securityProfile.SecureVMDiskEncryptionSetId)
				}
				d.Set("trusted_launch_enabled", trustedLaunchEnabled)
				d.Set("security_type", securityType)
				d.Set("secure_vm_disk_encryption_set_id", secureVMDiskEncryptionSetId)

				d.Set("on_demand_bursting_enabled", pointer.From(props.BurstingEnabled))
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
