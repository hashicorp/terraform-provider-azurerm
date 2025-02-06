// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func resourceStorageQueue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageQueueCreate,
		Read:   resourceStorageQueueRead,
		Update: resourceStorageQueueUpdate,
		Delete: resourceStorageQueueDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := queues.ParseQueueID(id, storageDomainSuffix)
			return err
		}),

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

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageQueueCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %v", accountName, queueName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	queuesDataPlaneClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %v", err)
	}

	// Determine the queue endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeQueue)
	if err != nil {
		return fmt.Errorf("determining Queue endpoint: %v", err)
	}

	// Parse the queue endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := queues.NewQueueID(*accountId, queueName).ID()

	exists, err := queuesDataPlaneClient.Exists(ctx, queueName)
	if err != nil {
		return fmt.Errorf("checking for existing %s: %v", id, err)
	}
	if exists != nil && *exists {
		return tf.ImportAsExistsError("azurerm_storage_queue", id)
	}

	if err = queuesDataPlaneClient.Create(ctx, queueName, metaData); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id)

	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %v", id.AccountId.AccountName, id.QueueName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %v", err)
	}

	if err = client.UpdateMetaData(ctx, id.QueueName, metaData); err != nil {
		return fmt.Errorf("updating MetaData for %s: %v", id, err)
	}

	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %v", id.AccountId.AccountName, id.QueueName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.QueueName, id.AccountId.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %v", err)
	}

	queue, err := client.Get(ctx, id.QueueName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}
	if queue == nil {
		log.Printf("[INFO] Storage Queue %q no longer exists, removing from state...", id.QueueName)
		d.SetId("")
		return nil
	}

	d.Set("name", id.QueueName)
	d.Set("storage_account_name", id.AccountId.AccountName)

	if err := d.Set("metadata", FlattenMetaData(queue.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %s", err)
	}

	resourceManagerId := parse.NewStorageQueueResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, id.AccountId.AccountName, "default", id.QueueName)
	d.Set("resource_manager_id", resourceManagerId.ID())

	return nil
}

func resourceStorageQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queues.ParseQueueID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Queue %q: %s", id.AccountId.AccountName, id.QueueName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Queue %q (Account %s) - assuming removed & removing from state", id.QueueName, id.AccountId.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %v", err)
	}

	if err = client.Delete(ctx, id.QueueName); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
