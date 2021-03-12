package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
)

func resourceStorageTableEntity() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageTableEntityCreateUpdate,
		Read:   resourceStorageTableEntityRead,
		Update: resourceStorageTableEntityCreateUpdate,
		Delete: resourceStorageTableEntityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageTableName,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateStorageAccountName,
			},
			"partition_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"row_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"entity": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStorageTableEntityCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	accountName := d.Get("storage_account_name").(string)
	tableName := d.Get("table_name").(string)
	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)
	entity := d.Get("entity").(map[string]interface{})

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", accountName, tableName, err)
	}
	if account == nil {
		if d.IsNewResource() {
			return fmt.Errorf("Unable to locate Account %q for Storage Table %q", accountName, tableName)
		} else {
			log.Printf("[DEBUG] Unable to locate Account %q for Storage Table %q - assuming removed & removing from state", accountName, tableName)
			d.SetId("")
			return nil
		}
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Entity Client: %s", err)
	}

	if d.IsNewResource() {
		input := entities.GetEntityInput{
			PartitionKey:  partitionKey,
			RowKey:        rowKey,
			MetaDataLevel: entities.NoMetaData,
		}
		existing, err := client.Get(ctx, accountName, tableName, input)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", partitionKey, rowKey, tableName, accountName, account.ResourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			id := client.GetResourceID(accountName, tableName, partitionKey, rowKey)
			return tf.ImportAsExistsError("azurerm_storage_table_entity", id)
		}
	}

	input := entities.InsertOrMergeEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity:       entity,
	}

	if _, err := client.InsertOrMerge(ctx, accountName, tableName, input); err != nil {
		return fmt.Errorf("Error creating Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %+v", partitionKey, rowKey, tableName, accountName, account.ResourceGroup, err)
	}

	resourceID := client.GetResourceID(accountName, tableName, partitionKey, rowKey)
	d.SetId(resourceID)

	return resourceStorageTableEntityRead(d, meta)
}

func resourceStorageTableEntityRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", id.AccountName, id.TableName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Table Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := client.Get(ctx, id.AccountName, id.TableName, input)
	if err != nil {
		return fmt.Errorf("Error retrieving Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("storage_account_name", id.AccountName)
	d.Set("table_name", id.TableName)
	d.Set("partition_key", id.PartitionKey)
	d.Set("row_key", id.RowKey)
	if err := d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("Error setting `entity` for Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func resourceStorageTableEntityDelete(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Table %q: %s", id.AccountName, id.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("Storage Account %q was not found!", id.AccountName)
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	input := entities.DeleteEntityInput{
		PartitionKey: id.PartitionKey,
		RowKey:       id.RowKey,
	}

	if _, err := client.Delete(ctx, id.AccountName, id.TableName, input); err != nil {
		return fmt.Errorf("Error deleting Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

// The api returns extra information that we already have. We'll remove it here before setting it in state.
func flattenEntity(entity map[string]interface{}) map[string]interface{} {
	delete(entity, "PartitionKey")
	delete(entity, "RowKey")
	delete(entity, "Timestamp")

	return entity
}
