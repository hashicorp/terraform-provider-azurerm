// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceImages() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceImagesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"tags_filter": commonschema.Tags(),

			"images": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": commonschema.LocationComputed(),

						"zone_resilient": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"os_disk": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"blob_uri": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"caching": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"managed_disk_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"os_state": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"size_gb": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},

									"disk_encryption_set_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"data_disk": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"blob_uri": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"caching": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"lun": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"managed_disk_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"size_gb": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"tags": commonschema.TagsDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceImagesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	filterTags := tags.Expand(d.Get("tags_filter").(map[string]interface{}))

	resourceGroupId := commonids.NewResourceGroupID(subscriptionId, d.Get("resource_group_name").(string))
	resp, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
	if err != nil {
		return fmt.Errorf("retrieving Images within %s: %+v", resourceGroupId, err)
	}

	virtualMachineImages := resp.Items
	if filterTags != nil && len(*filterTags) > 0 {
		virtualMachineImages = filterToImagesMatchingTags(virtualMachineImages, *filterTags)
	}
	if len(virtualMachineImages) == 0 {
		return fmt.Errorf("no images were found that match the specified tags")
	}
	flattenedImages := flattenImages(virtualMachineImages)
	if err := d.Set("images", flattenedImages); err != nil {
		return fmt.Errorf("setting `images`: %+v", err)
	}

	resourceId := resourceIdForImagesDataSource(resourceGroupId, *filterTags)
	d.SetId(resourceId)

	d.Set("resource_group_name", resourceGroupId.ResourceGroupName)

	return nil
}

func resourceIdForImagesDataSource(resourceGroupId commonids.ResourceGroupId, filterTags map[string]string) string {
	tagsId := ""
	tagKeys := make([]string, 0)
	for key := range filterTags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)
	for _, key := range tagKeys {
		value := ""
		if v, ok := filterTags[key]; ok {
			value = v
		}
		tagsId += fmt.Sprintf("[%s:%s]", key, value)
	}
	if tagsId == "" {
		tagsId = "[]"
	}
	return fmt.Sprintf("resourceGroups/%s/tags/%s/images", resourceGroupId.ResourceGroupName, tagsId)
}

func flattenImages(input []images.Image) []interface{} {
	output := make([]interface{}, 0)
	for _, item := range input {
		output = append(output, flattenImage(item))
	}
	return output
}

func filterToImagesMatchingTags(input []images.Image, filterTags map[string]string) []images.Image {
	output := make([]images.Image, 0)

	for _, item := range input {
		tagsMatch := true
		if item.Tags == nil {
			tagsMatch = false
		} else {
			for tagKey, tagValue := range filterTags {
				otherVal, exists := (*item.Tags)[tagKey]
				if !exists || tagValue != otherVal {
					tagsMatch = false
					break
				}
			}
		}

		if tagsMatch {
			output = append(output, item)
		}
	}

	return output
}

func flattenImage(input images.Image) map[string]interface{} {
	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	zoneResilient := false
	osDisk := make([]interface{}, 0)
	dataDisks := make([]interface{}, 0)
	if props := input.Properties; props != nil {
		osDisk = flattenImagesOSDisk(props.StorageProfile)
		dataDisks = flattenImagesDataDisks(props.StorageProfile)

		if props.StorageProfile != nil && props.StorageProfile.ZoneResilient != nil {
			zoneResilient = *props.StorageProfile.ZoneResilient
		}
	}

	return map[string]interface{}{
		"location":       location.Normalize(input.Location),
		"data_disk":      dataDisks,
		"name":           name,
		"os_disk":        osDisk,
		"tags":           tags.Flatten(input.Tags),
		"zone_resilient": zoneResilient,
	}
}

func flattenImagesOSDisk(input *images.ImageStorageProfile) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v := input.OsDisk; v != nil {
			blobUri := ""
			if uri := v.BlobUri; uri != nil {
				blobUri = *uri
			}
			caching := ""
			if v.Caching != nil {
				caching = string(*v.Caching)
			}
			diskSizeGB := 0
			if v.DiskSizeGB != nil {
				diskSizeGB = int(*v.DiskSizeGB)
			}
			managedDiskId := ""
			if disk := v.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
			diskEncryptionSetId := ""
			if set := v.DiskEncryptionSet; set != nil && set.Id != nil {
				diskEncryptionSetId = *set.Id
			}
			output = append(output, map[string]interface{}{
				"blob_uri":               blobUri,
				"caching":                caching,
				"managed_disk_id":        managedDiskId,
				"os_type":                string(v.OsType),
				"os_state":               string(v.OsState),
				"size_gb":                diskSizeGB,
				"disk_encryption_set_id": diskEncryptionSetId,
			})
		}
	}

	return output
}

func flattenImagesDataDisks(input *images.ImageStorageProfile) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v := input.DataDisks; v != nil {
			for _, disk := range *input.DataDisks {
				blobUri := ""
				if disk.BlobUri != nil {
					blobUri = *disk.BlobUri
				}
				caching := ""
				if disk.Caching != nil {
					caching = string(*disk.Caching)
				}
				diskSizeGb := 0
				if disk.DiskSizeGB != nil {
					diskSizeGb = int(*disk.DiskSizeGB)
				}
				managedDiskId := ""
				if disk.ManagedDisk != nil && disk.ManagedDisk.Id != nil {
					managedDiskId = *disk.ManagedDisk.Id
				}
				output = append(output, map[string]interface{}{
					"blob_uri":        blobUri,
					"caching":         caching,
					"lun":             int(disk.Lun),
					"managed_disk_id": managedDiskId,
					"size_gb":         diskSizeGb,
				})
			}
		}
	}

	return output
}
