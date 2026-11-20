// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table_entity

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
)

func resourceStorageTableEntityDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	var tableName string
	var accountName string
	var account *client.AccountDetails

	tableIdRaw, ok := d.GetOk("storage_table_id")
	if !ok || tableIdRaw.(string) == "" {
		return fmt.Errorf("`storage_table_id` is required")
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
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Entity Client for %s: %+v", account.StorageAccountId, err)
	}

	input := entities.DeleteEntityInput{
		PartitionKey: id.PartitionKey,
		RowKey:       id.RowKey,
	}

	if _, err = client.Delete(ctx, tableName, input); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
