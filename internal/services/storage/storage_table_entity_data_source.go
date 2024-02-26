// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/table/entities"
)

func dataSourceStorageTableEntity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageTableEntityRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"table_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageTableName,
			},

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageAccountName,
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
}

func dataSourceStorageTableEntityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	tableName := d.Get("table_name").(string)
	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
	}
	if account == nil {
		return fmt.Errorf("the parent Storage Account %s was not found", accountName)
	}

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Table Entity Client for Storage Account %q (Resource Group %q): %v", accountName, account.ResourceGroup, err)
	}

	accountId, err := accounts.ParseAccountID(accountName, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := entities.NewEntityID(*accountId, tableName, partitionKey, rowKey)

	input := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := client.Get(ctx, tableName, input)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("storage_account_name", accountName)
	d.Set("table_name", tableName)
	d.Set("partition_key", partitionKey)
	d.Set("row_key", rowKey)

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
