package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
)

func dataSourceStorageTableEntity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStorageTableEntityRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageTableName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"partition_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"row_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"entity": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceStorageTableEntityRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	storageAccountName := d.Get("storage_account_name").(string)
	tableName := d.Get("table_name").(string)
	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)

	account, err := storageClient.FindAccount(ctx, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", storageAccountName, tableName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", tableName, storageAccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Table Entity Client for Storage Account %q (Resource Group %q): %s", storageAccountName, account.ResourceGroup, err)
	}

	id := client.GetResourceID(storageAccountName, tableName, partitionKey, rowKey)

	input := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := client.Get(ctx, storageAccountName, tableName, input)
	if err != nil {
		return fmt.Errorf("Error retrieving Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", partitionKey, rowKey, tableName, storageAccountName, account.ResourceGroup, err)
	}

	d.Set("storage_account_name", storageAccountName)
	d.Set("table_name", tableName)
	d.Set("partition_key", partitionKey)
	d.Set("row_key", rowKey)
	if err := d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("Error setting `entity` for Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", partitionKey, rowKey, tableName, storageAccountName, account.ResourceGroup, err)
	}
	d.SetId(id)

	return nil
}
