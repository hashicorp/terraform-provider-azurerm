// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

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

func dataSourceSharedImageVersions() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSharedImageVersionsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"tags_filter": commonschema.Tags(),

			"images": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": commonschema.LocationComputed(),

						"managed_image_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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

						"tags": commonschema.TagsDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceSharedImageVersionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImageVersionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := galleryimageversions.NewGalleryImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("image_name").(string))
	filterTags := tags.Expand(d.Get("tags_filter").(map[string]interface{}))

	resp, err := client.ListByGalleryImageComplete(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.LatestHttpResponse) {
			return fmt.Errorf("no versions were found for %s", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	flattenedImages := flattenSharedImageVersions(resp.Items, filterTags)
	if len(flattenedImages) == 0 {
		return fmt.Errorf("unable to find any images")
	}

	d.SetId(fmt.Sprintf("%s-%s-%s", id.ImageName, id.GalleryName, id.ResourceGroupName))

	d.Set("image_name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if err := d.Set("images", flattenedImages); err != nil {
		return fmt.Errorf("setting `images`: %+v", err)
	}

	return nil
}

func flattenSharedImageVersions(input []galleryimageversions.GalleryImageVersion, filterTags *map[string]string) []interface{} {
	results := make([]interface{}, 0)

	if len(input) == 0 {
		return results
	}

	for _, imageVersion := range input {
		flattenedImageVersion := flattenSharedImageVersion(imageVersion)
		found := true
		// Loop through our filter tags and see if they match
		if filterTags != nil {
			for k, v := range *filterTags {
				if imageVersion.Tags != nil {
					// If the tags don't match, return false
					if v != (*imageVersion.Tags)[k] {
						found = false
						break
					}
				}
			}
		}

		if found {
			results = append(results, flattenedImageVersion)
		}
	}

	return results
}

func flattenSharedImageVersion(input galleryimageversions.GalleryImageVersion) map[string]interface{} {
	output := make(map[string]interface{})

	output["id"] = input.Id
	output["name"] = input.Name
	output["location"] = location.Normalize(input.Location)

	if props := input.Properties; props != nil {
		if profile := props.PublishingProfile; profile != nil {
			output["exclude_from_latest"] = profile.ExcludeFromLatest
			output["target_region"] = flattenSharedImageVersionDataSourceTargetRegions(profile.TargetRegions)
		}

		if source := props.StorageProfile.Source; source != nil {
			output["managed_image_id"] = source.Id
		}
	}

	output["tags"] = tags.Flatten(input.Tags)

	return output
}
