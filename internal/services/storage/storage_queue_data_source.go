// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/queueservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func dataSourceStorageQueue() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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

			"metadata": MetaDataComputedSchema(),

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
	}

	return r
}

func dataSourceStorageQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	queueClient := meta.(*clients.Client).Storage.ResourceManager.QueueService
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueName := d.Get("name").(string)

	if !features.FivePointOh() {
		if accountName := d.Get("storage_account_name").(string); accountName != "" {
			storageClient := meta.(*clients.Client).Storage
			account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
			if err != nil {
				return fmt.Errorf("retrieving Storage Account %q for Queue %q: %v", accountName, queueName, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q for Queue %q", accountName, queueName)
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
			accountId, err := accounts.ParseAccountID(pointer.From(endpoint), storageClient.StorageDomainSuffix)
			if err != nil {
				return err
			}

			id := queues.NewQueueID(*accountId, queueName)

			props, err := queuesDataPlaneClient.Get(ctx, queueName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			d.SetId(id.ID())

			d.Set("name", queueName)
			d.Set("storage_account_name", accountName)

			if props != nil {
				if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
					return fmt.Errorf("setting `metadata`: %v", err)
				}
			}

			resourceManagerId := parse.NewStorageQueueResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", queueName)
			d.Set("resource_manager_id", resourceManagerId.ID())
			d.Set("url", id.ID())

			return nil
		}
	}

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := queueservice.NewQueueID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, queueName)

	resp, err := queueClient.QueueGet(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %v", id, err)
	}

	d.Set("name", id.QueueName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if metadata := props.Metadata; metadata != nil {
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
	}

	d.SetId(id.ID())

	return nil
}
