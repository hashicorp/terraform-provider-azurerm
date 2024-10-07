// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/table/tables"
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

	if !features.FourPointOhBeta() {
		resource.Schema["storage_table_id"].Required = false
		resource.Schema["storage_table_id"].Optional = true
		resource.Schema["storage_table_id"].Computed = true
		resource.Schema["storage_table_id"].ConflictsWith = []string{"table_name", "storage_account_name"}

		resource.Schema["table_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			Deprecated:    "the `table_name` and `storage_account_name` properties have been superseded by the `storage_table_id` property and will be removed in version 4.0 of the AzureRM provider",
			ConflictsWith: []string{"storage_table_id"},
			RequiredWith:  []string{"storage_account_name"},
			ValidateFunc:  validate.StorageTableName,
		}

		resource.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			Deprecated:    "the `table_name` and `storage_account_name` properties have been superseded by the `storage_table_id` property and will be removed in version 4.0 of the AzureRM provider",
			ConflictsWith: []string{"storage_table_id"},
			RequiredWith:  []string{"table_name"},
			ValidateFunc:  validate.StorageAccountName,
		}
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
	} else if !features.FourPointOhBeta() {
		// TODO: this is needed until `table_name` / `storage_account_name` are removed in favor of `storage_table_id` in v4.0
		// we will retrieve the storage account twice but this will make it easier to refactor later
		storageAccountName := d.Get("storage_account_name").(string)

		account, err := storageClient.FindAccount(ctx, subscriptionId, storageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q: %v", storageAccountName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", storageAccountName)
		}

		// Determine the table endpoint, so we can build a data plane ID
		endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
		if err != nil {
			return fmt.Errorf("determining Table endpoint: %v", err)
		}

		// Parse the table endpoint as a data plane account ID
		accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
		if err != nil {
			return fmt.Errorf("parsing Account ID: %v", err)
		}

		storageTableId = pointer.To(tables.NewTableID(*accountId, d.Get("table_name").(string)))
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

	if !features.FourPointOhBeta() {
		d.Set("storage_account_name", id.AccountId.AccountName)
		d.Set("table_name", id.TableName)
	}

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
