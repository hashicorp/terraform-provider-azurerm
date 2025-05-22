// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSharedImageVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageVersionCreate,
		Read:   resourceSharedImageVersionRead,
		Update: resourceSharedImageVersionUpdate,
		Delete: resourceSharedImageVersionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := galleryimageversions.ParseImageVersionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageVersionName,
			},

			"gallery_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"image_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"target_region": {
				// This needs to be a `TypeList` due to the `StateFunc` on the nested property `name`
				// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							StateFunc:        location.StateFunc,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},

						"regional_replica_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.DiskEncryptionSetID,
						},

						"exclude_from_latest_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						// The Service API doesn't support to update `storage_account_type`. So it has to recreate the resource for updating `storage_account_type`.
						// However, `ForceNew` cannot be used since resource would be recreated while adding or removing `target_region`.
						// And `CustomizeDiff` also cannot be used since it doesn't support in a `Set`.
						// So currently terraform would directly return the error message from Service API while updating this property. If this property needs to be updated, please recreate this pluginsdk.
						"storage_account_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(galleryimageversions.StorageAccountTypePremiumLRS),
								string(galleryimageversions.StorageAccountTypeStandardLRS),
								string(galleryimageversions.StorageAccountTypeStandardZRS),
							}, false),
							Default: string(galleryimageversions.StorageAccountTypeStandardLRS),
						},
					},
				},
			},

			"blob_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
				RequiredWith: []string{"storage_account_id"},
				ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"blob_uri"},
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"end_of_life_date": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"os_disk_snapshot_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
				// TODO -- add a validation function when snapshot has its own validation function
			},

			"managed_image_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					images.ValidateImageID,
					commonids.ValidateVirtualMachineID,
				),
				ExactlyOneOf: []string{"blob_uri", "os_disk_snapshot_id", "managed_image_id"},
			},

			"replication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimageversions.ReplicationModeFull),
					string(galleryimageversions.ReplicationModeShallow),
				}, false),
				Default: galleryimageversions.ReplicationModeFull,
			},

			"exclude_from_latest": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"deletion_of_replicated_locations_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("end_of_life_date", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) != "" && new.(string) == ""
			}),
		),
	}
}

func resourceSharedImageVersionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimageversions.NewImageVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_shared_image_version", id.ID())
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
				ReplicationMode:   pointer.To(galleryimageversions.ReplicationMode(d.Get("replication_mode").(string))),
				TargetRegions:     targetRegions,
			},
			SafetyProfile: &galleryimageversions.GalleryImageVersionSafetyProfile{
				AllowDeletionOfReplicatedLocations: utils.Bool(d.Get("deletion_of_replicated_locations_enabled").(bool)),
			},
			StorageProfile: galleryimageversions.GalleryImageVersionStorageProfile{},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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
				VirtualMachineId: utils.String(v.(string)),
			}
		} else {
			version.Properties.StorageProfile.Source = &galleryimageversions.GalleryArtifactVersionFullSource{
				Id: utils.String(v.(string)),
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

	if err := client.CreateOrUpdateThenPoll(ctx, id, version); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageVersionRead(d, meta)
}

func resourceSharedImageVersionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, galleryimageversions.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	payload := existing.Model
	if payload == nil {
		return fmt.Errorf("model is nil for %s", id)
	}

	if payload.Properties == nil {
		return fmt.Errorf("properties is nil for %s", id)
	}

	if payload.Properties.PublishingProfile == nil {
		payload.Properties.PublishingProfile = &galleryimageversions.GalleryArtifactPublishingProfileBase{}
	}

	if d.HasChange("target_region") {
		targetRegions, err := expandSharedImageVersionTargetRegions(d)
		if err != nil {
			return err
		}

		payload.Properties.PublishingProfile.TargetRegions = targetRegions
	}

	if d.HasChange("end_of_life_date") {
		endOfLifeDate, _ := time.Parse(time.RFC3339, d.Get("end_of_life_date").(string))
		payload.Properties.PublishingProfile.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if d.HasChange("exclude_from_latest") {
		payload.Properties.PublishingProfile.ExcludeFromLatest = pointer.To(d.Get("exclude_from_latest").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageVersionRead(d, meta)
}

func resourceSharedImageVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleryimageversions.DefaultGetOperationOptions())
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
				sourceID := ""
				if props.StorageProfile.OsDiskImage.Source.Id != nil {
					sourceID = *props.StorageProfile.OsDiskImage.Source.Id
				}

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
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceSharedImageVersionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimageversions.ParseImageVersionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// @tombuildsstuff: there appears to be an eventual consistency issue here
	timeout, _ := ctx.Deadline()
	log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   sharedImageVersionDeleteStateRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func sharedImageVersionDeleteStateRefreshFunc(ctx context.Context, client *galleryimageversions.GalleryImageVersionsClient, id galleryimageversions.ImageVersionId) pluginsdk.StateRefreshFunc {
	// Whilst the Shared Image Version is deleted quickly, it appears it's not actually finished replicating at this time
	// so the deletion of the parent Shared Image fails with "can not delete until nested resources are deleted"
	// ergo we need to poll on this for a bit
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Shared Image Version has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}

func expandSharedImageVersionTargetRegions(d *pluginsdk.ResourceData) (*[]galleryimageversions.TargetRegion, error) {
	vs := d.Get("target_region").([]interface{})
	results := make([]galleryimageversions.TargetRegion, 0)

	for _, v := range vs {
		input := v.(map[string]interface{})

		name := input["name"].(string)
		regionalReplicaCount := input["regional_replica_count"].(int)
		storageAccountType := input["storage_account_type"].(string)
		diskEncryptionSetId := input["disk_encryption_set_id"].(string)
		excludeFromLatest := input["exclude_from_latest_enabled"].(bool)

		output := galleryimageversions.TargetRegion{
			Name:                 name,
			ExcludeFromLatest:    pointer.To(excludeFromLatest),
			RegionalReplicaCount: pointer.To(int64(regionalReplicaCount)),
			StorageAccountType:   pointer.To(galleryimageversions.StorageAccountType(storageAccountType)),
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

func flattenSharedImageVersionTargetRegions(input *[]galleryimageversions.TargetRegion) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]interface{})

			output["name"] = location.Normalize(v.Name)

			if v.RegionalReplicaCount != nil {
				output["regional_replica_count"] = int(*v.RegionalReplicaCount)
			}

			if v.StorageAccountType != nil {
				output["storage_account_type"] = string(*v.StorageAccountType)
			}

			diskEncryptionSetId := ""
			if v.Encryption != nil && v.Encryption.OsDiskImage != nil && v.Encryption.OsDiskImage.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *v.Encryption.OsDiskImage.DiskEncryptionSetId
			}
			output["disk_encryption_set_id"] = diskEncryptionSetId

			output["exclude_from_latest_enabled"] = pointer.From(v.ExcludeFromLatest)

			results = append(results, output)
		}
	}

	return results
}
