package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSharedImageVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSharedImageVersionCreateUpdate,
		Read:   resourceArmSharedImageVersionRead,
		Update: resourceArmSharedImageVersionCreateUpdate,
		Delete: resourceArmSharedImageVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"managed_image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_region": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        azure.NormalizeLocation,
							DiffSuppressFunc: azure.SuppressLocationDiff,
						},

						"regional_replica_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
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

func resourceArmSharedImageVersionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	imageVersion := d.Get("name").(string)
	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	managedImageId := d.Get("managed_image_id").(string)
	excludeFromLatest := d.Get("exclude_from_latest").(bool)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	targetRegions := expandSharedImageVersionTargetRegions(d)
	t := d.Get("tags").(map[string]interface{})

	version := compute.GalleryImageVersion{
		Location: utils.String(location),
		GalleryImageVersionProperties: &compute.GalleryImageVersionProperties{
			PublishingProfile: &compute.GalleryImageVersionPublishingProfile{
				ExcludeFromLatest: utils.Bool(excludeFromLatest),
				TargetRegions:     targetRegions,
				Source: &compute.GalleryArtifactSource{
					ManagedImage: &compute.ManagedArtifact{
						ID: utils.String(managedImageId),
					},
				},
			},
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, galleryName, imageName, imageVersion, version)
	if err != nil {
		return fmt.Errorf("Error creating Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the creation of Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	// TODO: poll?

	read, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmSharedImageVersionRead(d, meta)
}

func resourceArmSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	imageVersion := id.Path["versions"]
	imageName := id.Path["images"]
	galleryName := id.Path["galleries"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, compute.ReplicationStatusTypesReplicationStatus)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image Version %q (Image %q / Gallery %q / Resource Group %q) was not found - removing from state", imageVersion, imageName, galleryName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("image_name", imageName)
	d.Set("gallery_name", galleryName)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			flattenedRegions := flattenSharedImageVersionTargetRegions(profile.TargetRegions)
			if err := d.Set("target_region", flattenedRegions); err != nil {
				return fmt.Errorf("Error setting `target_region`: %+v", err)
			}

			if source := profile.Source; source != nil {
				if image := source.ManagedImage; image != nil {
					d.Set("managed_image_id", image.ID)
				}
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSharedImageVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	imageVersion := id.Path["versions"]
	imageName := id.Path["images"]
	galleryName := id.Path["galleries"]
	resourceGroup := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroup, galleryName, imageName, imageVersion)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
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

		output := compute.TargetRegion{
			Name:                 utils.String(name),
			RegionalReplicaCount: utils.Int32(int32(regionalReplicaCount)),
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

			results = append(results, output)
		}
	}

	return results
}
