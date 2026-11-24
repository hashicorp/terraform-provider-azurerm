// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Image.",
		Read:        dataSourceImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name_regex": {
				Description:  "Regex pattern of the image to match.",
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ExactlyOneOf: []string{"name", "name_regex"},
			},
			"sort_descending": {
				Description: "By default when matching by regex, images are sorted by name in ascending order and the first match is chosen, to sort descending, set this flag.",
				Type:        pluginsdk.TypeBool,
				Optional:    true,
				Default:     false,
			},

			"name": {
				Description:  "The name of the Image.",
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "name_regex"},
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"zone_resilient": {
				Description: "Whether zone resiliency is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"os_disk": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"blob_uri": {
							Description: "The URI in Azure storage of the blob used to create the image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"caching": {
							Description: "The caching mode for the OS Disk, such as `ReadWrite`, `ReadOnly`, or `None`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"managed_disk_id": {
							Description: "The ID of the Managed Disk used as the OS Disk Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"os_state": {
							Description: "The State of the OS used in the Image, such as `Generalized`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"os_type": {
							Description: "The type of Operating System used on the OS Disk. Such as `Linux` or `Windows`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"size_gb": {
							Description: "The size of the OS Disk in GB.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
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
							Description: "The URI in Azure storage of the blob used to create the image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"caching": {
							Description: "The caching mode for the Data Disk, such as `ReadWrite`, `ReadOnly`, or `None`.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"lun": {
							Description: "The logical unit number of the data disk.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"managed_disk_id": {
							Description: "The ID of the Managed Disk used as the Data Disk Image.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"size_gb": {
							Description: "The size of this Data Disk in GB.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceImageRead(d *pluginsdk.ResourceData, meta any) error {
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

func flattenImageDataSourceOSDisk(input *images.ImageOSDisk) []any {
	output := make([]any, 0)

	if input != nil {
		blobUri := pointer.From(input.BlobUri)
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
		output = append(output, map[string]any{
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

func flattenImageDataSourceDataDisks(input *[]images.ImageDataDisk) []any {
	output := make([]any, 0)

	if input != nil {
		for _, disk := range *input {
			blobUri := pointer.From(disk.BlobUri)
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
			output = append(output, map[string]any{
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
