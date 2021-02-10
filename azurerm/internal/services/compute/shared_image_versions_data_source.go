package compute

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSharedImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSharedImageVersionsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags_filter": tags.Schema(),

			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": azure.SchemaLocationForDataSource(),

						"managed_image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"target_region": {
							Type:     schema.TypeList,
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

									"storage_account_type": {
										Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func dataSourceSharedImageVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	filterTags := tags.Expand(d.Get("tags_filter").(map[string]interface{}))

	resp, err := client.ListByGalleryImageComplete(ctx, resourceGroup, galleryName, imageName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			return fmt.Errorf("No Versions were found for Shared Image %q / Gallery %q / Resource Group %q", imageName, galleryName, resourceGroup)
		}
		return fmt.Errorf("retrieving Shared Image Versions (Image %q / Gallery %q / Resource Group %q): %+v", imageName, galleryName, resourceGroup, err)
	}

	images := make([]compute.GalleryImageVersion, 0)
	for resp.NotDone() {
		images = append(images, resp.Value())
		if err := resp.NextWithContext(ctx); err != nil {
			return fmt.Errorf("listing next page of images for Shared Image %q / Gallery %q / Resource Group %q: %+v", imageName, galleryName, resourceGroup, err)
		}
	}

	flattenedImages := flattenSharedImageVersions(images, filterTags)
	if len(flattenedImages) == 0 {
		return fmt.Errorf("unable to find any images")
	}

	d.SetId(time.Now().UTC().String())

	d.Set("image_name", imageName)
	d.Set("gallery_name", galleryName)
	d.Set("resource_group_name", resourceGroup)

	if err := d.Set("images", flattenedImages); err != nil {
		return fmt.Errorf("setting `images`: %+v", err)
	}

	return nil
}

func flattenSharedImageVersions(input []compute.GalleryImageVersion, filterTags map[string]*string) []interface{} {
	results := make([]interface{}, 0)

	for _, imageVersion := range input {
		flattenedIPAddress := flattenSharedImageVersion(imageVersion)
		found := true
		// Loop through our filter tags and see if they match
		for k, v := range filterTags {
			if v != nil {
				// If the tags don't match, return false
				if imageVersion.Tags[k] == nil || *v != *imageVersion.Tags[k] {
					found = false
				}
			}
		}

		if found {
			results = append(results, flattenedIPAddress)
		}
	}

	return results
}

func flattenSharedImageVersion(input compute.GalleryImageVersion) map[string]interface{} {
	output := make(map[string]interface{})

	output["name"] = input.Name

	if location := input.Location; location != nil {
		output["location"] = azure.NormalizeLocation(*location)
	}

	if props := input.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			output["exclude_from_latest"] = profile.ExcludeFromLatest
			output["target_region"] = flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)
		}

		if profile := props.StorageProfile; profile != nil {
			if source := profile.Source; source != nil {
				output["managed_image_id"] = source.ID
			}
		}
	}

	output["tags"] = tags.Flatten(input.Tags)

	return output
}
