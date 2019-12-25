package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmManagedDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagedDiskRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"zones": azure.SchemaZonesComputed(),

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

			"create_option": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"disk_iops_read_write": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"encryption_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_disk_encryption_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmManagedDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Managed Disk %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("[ERROR] Error making Read request on Azure Managed Disk %q (Resource Group %q): %s", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			flattenAzureRmManagedDiskCreationData(d, creationData)
		}
		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
		d.Set("os_type", props.OsType)
		if encryption := props.Encryption; encryption != nil {
			d.Set("encryption_type", string(encryption.Type))
			d.Set("managed_disk_encryption_set_id", encryption.DiskEncryptionSetID)
		}
	}

	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	return tags.FlattenAndSet(d, resp.Tags)
}
