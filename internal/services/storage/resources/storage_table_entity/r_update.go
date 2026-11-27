// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table_entity

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
)

func resourceStorageTableEntityUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := storageClient2.From(meta)

	ctx, cancel := timeouts.ForUpdate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	var tableName string
	var accountName string
	var account *storageClient2.AccountDetails

	tableIdRaw, ok := d.GetOk("storage_table_id")
	if !ok || tableIdRaw.(string) == "" {
		return errors.New("`storage_table_id` is required")
	}
	storageTableIdRaw := tableIdRaw.(string)

	storageTableId, err := tables.ParseTableID(storageTableIdRaw)
	if err != nil {
		return err
	}
	tableName = storageTableId.TableName
	accountName = storageTableId.StorageAccountName
	storageAccountId := commonids.NewStorageAccountID(storageTableId.SubscriptionId, storageTableId.ResourceGroupName, storageTableId.StorageAccountName)
	account, err = storageClient.GetAccount(ctx, storageAccountId)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
	}

	if account == nil {
		log.Printf("[DEBUG] Unable to locate Storage Account %q for Table %q - assuming removed & removing from state", accountName, tableName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Entity Client: %v", err)
	}

	input := entities.InsertOrMergeEntityInput{
		PartitionKey: d.Get("partition_key").(string),
		RowKey:       d.Get("row_key").(string),
		Entity:       d.Get("entity").(map[string]any),
	}

	if _, err = client.InsertOrMerge(ctx, tableName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageTableEntityRead(d, meta)
}
