package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"managed_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_region": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"regional_replica_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"exclude_from_latest": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleryImageVersionsClient
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			flattenedRegions := flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)
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

func flattenSharedImageVersionDataSourceTargetRegions(input *[]compute.TargetRegion) []interface{} {
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
