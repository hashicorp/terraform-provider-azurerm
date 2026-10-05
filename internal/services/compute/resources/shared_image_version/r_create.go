// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageVersionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimageversions.NewImageVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_shared_image_version", id.ID())
		}
	}

	targetRegions, err := expandSharedImageVersionTargetRegions(d)
	if err != nil {
		return err
	}

	version := galleryimageversions.GalleryImageVersion{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleryimageversions.GalleryImageVersionProperties{
			PublishingProfile: &galleryimageversions.GalleryArtifactPublishingProfileBase{
				ExcludeFromLatest: pointer.To(d.Get("exclude_from_latest").(bool)),
				ReplicationMode:   pointer.ToEnum[galleryimageversions.ReplicationMode](d.Get("replication_mode").(string)),
				TargetRegions:     targetRegions,
			},
			SafetyProfile: &galleryimageversions.GalleryImageVersionSafetyProfile{
				AllowDeletionOfReplicatedLocations: pointer.To(d.Get("deletion_of_replicated_locations_enabled").(bool)),
			},
			StorageProfile: galleryimageversions.GalleryImageVersionStorageProfile{},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("end_of_life_date"); ok {
		endOfLifeDate, _ := time.Parse(time.RFC3339, v.(string))
		version.Properties.PublishingProfile.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if v, ok := d.GetOk("managed_image_id"); ok {
		_, err := virtualmachines.ParseVirtualMachineID(v.(string))
		if err == nil {
			version.Properties.StorageProfile.Source = &galleryimageversions.GalleryArtifactVersionFullSource{
				VirtualMachineId: pointer.To(v.(string)),
			}
		} else {
			version.Properties.StorageProfile.Source = &galleryimageversions.GalleryArtifactVersionFullSource{
				Id: pointer.To(v.(string)),
			}
		}
	}

	if v, ok := d.GetOk("os_disk_snapshot_id"); ok {
		version.Properties.StorageProfile.OsDiskImage = &galleryimageversions.GalleryDiskImage{
			Source: &galleryimageversions.GalleryDiskImageSource{
				Id: pointer.To(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("blob_uri"); ok {
		version.Properties.StorageProfile.OsDiskImage = &galleryimageversions.GalleryDiskImage{
			Source: &galleryimageversions.GalleryDiskImageSource{
				StorageAccountId: pointer.To(d.Get("storage_account_id").(string)),
				Uri:              pointer.To(v.(string)),
			},
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, version, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	readCtx, cancelCtx := context.WithTimeout(ctx, 5*time.Minute)
	defer cancelCtx()
	if err = retry.RetryContext(readCtx, 5*time.Second, func() *retry.RetryError {
		read, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(read.HttpResponse) {
				return retry.RetryableError(fmt.Errorf("waiting for creation of %s", id))
			}
			return retry.NonRetryableError(err)
		}
		if read.Model == nil {
			return retry.RetryableError(fmt.Errorf("waiting for `model` to become available for %s", id))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageVersionRead(d, meta)
}

func expandSharedImageVersionTargetRegions(d *pluginsdk.ResourceData) (*[]galleryimageversions.TargetRegion, error) {
	vs := d.Get("target_region").([]any)
	results := make([]galleryimageversions.TargetRegion, 0)

	for _, v := range vs {
		input := v.(map[string]any)

		name := input["name"].(string)
		regionalReplicaCount := input["regional_replica_count"].(int)
		storageAccountType := input["storage_account_type"].(string)
		diskEncryptionSetId := input["disk_encryption_set_id"].(string)
		excludeFromLatest := input["exclude_from_latest_enabled"].(bool)

		output := galleryimageversions.TargetRegion{
			Name:                 name,
			ExcludeFromLatest:    pointer.To(excludeFromLatest),
			RegionalReplicaCount: pointer.To(int64(regionalReplicaCount)),
			StorageAccountType:   pointer.ToEnum[galleryimageversions.StorageAccountType](storageAccountType),
		}

		if diskEncryptionSetId != "" {
			if d.Get("replication_mode").(string) == string(galleryimageversions.ReplicationModeShallow) {
				return nil, fmt.Errorf("`disk_encryption_set_id` cannot be used when `replication_mode` is `Shallow`")
			}

			output.Encryption = &galleryimageversions.EncryptionImages{
				OsDiskImage: &galleryimageversions.OSDiskImageEncryption{
					DiskEncryptionSetId: pointer.To(diskEncryptionSetId),
				},
			}
		}

		results = append(results, output)
	}

	return &results, nil
}
