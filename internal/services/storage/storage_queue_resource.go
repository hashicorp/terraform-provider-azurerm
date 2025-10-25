// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/queueservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
	r := &pluginsdk.Resource{
		Create: resourceStorageQueueCreate,
		Read:   resourceStorageQueueRead,
		Update: resourceStorageQueueUpdate,
		Delete: resourceStorageQueueDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			if !features.FivePointOh() {
				if strings.HasPrefix(id, "/subscriptions/") {
					_, err := queueservice.ParseQueueID(id)
					return err
				}
				_, err := queues.ParseQueueID(id, storageDomainSuffix)
				return err
			}

			_, err := queueservice.ParseQueueID(id)
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

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"metadata": MetaDataSchema(),

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.StorageAccountName,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
			Deprecated:   "the `storage_account_name` property has been deprecated in favour of `storage_account_id` and will be removed in version 5.0 of the Provider.",
		}

		r.Schema["storage_account_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
		}

		r.Schema["resource_manager_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Computed:   true,
			Deprecated: "the `resource_manager_id` property has been deprecated in favour of `id` and will be removed in version 5.0 of the Provider.",
		}

		r.CustomizeDiff = func(ctx context.Context, diff *pluginsdk.ResourceDiff, i interface{}) error {
			// Resource Manager ID in use, but change to `storage_account_id` should recreate
			if strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_id") {
				return diff.ForceNew("storage_account_id")
			}

			// using legacy Data Plane ID but attempting to change the storage_account_name should recreate
			if diff.Id() != "" && !strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_name") {
				// converting from storage_account_id to the deprecated storage_account_name is not supported
				oldAccountId, _ := diff.GetChange("storage_account_id")
				oldName, newName := diff.GetChange("storage_account_name")

				if oldAccountId.(string) != "" && newName.(string) != "" {
					return diff.ForceNew("storage_account_name")
				}

				if oldName.(string) != "" && newName.(string) != "" {
					return diff.ForceNew("storage_account_name")
				}
			}

			return nil
		}
	}

	return r
}

func resourceStorageQueueCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	queueClient := meta.(*clients.Client).Storage.ResourceManager.QueueService
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueName := d.Get("name").(string)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := ExpandMetaData(metaDataRaw)

	if !features.FivePointOh() {
		if accountName := d.Get("storage_account_name").(string); accountName != "" {
			storageClient := meta.(*clients.Client).Storage

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
	}

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := queueservice.NewQueueID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, queueName)

	existing, err := queueClient.QueueGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %q: %v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_queue", id.ID())
	}

	payload := queueservice.StorageQueue{
		Properties: &queueservice.QueueProperties{
			Metadata: &metaData,
		},
	}

	if _, err := queueClient.QueueCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	queueClient := meta.(*clients.Client).Storage.ResourceManager.QueueService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage

		id, err := queues.ParseQueueID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

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

		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := ExpandMetaData(metaDataRaw)

		if err = client.UpdateMetaData(ctx, id.QueueName, metaData); err != nil {
			return fmt.Errorf("updating MetaData for %s: %v", id, err)
		}

		return resourceStorageQueueRead(d, meta)
	}

	id, err := queueservice.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	existing, err := queueClient.QueueGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("unexpected null model after retrieving %v", id)
	}

	payload := queueservice.StorageQueue{
		Properties: existing.Model.Properties,
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		payload.Properties.Metadata = pointer.To(ExpandMetaData(metaDataRaw))
	}

	if _, err := queueClient.QueueCreate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageQueueRead(d, meta)
}

func resourceStorageQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	queueClient := meta.(*clients.Client).Storage.ResourceManager.QueueService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() {
		if !strings.HasPrefix(d.Id(), "/subscriptions/") {
			if said := d.Get("storage_account_id").(string); said == "" {
				storageClient := meta.(*clients.Client).Storage

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
				d.Set("url", id.ID())

				return nil
			} else {
				// Deal with the ID changing if the user changes from `storage_account_name` to `storage_account_id`
				accountId, err := commonids.ParseStorageAccountID(said)
				if err != nil {
					return err
				}

				id := queueservice.NewQueueID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))
				d.SetId(id.ID())
				// Continue the code flow outside this block
			}
		}
	}

	id, err := queueservice.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	existing, err := queueClient.QueueGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	d.Set("name", id.QueueName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())

	if model := existing.Model; model != nil {
		if prop := model.Properties; prop != nil {
			if metadata := prop.Metadata; metadata != nil {
				if err := d.Set("metadata", FlattenMetaData(*metadata)); err != nil {
					return fmt.Errorf("setting `metadata`: %s", err)
				}
			}
		}
	}

	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Queue %q: %v", id, err)
	}
	// Determine the queue endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeQueue)
	if err != nil {
		return fmt.Errorf("determining Queue endpoint: %v", err)
	}
	// Parse the queue endpoint as a data plane account ID
	accountDpId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}
	d.Set("url", queues.NewQueueID(*accountDpId, id.QueueName).ID())

	if !features.FivePointOh() {
		d.Set("resource_manager_id", id.ID())
		d.Set("storage_account_name", "")
	}

	return nil
}

func resourceStorageQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	queueClient := meta.(*clients.Client).Storage.ResourceManager.QueueService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage
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

	id, err := queueservice.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := queueClient.QueueDelete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}
