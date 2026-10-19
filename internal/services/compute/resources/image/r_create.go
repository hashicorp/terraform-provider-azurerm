// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceImageCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := images.NewImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, images.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_image", id.ID())
			}
		}
	}

	props := images.ImageProperties{
		HyperVGeneration: pointer.ToEnum[images.HyperVGenerationTypes](d.Get("hyper_v_generation").(string)),
	}

	sourceVM := images.SubResource{}
	if v, ok := d.GetOk("source_virtual_machine_id"); ok {
		sourceVM.Id = pointer.To(v.(string))
	}

	storageProfile := images.ImageStorageProfile{
		OsDisk:        expandImageOSDisk(d.Get("os_disk").([]any)),
		DataDisks:     expandImageDataDisks(d.Get("data_disk").([]any)),
		ZoneResilient: pointer.To(d.Get("zone_resilient").(bool)),
	}

	// either source VM or storage profile can be specified, but not both
	if sourceVM.Id == nil {
		// if both sourceVM and storageProfile are empty, return an error
		if storageProfile.OsDisk == nil && (storageProfile.DataDisks == nil || len(*storageProfile.DataDisks) == 0) {
			return fmt.Errorf("[ERROR] Cannot create image when both source VM and storage profile are empty")
		}

		props.StorageProfile = &storageProfile
	} else {
		// creating an image from source VM
		props.SourceVirtualMachine = &sourceVM
	}

	payload := images.Image{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &props,
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceImageRead(d, meta)
}

func expandImageOSDisk(input []any) *images.ImageOSDisk {
	if len(input) > 0 {
		config := input[0].(map[string]any)

		out := &images.ImageOSDisk{}

		if v := config["os_type"].(string); v != "" {
			out.OsType = images.OperatingSystemTypes(v)
		}

		if v := config["os_state"].(string); v != "" {
			out.OsState = images.OperatingSystemStateTypes(v)
		}
		managedDiskID := config["managed_disk_id"].(string)
		if managedDiskID != "" {
			out.ManagedDisk = &images.SubResource{
				Id: &managedDiskID,
			}
		}

		out.BlobUri = pointer.To(config["blob_uri"].(string))

		if v := config["caching"].(string); v != "" {
			out.Caching = pointer.ToEnum[images.CachingTypes](v)
		}

		if size := config["size_gb"]; size != 0 {
			out.DiskSizeGB = pointer.To(int64(size.(int)))
		}

		if id := config["disk_encryption_set_id"].(string); id != "" {
			out.DiskEncryptionSet = &images.SubResource{
				Id: pointer.To(id),
			}
		}

		out.StorageAccountType = pointer.ToEnum[images.StorageAccountTypes](config["storage_type"].(string))

		return out
	}

	return nil
}

func expandImageDataDisks(disks []any) *[]images.ImageDataDisk {
	output := make([]images.ImageDataDisk, 0)
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]any)

		item := images.ImageDataDisk{
			BlobUri: pointer.To(config["blob_uri"].(string)),
			Lun:     int64(config["lun"].(int)),
		}

		if size := config["size_gb"]; size != 0 {
			item.DiskSizeGB = pointer.To(int64(size.(int)))
		}

		if v := config["caching"].(string); v != "" {
			item.Caching = pointer.ToEnum[images.CachingTypes](v)
		}

		if managedDiskID := config["managed_disk_id"].(string); managedDiskID != "" {
			item.ManagedDisk = &images.SubResource{
				Id: &managedDiskID,
			}
		}

		if id := config["disk_encryption_set_id"].(string); id != "" {
			item.DiskEncryptionSet = &images.SubResource{
				Id: pointer.To(id),
			}
		}

		item.StorageAccountType = pointer.ToEnum[images.StorageAccountTypes](config["storage_type"].(string))

		output = append(output, item)
	}

	return &output
}
