package compute

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	resourceGroup := d.Get("resource_group_name").(string)

	var img compute.Image

	if !nameRegexOk {
		var err error
		if img, err = client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(img.Response) {
				return fmt.Errorf("image %q (Resource Group: %s) was not found", name, resourceGroup)
			}
			return fmt.Errorf("making Read request on image %q (Resource Group: %s): %s", name, resourceGroup, err)
		}
	} else {
		r := regexp.MustCompile(nameRegex.(string))

		list := make([]compute.Image, 0)
		resp, err := client.ListByResourceGroupComplete(ctx, resourceGroup)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response().Response) {
				return fmt.Errorf("no Images were found for Resource Group %q", resourceGroup)
			}
			return fmt.Errorf("getting list of images (resource group %q): %+v", resourceGroup, err)
		}

		for resp.NotDone() {
			img = resp.Value()
			if r.Match(([]byte)(*img.Name)) {
				list = append(list, img)
			}
			err = resp.NextWithContext(ctx)

			if err != nil {
				return err
			}
		}

		if 1 > len(list) {
			return fmt.Errorf("no Images were found for Resource Group %q", resourceGroup)
		}

		if len(list) > 1 {
			desc := d.Get("sort_descending").(bool)
			log.Printf("[DEBUG] Image - multiple results found and `sort_descending` is set to: %t", desc)

			sort.Slice(list, func(i, j int) bool {
				return (!desc && *list[i].Name < *list[j].Name) ||
					(desc && *list[i].Name > *list[j].Name)
			})
		}
		img = list[0]
	}

	if img.Name == nil {
		return fmt.Errorf("image name is empty in Resource Group %s", resourceGroup)
	}

	id := parse.NewImageID(subscriptionId, resourceGroup, *img.Name)

	d.SetId(id.ID())
	d.Set("name", img.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := img.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if profile := img.StorageProfile; profile != nil {
		if disk := profile.OsDisk; disk != nil {
			if err := d.Set("os_disk", flattenAzureRmImageOSDisk(disk)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image OS Disk error: %+v", err)
			}
		}

		if disks := profile.DataDisks; disks != nil {
			if err := d.Set("data_disk", flattenAzureRmImageDataDisks(disks)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image Data Disks error: %+v", err)
			}
		}

		d.Set("zone_resilient", profile.ZoneResilient)
	}

	return tags.FlattenAndSet(d, img.Tags)
}
