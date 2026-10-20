// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceImageRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := images.ParseImageID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, images.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ImageName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			hyperVGeneration := ""
			if props.HyperVGeneration != nil {
				hyperVGeneration = string(*props.HyperVGeneration)
			}
			d.Set("hyper_v_generation", hyperVGeneration)

			// either source VM or storage profile can be specified, but not both
			if props.SourceVirtualMachine != nil && props.SourceVirtualMachine.Id != nil {
				d.Set("source_virtual_machine_id", pointer.From(props.SourceVirtualMachine.Id))
			} else {
				if err := d.Set("os_disk", flattenImageOSDisk(props.StorageProfile)); err != nil {
					return fmt.Errorf("setting `os_disk`: %+v", err)
				}
				if err := d.Set("data_disk", flattenImageDataDisks(props.StorageProfile)); err != nil {
					return fmt.Errorf("setting `data_disk`: %+v", err)
				}
				zoneResilient := false
				if props.StorageProfile != nil && props.StorageProfile.ZoneResilient != nil {
					zoneResilient = *props.StorageProfile.ZoneResilient
				}
				d.Set("zone_resilient", zoneResilient)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenImageOSDisk(input *images.ImageStorageProfile) []any {
	output := make([]any, 0)

	if input != nil {
		if v := input.OsDisk; v != nil {
			blobUri := pointer.From(v.BlobUri)
			caching := ""
			if v.Caching != nil {
				caching = string(*v.Caching)
			}
			diskSizeGB := 0
			if v.DiskSizeGB != nil {
				diskSizeGB = int(*v.DiskSizeGB)
			}
			managedDiskId := ""
			if disk := v.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
			diskEncryptionSetId := ""
			if set := v.DiskEncryptionSet; set != nil && set.Id != nil {
				encryptionId, _ := commonids.ParseDiskEncryptionSetIDInsensitively(*set.Id)
				diskEncryptionSetId = encryptionId.ID()
			}

			properties := map[string]any{
				"blob_uri":               blobUri,
				"caching":                caching,
				"managed_disk_id":        managedDiskId,
				"os_type":                string(v.OsType),
				"os_state":               string(v.OsState),
				"size_gb":                diskSizeGB,
				"disk_encryption_set_id": diskEncryptionSetId,
			}

			storageType := ""
			if v.StorageAccountType != nil {
				storageType = string(*v.StorageAccountType)
			}
			properties["storage_type"] = storageType

			output = append(output, properties)
		}
	}

	return output
}

func flattenImageDataDisks(input *images.ImageStorageProfile) []any {
	output := make([]any, 0)

	if input != nil {
		if v := input.DataDisks; v != nil {
			for _, disk := range *input.DataDisks {
				blobUri := pointer.From(disk.BlobUri)
				caching := ""
				if disk.Caching != nil {
					caching = string(*disk.Caching)
				}
				diskSizeGb := 0
				if disk.DiskSizeGB != nil {
					diskSizeGb = int(*disk.DiskSizeGB)
				}
				managedDiskId := ""
				if disk.ManagedDisk != nil && disk.ManagedDisk.Id != nil {
					managedDiskId = *disk.ManagedDisk.Id
				}
				diskEncryptionSetId := ""
				if set := disk.DiskEncryptionSet; set != nil && set.Id != nil {
					encryptionId, _ := commonids.ParseDiskEncryptionSetIDInsensitively(*set.Id)
					diskEncryptionSetId = encryptionId.ID()
				}

				properties := map[string]any{
					"blob_uri":               blobUri,
					"caching":                caching,
					"lun":                    int(disk.Lun),
					"managed_disk_id":        managedDiskId,
					"size_gb":                diskSizeGb,
					"disk_encryption_set_id": diskEncryptionSetId,
				}

				storageType := ""
				if disk.StorageAccountType != nil {
					storageType = string(*disk.StorageAccountType)
				}
				properties["storage_type"] = storageType

				output = append(output, properties)
			}
		}
	}

	return output
}
