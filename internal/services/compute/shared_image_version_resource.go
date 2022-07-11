package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSharedImageVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageVersionCreateUpdate,
		Read:   resourceSharedImageVersionRead,
		Update: resourceSharedImageVersionCreateUpdate,
		Delete: resourceSharedImageVersionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SharedImageVersionID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

						// The Service API doesn't support to update `storage_account_type`. So it has to recreate the resource for updating `storage_account_type`.
						// However, `ForceNew` cannot be used since resource would be recreated while adding or removing `target_region`.
						// And `CustomizeDiff` also cannot be used since it doesn't support in a `Set`.
						// So currently terraform would directly return the error message from Service API while updating this property. If this property needs to be updated, please recreate this pluginsdk.
						"storage_account_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StorageAccountTypePremiumLRS),
								string(compute.StorageAccountTypeStandardLRS),
								string(compute.StorageAccountTypeStandardZRS),
							}, false),
							Default: string(compute.StorageAccountTypeStandardLRS),
						},
					},
				},
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
				ExactlyOneOf: []string{"os_disk_snapshot_id", "managed_image_id"},
				// TODO -- add a validation function when snapshot has its own validation function
			},

			"managed_image_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validate.ImageID,
					validate.VirtualMachineID,
				),
				ExactlyOneOf: []string{"os_disk_snapshot_id", "managed_image_id"},
			},

			"replication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.ReplicationModeFull),
					string(compute.ReplicationModeShallow),
				}, false),
				Default: compute.ReplicationModeFull,
			},

			"exclude_from_latest": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("end_of_life_date", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) != "" && new.(string) == ""
			}),
		),
	}
}

func resourceSharedImageVersionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedImageVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_shared_image_version", id.ID())
		}
	}

	targetRegions, err := expandSharedImageVersionTargetRegions(d)
	if err != nil {
		return err
	}

	version := compute.GalleryImageVersion{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		GalleryImageVersionProperties: &compute.GalleryImageVersionProperties{
			PublishingProfile: &compute.GalleryImageVersionPublishingProfile{
				ExcludeFromLatest: utils.Bool(d.Get("exclude_from_latest").(bool)),
				ReplicationMode:   compute.ReplicationMode(d.Get("replication_mode").(string)),
				TargetRegions:     targetRegions,
			},
			StorageProfile: &compute.GalleryImageVersionStorageProfile{},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("end_of_life_date"); ok {
		endOfLifeDate, _ := time.Parse(time.RFC3339, v.(string))
		version.GalleryImageVersionProperties.PublishingProfile.EndOfLifeDate = &date.Time{
			Time: endOfLifeDate,
		}
	}

	if v, ok := d.GetOk("managed_image_id"); ok {
		version.GalleryImageVersionProperties.StorageProfile.Source = &compute.GalleryArtifactVersionSource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("os_disk_snapshot_id"); ok {
		version.GalleryImageVersionProperties.StorageProfile.OsDiskImage = &compute.GalleryOSDiskImage{
			Source: &compute.GalleryArtifactVersionSource{
				ID: utils.String(v.(string)),
			},
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, version)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageVersionRead(d, meta)
}

func resourceSharedImageVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, compute.ReplicationStatusTypesReplicationStatus)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image Version %q (Image %q / Gallery %q / Resource Group %q) was not found - removing from state", id.VersionName, id.ImageName, id.GalleryName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", id.VersionName, id.ImageName, id.GalleryName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("image_name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			if v := profile.EndOfLifeDate; v != nil {
				d.Set("end_of_life_date", profile.EndOfLifeDate.Format(time.RFC3339))
			}

			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			replicationMode := string(compute.ReplicationModeFull)
			if profile.ReplicationMode != "" {
				replicationMode = string(profile.ReplicationMode)
			}
			d.Set("replication_mode", replicationMode)

			if err := d.Set("target_region", flattenSharedImageVersionTargetRegions(profile.TargetRegions)); err != nil {
				return fmt.Errorf("setting `target_region`: %+v", err)
			}
		}

		if profile := props.StorageProfile; profile != nil {
			if source := profile.Source; source != nil {
				d.Set("managed_image_id", source.ID)
			}

			osDiskSnapShotID := ""
			if profile.OsDiskImage != nil && profile.OsDiskImage.Source != nil && profile.OsDiskImage.Source.ID != nil {
				osDiskSnapShotID = *profile.OsDiskImage.Source.ID
			}
			d.Set("os_disk_snapshot_id", osDiskSnapShotID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSharedImageVersionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageVersionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
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

func sharedImageVersionDeleteStateRefreshFunc(ctx context.Context, client *compute.GalleryImageVersionsClient, id parse.SharedImageVersionId) pluginsdk.StateRefreshFunc {
	// Whilst the Shared Image Version is deleted quickly, it appears it's not actually finished replicating at this time
	// so the deletion of the parent Shared Image fails with "can not delete until nested resources are deleted"
	// ergo we need to poll on this for a bit
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, "")
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Shared Image Version has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}

func expandSharedImageVersionTargetRegions(d *pluginsdk.ResourceData) (*[]compute.TargetRegion, error) {
	vs := d.Get("target_region").([]interface{})
	results := make([]compute.TargetRegion, 0)

	for _, v := range vs {
		input := v.(map[string]interface{})

		name := input["name"].(string)
		regionalReplicaCount := input["regional_replica_count"].(int)
		storageAccountType := input["storage_account_type"].(string)
		diskEncryptionSetId := input["disk_encryption_set_id"].(string)

		output := compute.TargetRegion{
			Name:                 utils.String(name),
			RegionalReplicaCount: utils.Int32(int32(regionalReplicaCount)),
			StorageAccountType:   compute.StorageAccountType(storageAccountType),
		}

		if diskEncryptionSetId != "" {
			if d.Get("replication_mode").(string) == string(compute.ReplicationModeShallow) {
				return nil, fmt.Errorf("`disk_encryption_set_id` cannot be used when `replication_mode` is `Shallow`")
			}

			output.Encryption = &compute.EncryptionImages{
				OsDiskImage: &compute.OSDiskImageEncryption{
					DiskEncryptionSetID: utils.String(diskEncryptionSetId),
				},
			}
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenSharedImageVersionTargetRegions(input *[]compute.TargetRegion) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]interface{})

			if v.Name != nil {
				output["name"] = azure.NormalizeLocation(*v.Name)
			}

			if v.RegionalReplicaCount != nil {
				output["regional_replica_count"] = int(*v.RegionalReplicaCount)
			}

			output["storage_account_type"] = string(v.StorageAccountType)

			diskEncryptionSetId := ""
			if v.Encryption != nil && v.Encryption.OsDiskImage != nil && v.Encryption.OsDiskImage.DiskEncryptionSetID != nil {
				diskEncryptionSetId = *v.Encryption.OsDiskImage.DiskEncryptionSetID
			}
			output["disk_encryption_set_id"] = diskEncryptionSetId

			results = append(results, output)
		}
	}

	return results
}
