// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
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

			"tags_filter": tags.Schema(),

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

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceSharedImageVersionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	d.SetId(fmt.Sprintf("%s-%s-%s", imageName, galleryName, resourceGroup))

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
		flattenedImageVersion := flattenSharedImageVersion(imageVersion)
		found := true
		// Loop through our filter tags and see if they match
		for k, v := range filterTags {
			if v != nil {
				// If the tags don't match, return false
				if imageVersion.Tags[k] == nil || *v != *imageVersion.Tags[k] {
					found = false
					break
				}
			}
		}

		if found {
			results = append(results, flattenedImageVersion)
		}
	}

	return results
}

func flattenSharedImageVersion(input compute.GalleryImageVersion) map[string]interface{} {
	output := make(map[string]interface{})

	output["id"] = input.ID
	output["name"] = input.Name
	output["location"] = location.NormalizeNilable(input.Location)

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
