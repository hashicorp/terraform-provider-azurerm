// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name_regex": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ExactlyOneOf: []string{"name", "name_regex"},
			},
			"sort_descending": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "name_regex"},
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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
	}
}

func dataSourceImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	nameRegex, nameRegexOk := d.GetOk("name_regex")

	var id images.ImageId
	var model images.Image

	if !nameRegexOk {
		var err error
		id = images.NewImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
		result, err := client.Get(ctx, id, images.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(result.HttpResponse) {
				return fmt.Errorf("%s was not found", id)
			}
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if result.Model == nil {
			return fmt.Errorf("%s was not found", id)
		}
		model = *result.Model
	} else {
		r := regexp.MustCompile(nameRegex.(string))
		resourceGroupId := commonids.NewResourceGroupID(subscriptionId, d.Get("resource_group_name").(string))
		listResponse, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
		if err != nil {
			return fmt.Errorf("listing Images within %s: %+v", resourceGroupId, err)
		}
		results := make([]images.Image, 0)
		for _, v := range listResponse.Items {
			if r.Match(([]byte)(*v.Name)) {
				results = append(results, v)
			}
		}

		if len(results) == 0 {
			return fmt.Errorf("no Images were found within %s", resourceGroupId)
		}
		if len(results) > 1 {
			desc := d.Get("sort_descending").(bool)
			log.Printf("[DEBUG] Image - multiple results found and `sort_descending` is set to: %t", desc)

			sort.Slice(results, func(i, j int) bool {
				return (!desc && *results[i].Name < *results[j].Name) ||
					(desc && *results[i].Name > *results[j].Name)
			})
		}
		model = results[0]
		if model.Name == nil {
			return fmt.Errorf("image name is null for the first Image in %s: %+v", resourceGroupId, model)
		}

		id = images.NewImageID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroupName, *model.Name)
	}

	d.SetId(id.ID())
	d.Set("name", id.ImageName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.Normalize(model.Location))
	if props := model.Properties; props != nil {
		if profile := props.StorageProfile; profile != nil {
			if err := d.Set("os_disk", flattenImageDataSourceOSDisk(profile.OsDisk)); err != nil {
				return fmt.Errorf("setting `os_disk`: %+v", err)
			}

			if err := d.Set("data_disk", flattenImageDataSourceDataDisks(profile.DataDisks)); err != nil {
				return fmt.Errorf("setting `data_disk`: %+v", err)
			}
			d.Set("zone_resilient", profile.ZoneResilient)
		}
	}

	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return fmt.Errorf("setting `tags`: %+v", err)
	}

	return nil
}

func flattenImageDataSourceOSDisk(input *images.ImageOSDisk) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		blobUri := ""
		if uri := input.BlobUri; uri != nil {
			blobUri = *uri
		}
		caching := ""
		if input.Caching != nil {
			caching = string(*input.Caching)
		}
		diskSizeGB := 0
		if input.DiskSizeGB != nil {
			diskSizeGB = int(*input.DiskSizeGB)
		}
		managedDiskId := ""
		if disk := input.ManagedDisk; disk != nil && disk.Id != nil {
			managedDiskId = *disk.Id
		}
		output = append(output, map[string]interface{}{
			"blob_uri":        blobUri,
			"caching":         caching,
			"managed_disk_id": managedDiskId,
			"os_type":         string(input.OsType),
			"os_state":        string(input.OsState),
			"size_gb":         diskSizeGB,
		})
	}

	return output
}

func flattenImageDataSourceDataDisks(input *[]images.ImageDataDisk) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, disk := range *input {
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

	return output
}
