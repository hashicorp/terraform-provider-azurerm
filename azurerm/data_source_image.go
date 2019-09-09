package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmImageRead,
		Schema: map[string]*schema.Schema{

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.ValidateRegexp,
				ConflictsWith: []string{"name"},
			},
			"sort_descending": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name_regex"},
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"zone_resilient": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"os_disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"caching": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"managed_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"data_disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"caching": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lun": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"managed_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.ImagesClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)

	name := d.Get("name").(string)
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if name == "" && !nameRegexOk {
		return fmt.Errorf("[ERROR] either name or name_regex is required")
	}

	var img compute.Image

	if !nameRegexOk {
		var err error
		if img, err = client.Get(ctx, resGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(img.Response) {
				return fmt.Errorf("Error: Image %q (Resource Group %q) was not found", name, resGroup)
			}
			return fmt.Errorf("[ERROR] Error making Read request on Azure Image %q (resource group %q): %+v", name, resGroup, err)
		}
	} else {
		r := regexp.MustCompile(nameRegex.(string))

		list := make([]compute.Image, 0)
		resp, err := client.ListByResourceGroupComplete(ctx, resGroup)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response().Response) {
				return fmt.Errorf("Error: Unable to list images for Resource Group %q", resGroup)
			}
			return fmt.Errorf("[ERROR] Error getting list of images (resource group %q): %+v", resGroup, err)
		}

		for resp.NotDone() {
			img = resp.Value()
			if r.Match(([]byte)(*img.Name)) {
				list = append(list, img)
			}
			err = resp.Next()

			if err != nil {
				return err
			}
		}

		if len(list) < 1 {
			d.SetId("")
			return nil
		}

		if len(list) > 1 {
			desc := d.Get("sort_descending").(bool)
			log.Printf("[DEBUG] arm_image - multiple results found and `sort_descending` is set to: %t", desc)

			sort.Slice(list, func(i, j int) bool {
				return (!desc && *list[i].Name < *list[j].Name) ||
					(desc && *list[i].Name > *list[j].Name)
			})
		}
		img = list[0]

	}

	d.SetId(*img.ID)
	d.Set("name", img.Name)
	d.Set("resource_group_name", resGroup)
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
