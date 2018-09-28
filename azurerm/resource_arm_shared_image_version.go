package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"managed_image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"regions": {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				StateFunc:        azureRMNormalizeLocation,
				DiffSuppressFunc: azureRMSuppressLocationDiff,
			},

			"exclude_from_latest": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSharedImageVersionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).galleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	imageVersion := d.Get("name").(string)
	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	managedImageId := d.Get("managed_image_id").(string)
	excludeFromLatest := d.Get("exclude_from_latest").(bool)

	regions := expandSharedImageVersionRegions(d)
	tags := d.Get("tags").(map[string]interface{})

	version := compute.GalleryImageVersion{
		Location: utils.String(location),
		GalleryImageVersionProperties: &compute.GalleryImageVersionProperties{
			PublishingProfile: &compute.GalleryImageVersionPublishingProfile{
				ExcludeFromLatest: utils.Bool(excludeFromLatest),
				Regions:           regions,
				Source: &compute.GalleryArtifactSource{
					ManagedImage: &compute.ManagedArtifact{
						ID: utils.String(managedImageId),
					},
				},
			},
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, galleryName, imageName, imageVersion, version)
	if err != nil {
		return fmt.Errorf("Error creating Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the creation of Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmSharedImageVersionRead(d, meta)
}

func resourceArmSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).galleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.GalleryImageVersionProperties; props != nil {
		// `targetRegions` is returned in the API Response but isn't exposed there.
		// TODO: replace this once this fields exposed
		// BUG: https://github.com/Azure/azure-sdk-for-go/issues/2855
		if status := props.ReplicationStatus; status != nil {
			flattenedRegions := flattenSharedImageVersionRegions(status.Summary)
			if err := d.Set("regions", flattenedRegions); err != nil {
				return fmt.Errorf("Error flattening `regions`: %+v", err)
			}
		}

		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			if source := profile.Source; source != nil {
				if image := source.ManagedImage; image != nil {
					d.Set("managed_image_id", image.ID)
				}
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSharedImageVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).galleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
		}
	}

	return nil
}

func expandSharedImageVersionRegions(d *schema.ResourceData) *[]string {
	vs := d.Get("regions").(*schema.Set)
	output := make([]string, 0)

	for _, v := range vs.List() {
		output = append(output, v.(string))
	}

	return &output
}

func flattenSharedImageVersionRegions(input *[]compute.RegionalReplicationStatus) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			if v.Region != nil {
				output = append(output, azureRMNormalizeLocation(*v.Region))
			}
		}
	}

	return output
}
