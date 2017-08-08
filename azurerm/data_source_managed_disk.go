package azurerm

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmManagedDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagedDiskRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"storage_account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	diskClient := meta.(*ArmClient).diskClient

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := diskClient.Get(resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %s (resource group %s): %s", name, resGroup, err)
	}

	d.SetId(*resp.ID)
	if resp.Properties != nil {
		flattenAzureRmManagedDiskProperties(d, resp.Properties)
	}

	if resp.CreationData != nil {
		flattenAzureRmManagedDiskCreationData(d, resp.CreationData)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
