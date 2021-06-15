package compute

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceManagedDisk() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagedDiskRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"create_option": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_iops_read_write": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"image_reference_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),

			"zones": azure.SchemaZonesComputed(),
		},
	}
}

func dataSourceManagedDiskRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if sku := resp.Sku; sku != nil {
		d.Set("storage_account_type", string(sku.Name))
	}

	if props := resp.DiskProperties; props != nil {
		if creationData := props.CreationData; creationData != nil {
			d.Set("create_option", string(creationData.CreateOption))

			imageReferenceID := ""
			if creationData.ImageReference != nil && creationData.ImageReference.ID != nil {
				imageReferenceID = *creationData.ImageReference.ID
			}
			d.Set("image_reference_id", imageReferenceID)

			d.Set("source_resource_id", creationData.SourceResourceID)
			d.Set("source_uri", creationData.SourceURI)
			d.Set("storage_account_id", creationData.StorageAccountID)
		}

		d.Set("disk_size_gb", props.DiskSizeGB)
		d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
		d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
		d.Set("os_type", props.OsType)

		diskEncryptionSetId := ""
		if props.Encryption != nil && props.Encryption.DiskEncryptionSetID != nil {
			diskEncryptionSetId = *props.Encryption.DiskEncryptionSetID
		}
		d.Set("disk_encryption_set_id", diskEncryptionSetId)
	}

	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	return tags.FlattenAndSet(d, resp.Tags)
}
