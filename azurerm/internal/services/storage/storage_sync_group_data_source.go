package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageSyncId, err := parse.StorageSyncServiceID(d.Get("storage_sync_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, storageSyncId.ResourceGroup, storageSyncId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sync Group %q does not exist within Storage Sync %q / Resource Group %q", name, storageSyncId.Name, storageSyncId.ResourceGroup)
		}
		return fmt.Errorf("retrieving Sync Group %q (Storage Sync %q / Resource Group %q): %+v", name, storageSyncId.Name, storageSyncId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("ID is nil for Sync Group %q (Storage Sync %q / Resource Group %q)", name, storageSyncId.Name, storageSyncId.ResourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("storage_sync_id", storageSyncId.ID(subscriptionId))
	return nil
}
