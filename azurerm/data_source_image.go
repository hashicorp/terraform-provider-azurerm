package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmImageRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
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

	resp, err := client.Get(resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Image %q (resource group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if profile := resp.StorageProfile; profile != nil {
		if disk := profile.OsDisk; disk != nil {
			if err := d.Set("os_disk", flattenAzureRmImageOSDisk(d, disk)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image OS Disk error: %+v", err)
			}
		}

		if disks := resp.StorageProfile.DataDisks; disks != nil {
			if err := d.Set("data_disk", flattenAzureRmImageDataDisks(d, disks)); err != nil {
				return fmt.Errorf("[DEBUG] Error setting AzureRM Image Data Disks error: %+v", err)
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
