package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStorageSyncGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageSyncGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageSyncId,
			},
		},
	}
}

func dataSourceArmStorageSyncGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageSyncId, err := parsers.ParseStorageSyncID(d.Get("storage_sync_id").(string))

	resp, err := client.Get(ctx, ssId.ResourceGroup, ssId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO]  does not exist - removing from state", name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Sync Group %q (Storage Sync Name %q / Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Storage Sync Group %q (Storage Sync Name %q /Resource Group %q) ID is empty or nil", name, ssId.Name, ssId.ResourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("storage_sync_id", storageSyncId)
	return nil
}
