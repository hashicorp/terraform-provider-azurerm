// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_queue

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageQueueRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"metadata": helpers.MetaDataComputedSchema(),

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageQueueRead(d *pluginsdk.ResourceData, meta any) error {
	queueClient := client.From(meta).ResourceManager.StorageQueues
	ctx, cancel := timeouts.ForRead(clientsAccount.StopContext(meta), d)
	defer cancel()

	queueName := d.Get("name").(string)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := storagequeues.NewQueueID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, queueName)

	resp, err := queueClient.QueueGet(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %v", id, err)
	}

	d.Set("name", id.QueueName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if metadata := props.Metadata; metadata != nil {
				if err := d.Set("metadata", helpers.FlattenMetaData(*metadata)); err != nil {
					return fmt.Errorf("setting `metadata`: %s", err)
				}
			}
		}
	}

	account, err := client.From(meta).GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Queue %q: %v", id, err)
	}
	// Determine the queue endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeQueue)
	if err != nil {
		return fmt.Errorf("determining Queue endpoint: %v", err)
	}
	// Parse the queue endpoint as a data plane account ID
	accountDpId, err := accounts.ParseAccountID(*endpoint, client.From(meta).StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}
	d.Set("url", queues.NewQueueID(*accountDpId, id.QueueName).ID())

	d.SetId(id.ID())

	return nil
}
