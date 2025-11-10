// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func dataSourceStorageTableEntity() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceStorageTableEntityRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_table_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageTableDataPlaneID,
			},

			"partition_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"row_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"entity": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}

	return resource
}

func dataSourceStorageTableEntityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)

	var storageTableId *tables.TableId
	var err error
	if v, ok := d.GetOk("storage_table_id"); ok && v.(string) != "" {
		storageTableId, err = tables.ParseTableID(v.(string), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}
	}

	if storageTableId == nil {
		return fmt.Errorf("determining storage table ID")
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageTableId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %v", storageTableId.AccountId.AccountName, storageTableId.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("the parent Storage Account %s was not found", storageTableId.AccountId.AccountName)
	}

	dataPlaneClient, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Table Entity Client for %s: %+v", account.StorageAccountId, err)
	}

	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
	if err != nil {
		return fmt.Errorf("retrieving the table data plane endpoint: %v", err)
	}

	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := entities.NewEntityID(*accountId, storageTableId.TableName, partitionKey, rowKey)

	input := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := dataPlaneClient.Get(ctx, storageTableId.TableName, input)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("storage_table_id", storageTableId.ID())
	d.Set("partition_key", partitionKey)
	d.Set("row_key", rowKey)

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
