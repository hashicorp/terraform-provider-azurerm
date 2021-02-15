package compute

import (
	"context"
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

func dataSourceSharedImageVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSharedImageVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

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

			"os_disk_snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_disk_image_size_gb": {
				Type:     schema.TypeInt,
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
	}
}

func dataSourceSharedImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	imageVersion := d.Get("name").(string)
	imageName := d.Get("image_name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	image, err := obtainImage(client, ctx, resourceGroup, galleryName, imageName, imageVersion)
	if err != nil {
		return err
	}

	d.SetId(*image.ID)
	d.Set("name", image.Name)
	d.Set("image_name", imageName)
	d.Set("gallery_name", galleryName)
	d.Set("resource_group_name", resourceGroup)

	if location := image.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := image.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			if err := d.Set("target_region", flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)); err != nil {
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

			osDiskImageSize := 0
			if profile.OsDiskImage != nil && profile.OsDiskImage.SizeInGB != nil {
				osDiskImageSize = int(*profile.OsDiskImage.SizeInGB)
			}
			d.Set("os_disk_image_size_gb", osDiskImageSize)
		}
	}

	return tags.FlattenAndSet(d, image.Tags)
}

func obtainImage(client *compute.GalleryImageVersionsClient, ctx context.Context, resourceGroup string, galleryName string, galleryImageName string, galleryImageVersionName string) (*compute.GalleryImageVersion, error) {
	notFoundError := fmt.Errorf("A Version was not found for Shared Image %q / Gallery %q / Resource Group %q", galleryImageName, galleryName, resourceGroup)

	switch galleryImageVersionName {
	case "latest":
		images, err := client.ListByGalleryImage(ctx, resourceGroup, galleryName, galleryImageName)
		if err != nil {
			if utils.ResponseWasNotFound(images.Response().Response) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("retrieving Shared Image Versions (Image %q / Gallery %q / Resource Group %q): %+v", galleryImageName, galleryName, resourceGroup, err)
		}

		// the last image in the list is the latest version
		if len(images.Values()) > 0 {
			image := images.Values()[len(images.Values())-1]
			return &image, nil
		}

		return nil, notFoundError

	case "recent":
		images, err := client.ListByGalleryImage(ctx, resourceGroup, galleryName, galleryImageName)
		if err != nil {
			if utils.ResponseWasNotFound(images.Response().Response) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("retrieving Shared Image Versions (Image %q / Gallery %q / Resource Group %q): %+v", galleryImageName, galleryName, resourceGroup, err)
		}
		var image *compute.GalleryImageVersion
		var recentDate *time.Time
		// compare dates until we find the image that was updated most recently
		for _, currImage := range images.Values() {
			if profile := currImage.PublishingProfile; profile != nil {
				if profile.PublishedDate != nil && (recentDate == nil || profile.PublishedDate.Time.After(*recentDate)) {
					recentDate = &profile.PublishedDate.Time
					image = &currImage
				}
			}
		}

		if image != nil {
			return image, nil
		}

		return nil, notFoundError

	default:
		image, err := client.Get(ctx, resourceGroup, galleryName, galleryImageName, galleryImageVersionName, compute.ReplicationStatusTypesReplicationStatus)
		if err != nil {
			if utils.ResponseWasNotFound(image.Response) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("Error retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", galleryImageVersionName, galleryImageName, galleryName, resourceGroup, err)
		}

		return &image, nil
	}
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

			output["storage_account_type"] = string(v.StorageAccountType)

			results = append(results, output)
		}
	}

	return results
}
