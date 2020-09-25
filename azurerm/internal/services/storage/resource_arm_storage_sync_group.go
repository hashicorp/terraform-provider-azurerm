package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageSyncGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageSyncGroupCreate,
		Read:   resourceArmStorageSyncGroupRead,
		Delete: resourceArmStorageSyncGroupDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parsers.StorageSyncGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncId,
			},
		},
	}
}

func resourceArmStorageSyncGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	ssId, _ := parsers.ParseStorageSyncID(d.Get("storage_sync_id").(string))

	existing, err := client.Get(ctx, ssId.ResourceGroup, ssId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_storage_sync_group", *existing.ID)
	}

	if _, err := client.Create(ctx, ssId.ResourceGroup, ssId.Name, name, storagesync.SyncGroupCreateParameters{}); err != nil {
		return fmt.Errorf("creating Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, ssId.ResourceGroup, ssId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", name, ssId.Name, ssId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q) ID is empty or nil", name, ssId.Name, ssId.ResourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmStorageSyncGroupRead(d, meta)
}

func resourceArmStorageSyncGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ssClient := meta.(*clients.Client).Storage.SyncServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.StorageSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Sync Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", id.Name, id.StorageSyncName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	ssResp, err := ssClient.Get(ctx, id.ResourceGroup, id.StorageSyncName)
	if err != nil {
		return fmt.Errorf("reading Storage Sync %q (Resource Group %q): %+v", id.StorageSyncName, id.ResourceGroup, err)
	}

	if ssResp.ID == nil || *ssResp.ID == "" {
		return fmt.Errorf("reading Storage Sync %q (Resource Group %q) ID is empty or nil", id.StorageSyncName, id.ResourceGroup)
	}

	d.Set("storage_sync_id", ssResp.ID)

	return nil
}

func resourceArmStorageSyncGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.StorageSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.StorageSyncName, id.Name); err != nil {
		return fmt.Errorf("deleting Storage Sync Group (Storage Sync Group Name %q / Storage Sync Name %q /Resource Group %q): %+v", id.Name, id.StorageSyncName, id.ResourceGroup, err)
	}

	return nil
}
