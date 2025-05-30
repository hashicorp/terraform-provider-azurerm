// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceSharedImageVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimageversions.NewImageVersionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string), d.Get("name").(string))
	sortBySemVer := d.Get("sort_versions_by_semver").(bool)

	image, err := obtainImage(client, ctx, id, sortBySemVer)
	if err != nil {
		return err
	}

	name := ""
	if image.Name != nil {
		name = *image.Name
	}

	exactId := galleryimageversions.NewImageVersionID(subscriptionId, id.ResourceGroupName, id.GalleryName, id.ImageName, name)
	d.SetId(exactId.ID())
	d.Set("name", name)
	d.Set("image_name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("sort_versions_by_semver", sortBySemVer)

	d.Set("location", location.Normalize(image.Location))

	if props := image.Properties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			d.Set("exclude_from_latest", profile.ExcludeFromLatest)

			if err := d.Set("target_region", flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)); err != nil {
				return fmt.Errorf("setting `target_region`: %+v", err)
			}
		}

		if source := props.StorageProfile.Source; source != nil {
			d.Set("managed_image_id", source.Id)
			osDiskSnapShotID := ""
			if props.StorageProfile.OsDiskImage != nil && props.StorageProfile.OsDiskImage.Source != nil && props.StorageProfile.OsDiskImage.Source.Id != nil {
				osDiskSnapShotID = *props.StorageProfile.OsDiskImage.Source.Id
			}
			d.Set("os_disk_snapshot_id", osDiskSnapShotID)

			osDiskImageSize := 0
			if props.StorageProfile.OsDiskImage != nil && props.StorageProfile.OsDiskImage.SizeInGB != nil {
				osDiskImageSize = int(*props.StorageProfile.OsDiskImage.SizeInGB)
			}
			d.Set("os_disk_image_size_gb", osDiskImageSize)
		}
		return tags.FlattenAndSet(d, image.Tags)
	}
	return nil
}

func obtainImage(client *galleryimageversions.GalleryImageVersionsClient, ctx context.Context, id galleryimageversions.ImageVersionId, sortBySemVer bool) (*galleryimageversions.GalleryImageVersion, error) {
	galleryImageId := galleryimageversions.NewGalleryImageID(id.SubscriptionId, id.ResourceGroupName, id.GalleryName, id.ImageName)

	notFoundError := fmt.Errorf("a version was not found for %s", galleryImageId)

	switch id.VersionName {
	case "latest":
		resp, err := client.ListByGalleryImageComplete(ctx, galleryImageId)
		if err != nil {
			if response.WasNotFound(resp.LatestHttpResponse) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("retrieving `latest` versions for %s: %+v", galleryImageId, err)
		}

		images := resp.Items

		// the last image in the list is the latest version
		if len(images) > 0 {
			if sortBySemVer {
				var errs []error
				images, errs = sortSharedImageVersions(images)
				if len(errs) > 0 {
					return nil, fmt.Errorf("parsing version(s): %v", errs)
				}
			}

			for i := len(images) - 1; i >= 0; i-- {
				if prop := images[i].Properties; prop == nil || prop.PublishingProfile == nil || prop.PublishingProfile.ExcludeFromLatest == nil || !*prop.PublishingProfile.ExcludeFromLatest {
					return &(images[i]), nil
				}
			}
		}
		return nil, notFoundError

	case "recent":
		resp, err := client.ListByGalleryImageComplete(ctx, galleryImageId)
		if err != nil {
			if response.WasNotFound(resp.LatestHttpResponse) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("retrieving `recent` versions for %s: %+v", galleryImageId, err)
		}
		var image *galleryimageversions.GalleryImageVersion
		var recentDate *time.Time
		// compare dates until we find the image that was updated most recently
		for _, currentImage := range resp.Items {
			if profile := currentImage.Properties.PublishingProfile; profile != nil {
				if profile.PublishedDate != nil {
					publishedDate, err := time.Parse(time.RFC3339, *profile.PublishedDate)
					if err != nil {
						return nil, fmt.Errorf("parsing published date for %s: %+v", galleryImageId, err)
					}
					if recentDate == nil || publishedDate.After(*recentDate) {
						recentDate = pointer.To(publishedDate)
						image = pointer.To(currentImage)
					}
				}
			}
		}

		if image != nil {
			return image, nil
		}

		return nil, notFoundError

	default:
		image, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(image.HttpResponse) {
				return nil, notFoundError
			}
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if image.Model == nil {
			return nil, fmt.Errorf("model is nil for %s", id)
		}

		return image.Model, nil
	}
}

func flattenSharedImageVersionDataSourceTargetRegions(input *[]galleryimageversions.TargetRegion) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]interface{})

			output["name"] = location.Normalize(v.Name)

			if v.RegionalReplicaCount != nil {
				output["regional_replica_count"] = int(*v.RegionalReplicaCount)
			}

			if v.StorageAccountType != nil {
				output["storage_account_type"] = string(*v.StorageAccountType)
			}

			results = append(results, output)
		}
	}

	return results
}
