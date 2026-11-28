// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table_entity

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
)

func resourceStorageTableEntityCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)

	var tableName string
	var accountName string
	var account *client.AccountDetails
	var err error

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
		return fmt.Errorf("the parent Storage Account %s was not found", accountName)
	}

	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
	if err != nil {
		return fmt.Errorf("retrieving the table data plane endpoint: %v", err)
	}

	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	id := entities.NewEntityID(*accountId, tableName, partitionKey, rowKey)

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Entity Client: %v", err)
	}

	getEntityInput := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, tableName, getEntityInput)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) && !response.WasForbidden(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_table_entity", id.ID())
		}
	}

	input := entities.InsertOrMergeEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity:       d.Get("entity").(map[string]any),
	}

	if _, err = client.InsertOrMerge(ctx, tableName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageTableEntityRead(d, meta)
}
