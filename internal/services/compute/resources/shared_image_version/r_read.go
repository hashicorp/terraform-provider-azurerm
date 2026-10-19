// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageVersionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleryimageversions.GetOperationOptions{
		Expand: pointer.To(galleryimageversions.ReplicationStatusTypesUefiSettings),
	})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.VersionName)
	d.Set("image_name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			if profile := props.PublishingProfile; profile != nil {
				if v := profile.EndOfLifeDate; v != nil {
					d.Set("end_of_life_date", profile.EndOfLifeDate)
				}

				d.Set("exclude_from_latest", profile.ExcludeFromLatest)

				replicationMode := string(galleryimageversions.ReplicationModeFull)
				if profile.ReplicationMode != nil {
					replicationMode = string(*profile.ReplicationMode)
				}
				d.Set("replication_mode", replicationMode)

				if err := d.Set("target_region", flattenSharedImageVersionTargetRegions(profile.TargetRegions)); err != nil {
					return fmt.Errorf("setting `target_region`: %+v", err)
				}
			}

			if source := props.StorageProfile.Source; source != nil {
				if source.Id != nil {
					d.Set("managed_image_id", source.Id)
				}

				if source.VirtualMachineId != nil {
					d.Set("managed_image_id", source.VirtualMachineId)
				}
			}

			blobURI := ""
			if props.StorageProfile.OsDiskImage != nil && props.StorageProfile.OsDiskImage.Source != nil && props.StorageProfile.OsDiskImage.Source.Uri != nil {
				blobURI = *props.StorageProfile.OsDiskImage.Source.Uri
			}
			d.Set("blob_uri", blobURI)

			osDiskSnapShotID := ""
			storageAccountID := ""
			if props.StorageProfile.OsDiskImage != nil && props.StorageProfile.OsDiskImage.Source != nil {
				sourceID := pointer.From(props.StorageProfile.OsDiskImage.Source.Id)

				if props.StorageProfile.OsDiskImage.Source.StorageAccountId != nil {
					sourceID = *props.StorageProfile.OsDiskImage.Source.StorageAccountId
				}

				if blobURI == "" {
					osDiskSnapShotID = sourceID
				} else {
					storageAccountID = sourceID
				}
			}

			d.Set("os_disk_snapshot_id", osDiskSnapShotID)
			d.Set("storage_account_id", storageAccountID)

			if safetyProfile := props.SafetyProfile; safetyProfile != nil {
				d.Set("deletion_of_replicated_locations_enabled", pointer.From(safetyProfile.AllowDeletionOfReplicatedLocations))
			}

			if securityProfile := props.SecurityProfile; securityProfile != nil && securityProfile.UefiSettings != nil {
				if err := d.Set("uefi", flattenUefiSettings(securityProfile.UefiSettings)); err != nil {
					return fmt.Errorf("setting `uefi`: %+v", err)
				}
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
