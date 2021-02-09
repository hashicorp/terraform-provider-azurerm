package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSharedImageVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedImageVersionCreateUpdate,
		Read:   resourceSharedImageVersionRead,
		Update: resourceSharedImageVersionCreateUpdate,
		Delete: resourceSharedImageVersionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SharedImageVersionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageVersionName,
			},

			"gallery_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"image_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"target_region": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        location.StateFunc,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},

						"regional_replica_count": {
							Type:     schema.TypeInt,
							Required: true,
						},

						// The Service API doesn't support to update `storage_account_type`. So it has to recreate the resource for updating `storage_account_type`.
						// However, `ForceNew` cannot be used since resource would be recreated while adding or removing `target_region`.
						// And `CustomizeDiff` also cannot be used since it doesn't support in a `Set`.
						// So currently terraform would directly return the error message from Service API while updating this property. If this property needs to be updated, please recreate this resource.
						"storage_account_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StorageAccountTypeStandardLRS),
								string(compute.StorageAccountTypeStandardZRS),
							}, false),
							Default: string(compute.StorageAccountTypeStandardLRS),
						},
					},
				},
			},

			"os_disk_snapshot_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"os_disk_snapshot_id", "managed_image_id"},
				// TODO -- add a validation function when snapshot has its own validation function
			},

			"managed_image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validate.ImageID,
					validate.VirtualMachineID,
				),
				ExactlyOneOf: []string{"os_disk_snapshot_id", "managed_image_id"},
			},

			"exclude_from_latest": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSharedImageVersionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	imageVersion := d.Get("name").(string)
	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_shared_image_version", *existing.ID)
		}
	}

	version := compute.GalleryImageVersion{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		GalleryImageVersionProperties: &compute.GalleryImageVersionProperties{
			PublishingProfile: &compute.GalleryImageVersionPublishingProfile{
				ExcludeFromLatest: utils.Bool(d.Get("exclude_from_latest").(bool)),
				TargetRegions:     expandSharedImageVersionTargetRegions(d),
			},
			StorageProfile: &compute.GalleryImageVersionStorageProfile{},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	future, err := client.CreateOrUpdate(ctx, resourceGroup, galleryName, imageName, imageVersion, version)
	if err != nil {
		return fmt.Errorf("Error creating Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the creation of Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceSharedImageVersionRead(d, meta)
}

func resourceSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", id.VersionName, id.ImageName, id.GalleryName, id.ResourceGroup, err)
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
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			if err := d.Set("target_region", flattenSharedImageVersionTargetRegions(profile.TargetRegions)); err != nil {
				return fmt.Errorf("Error setting `target_region`: %+v", err)
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

func resourceSharedImageVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageVersionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", id.VersionName, id.ImageName, id.GalleryName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", id.VersionName, id.ImageName, id.GalleryName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandSharedImageVersionTargetRegions(d *schema.ResourceData) *[]compute.TargetRegion {
	vs := d.Get("target_region").(*schema.Set)
	results := make([]compute.TargetRegion, 0)

	for _, v := range vs.List() {
		input := v.(map[string]interface{})

		name := input["name"].(string)
		regionalReplicaCount := input["regional_replica_count"].(int)
		storageAccountType := input["storage_account_type"].(string)

		output := compute.TargetRegion{
			Name:                 utils.String(name),
			RegionalReplicaCount: utils.Int32(int32(regionalReplicaCount)),
			StorageAccountType:   compute.StorageAccountType(storageAccountType),
		}
		results = append(results, output)
	}

	return &results
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

			results = append(results, output)
		}
	}

	return results
}
