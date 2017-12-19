package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

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

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).imageClient

	resGroup := d.Get("resource_group_name").(string)

	name := d.Get("name").(string)
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if name == "" && !nameRegexOk {
		return fmt.Errorf("[ERROR] either name or name_regex is required")
	}

	var img compute.Image

	if !nameRegexOk {
		var err error
		if img, err = client.Get(resGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(img.Response) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error making Read request on Azure Image %q (resource group %q): %+v", name, resGroup, err)
		}
	} else {
		r := regexp.MustCompile(nameRegex.(string))

		list := []compute.Image{}
		resp, err := client.ListByResourceGroup(resGroup)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error getting list of images (resource group %q): %+v", resGroup, err)
		}
		for _, ri := range *resp.Value {
			if r.Match(([]byte)(*ri.Name)) {
				list = append(list, ri)
			}
		}
		for resp.NextLink != nil && *resp.NextLink != "" {
			resp, err = client.ListByResourceGroupNextResults(resp)
			if err != nil {
				return fmt.Errorf("[ERROR] Unable to query next results for image list (resource group %q): %+v", resGroup, err)
			}
			for _, ri := range *resp.Value {
				if r.Match(([]byte)(*ri.Name)) {
					list = append(list, ri)
				}
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
	d.Set("location", azureRMNormalizeLocation(*img.Location))

	if profile := img.StorageProfile; profile != nil {
		if disk := profile.OsDisk; disk != nil {
			if err := d.Set("os_disk", flattenAzureRmImageOSDisk(d, disk)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image OS Disk error: %+v", err)
			}
		}

		if disks := img.StorageProfile.DataDisks; disks != nil {
			if err := d.Set("data_disk", flattenAzureRmImageDataDisks(d, disks)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image Data Disks error: %+v", err)
			}
		}
	}

	flattenAndSetTags(d, img.Tags)

	return nil
}
