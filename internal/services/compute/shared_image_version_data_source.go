package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSharedImageVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSharedImageVersionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageVersionName,
			},

			"gallery_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"image_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"managed_image_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_disk_snapshot_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_disk_image_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"sort_versions_by_semver": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"target_region": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"regional_replica_count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"storage_account_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"exclude_from_latest": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceSharedImageVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedImageVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string), d.Get("name").(string))
	sortBySemVer := d.Get("sort_versions_by_semver").(bool)

	image, err := obtainImage(client, ctx, id.ResourceGroup, id.GalleryName, id.ImageName, id.VersionName, sortBySemVer)
	if err != nil {
		return err
	}

	name := ""
	if image.Name != nil {
		name = *image.Name
	}

	exactId := parse.NewSharedImageVersionID(subscriptionId, id.ResourceGroup, id.GalleryName, id.ImageName, name)
	d.SetId(exactId.ID())
	d.Set("name", name)
	d.Set("image_name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sort_versions_by_semver", sortBySemVer)

	d.Set("location", location.NormalizeNilable(image.Location))

	if props := image.GalleryImageVersionProperties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			if err := d.Set("target_region", flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)); err != nil {
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

			osDiskImageSize := 0
			if profile.OsDiskImage != nil && profile.OsDiskImage.SizeInGB != nil {
				osDiskImageSize = int(*profile.OsDiskImage.SizeInGB)
			}
			d.Set("os_disk_image_size_gb", osDiskImageSize)
		}
	}

	return tags.FlattenAndSet(d, image.Tags)
}

func obtainImage(client *compute.GalleryImageVersionsClient, ctx context.Context, resourceGroup string, galleryName string, galleryImageName string, galleryImageVersionName string, sortBySemVer bool) (*compute.GalleryImageVersion, error) {
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
			values := images.Values()
			var errs []error
			if sortBySemVer {
				values, errs = sortSharedImageVersions(values)
				if len(errs) > 0 {
					return nil, fmt.Errorf("parsing version(s): %v", errs)
				}
			}
			image := values[len(values)-1]
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
			return nil, fmt.Errorf("retrieving Shared Image Version %q (Image %q / Gallery %q / Resource Group %q): %+v", galleryImageVersionName, galleryImageName, galleryName, resourceGroup, err)
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
