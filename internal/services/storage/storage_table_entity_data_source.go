// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	rmTables "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
				ValidateFunc: rmTables.ValidateTableID,
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

	if !features.FivePointOh() {
		resource.Schema["storage_table_id"].ValidateFunc = validation.Any(rmTables.ValidateTableID, storageValidate.StorageTableDataPlaneID)
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

	var tableName string
	var accountName string
	var account *client.AccountDetails
	var err error

	tableIdRaw, ok := d.GetOk("storage_table_id")
	if !ok || tableIdRaw.(string) == "" {
		return fmt.Errorf("`storage_table_id` is required")
	}
	storageTableIdRaw := tableIdRaw.(string)
	storageTableIdFmtd := ""

	if !features.FivePointOh() {
		if strings.HasPrefix(strings.ToLower(storageTableIdRaw), "/subscriptions/") {
			storageTableId, err := rmTables.ParseTableID(storageTableIdRaw)
			if err != nil {
				return err
			}
			storageTableIdFmtd = storageTableId.ID()
			tableName = storageTableId.TableName
			accountName = storageTableId.StorageAccountName
			storageAccountId := commonids.NewStorageAccountID(storageTableId.SubscriptionId, storageTableId.ResourceGroupName, storageTableId.StorageAccountName)
			account, err = storageClient.GetAccount(ctx, storageAccountId)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
			}
		} else {
			log.Printf("[WARN] `storage_table_id` is currently configured as a Data Plane URL. This legacy behavior has been deprecated and will be removed in version 5.0 of the AzureRM Provider. Please migrate to the Management Plane ID format.")
			storageTableId, err := tables.ParseTableID(storageTableIdRaw, storageClient.StorageDomainSuffix)
			if err != nil {
				return err
			}
			storageTableIdFmtd = storageTableId.ID()
			tableName = storageTableId.TableName
			accountName = storageTableId.AccountId.AccountName
			account, err = storageClient.FindAccount(ctx, subscriptionId, accountName)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
			}
		}
	} else {
		storageTableId, err := rmTables.ParseTableID(storageTableIdRaw)
		if err != nil {
			return err
		}
		storageTableIdFmtd = storageTableId.ID()
		tableName = storageTableId.TableName
		accountName = storageTableId.StorageAccountName
		storageAccountId := commonids.NewStorageAccountID(storageTableId.SubscriptionId, storageTableId.ResourceGroupName, storageTableId.StorageAccountName)
		account, err = storageClient.GetAccount(ctx, storageAccountId)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
		}
	}

	if account == nil {
		return fmt.Errorf("the parent Storage Account %s was not found", accountName)
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

	id := entities.NewEntityID(*accountId, tableName, partitionKey, rowKey)

	input := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := dataPlaneClient.Get(ctx, tableName, input)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("storage_table_id", storageTableIdFmtd)
	d.Set("partition_key", partitionKey)
	d.Set("row_key", rowKey)

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
