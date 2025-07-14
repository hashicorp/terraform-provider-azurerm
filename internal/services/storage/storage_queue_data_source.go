// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func dataSourceStorageQueue() *pluginsdk.Resource {
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

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"metadata": MetaDataComputedSchema(),

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

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
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
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

	return nil
}
