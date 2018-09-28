package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSharedImageVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSharedImageVersionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageVersionName,
			},

			"gallery_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"image_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"managed_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"regions": {
				Type:      schema.TypeSet,
				Computed:  true,
				Elem:      &schema.Schema{Type: schema.TypeString},
				StateFunc: azureRMNormalizeLocation,
			},

			"exclude_from_latest": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).galleryImageVersionsClient
	ctx := meta.(*ArmClient).StopContext

	imageVersion := d.Get("name").(string)
	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, galleryName, imageName, imageVersion, compute.ReplicationStatusTypesReplicationStatus)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image Version %q (Image %q / Gallery %q / Resource Group %q) was not found - removing from state", imageVersion, imageName, galleryName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", imageVersion, imageName, galleryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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
			flattenedRegions := flattenSharedImageVersionDataSourceRegions(status.Summary)
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

func flattenSharedImageVersionDataSourceRegions(input *[]compute.RegionalReplicationStatus) []interface{} {
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
