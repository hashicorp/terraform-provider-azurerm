package compute

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

func dataSourceArmSharedImageVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSharedImageVersionsRead,

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

			"location": azure.SchemaLocationForDataSource(),

			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

func dataSourceArmSharedImageVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.ListByGalleryImage(ctx, resourceGroup, galleryName, imageName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			log.Printf("[DEBUG] Shared Image Versions (Image %q / Gallery %q / Resource Group %q) was not found - removing from state", imageName, galleryName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Shared Image Versions (Image %q / Gallery %q / Resource Group %q): %+v", imageName, galleryName, resourceGroup, err)
	}
	return fmt.Errorf("%+v", resp.Response())

	d.SetId(time.Now().UTC().String())

	d.Set("image_name", imageName)
	d.Set("gallery_name", galleryName)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Response().Location; location != nil {
		d.Set("location", azure.NormalizeLocation(location))
	}
	/*
		if props := resp.GalleryImageVersionProperties; props != nil {
			if profile := props.PublishingProfile; profile != nil {
				d.Set("exclude_from_latest", profile.ExcludeFromLatest)

				flattenedRegions := flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)
				if err := d.Set("target_region", flattenedRegions); err != nil {
					return fmt.Errorf("Error setting `target_region`: %+v", err)
				}
			}

			if profile := props.StorageProfile; profile != nil {
				if source := profile.Source; source != nil {
					d.Set("managed_image_id", source.ID)
				}
			}
		}

		return tags.FlattenAndSet(d, resp.Tags)*/
	return nil
}
