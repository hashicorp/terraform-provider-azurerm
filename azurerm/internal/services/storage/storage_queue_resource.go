package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceStorageQueue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageQueueCreate,
		Read:   resourceStorageQueueRead,
		Update: resourceStorageQueueUpdate,
		Delete: resourceStorageQueueDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer:      pluginsdk.DefaultImporter(),
		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.QueueV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageQueueName,
			},

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"metadata": MetaDataSchema(),
		},
	}
}

func resourceStorageQueueCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %s", accountName, queueName, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q", accountName)
	}

	client, err := storageClient.QueuesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	resourceId := parse.NewStorageQueueDataPlaneId(accountName, storageClient.Environment.StorageEndpointSuffix, queueName).ID()

	exists, err := client.Exists(ctx, account.ResourceGroup, accountName, queueName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Queue %q (Storage Account %q): %s", queueName, accountName, err)
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_queue", resourceId)
	}

	if err := client.Create(ctx, account.ResourceGroup, accountName, queueName, metaData); err != nil {
		return fmt.Errorf("creating Queue %q (Account %q): %+v", queueName, accountName, err)
	}

	d.SetId(resourceId)
	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageQueueDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.QueuesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	if err := client.UpdateMetaData(ctx, account.ResourceGroup, id.AccountName, id.Name, metaData); err != nil {
		return fmt.Errorf("updating MetaData for Queue %q (Storage Account %q): %s", id.Name, id.AccountName, err)
	}

	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageQueueDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.Name, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.QueuesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	queue, err := client.Get(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Queue %q (Account %q): %+v", id.AccountName, id.Name, err)
	}
	if queue == nil {
		log.Printf("[INFO] Storage Queue %q no longer exists, removing from state...", id.Name)
		d.SetId("")
		return nil
	}

	d.Set("name", id.Name)
	d.Set("storage_account_name", id.AccountName)

	if err := d.Set("metadata", FlattenMetaData(queue.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %s", err)
	}

	return nil
}

func resourceStorageQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageQueueDataPlaneID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %s", id.AccountName, id.Name, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.Name, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.QueuesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	if err := client.Delete(ctx, account.ResourceGroup, id.AccountName, id.Name); err != nil {
		return fmt.Errorf("deleting Storage Queue %q (Account %q): %s", id.Name, id.AccountName, err)
	}

	return nil
}
