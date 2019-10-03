package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/entities"
)

func resourceArmStorageTableEntity() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageTableEntityCreateUpdate,
		Read:   resourceArmStorageTableEntityRead,
		Update: resourceArmStorageTableEntityCreateUpdate,
		Delete: resourceArmStorageTableEntityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageTableName,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},
			"partition_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"row_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
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

func resourceArmStorageTableEntityCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).Storage

	accountName := d.Get("storage_account_name").(string)
	tableName := d.Get("table_name").(string)
	partitionKey := d.Get("partition_key").(string)
	rowKey := d.Get("row_key").(string)
	entity := d.Get("entity").(map[string]interface{})

	resourceGroup, err := storageClient.FindResourceGroup(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Stable %q (Account %s): %s", tableName, accountName, err)
	}
	if resourceGroup == nil {
		if d.IsNewResource() {
			return fmt.Errorf("Unable to locate Resource Group for Storage Table %q (Account %s)", tableName, accountName)
		} else {
			log.Printf("[DEBUG] Unable to locate Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", tableName, accountName)
			d.SetId("")
			return nil
		}
	}

	client, err := storageClient.TableEntityClient(ctx, *resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error building Entity Client: %s", err)
	}

	if features.ShouldResourcesBeImported() {
		input := entities.GetEntityInput{
			PartitionKey:  partitionKey,
			RowKey:        rowKey,
			MetaDataLevel: entities.NoMetaData,
		}
		existing, err := client.Get(ctx, accountName, tableName, input)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", partitionKey, rowKey, tableName, accountName, *resourceGroup, err)
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
		return fmt.Errorf("Error creating Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %+v", partitionKey, rowKey, tableName, accountName, *resourceGroup, err)
	}

	resourceID := client.GetResourceID(accountName, tableName, partitionKey, rowKey)
	d.SetId(resourceID)

	return resourceArmStorageTableEntityRead(d, meta)
}

func resourceArmStorageTableEntityRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", id.TableName, id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TableEntityClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Table Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, *resourceGroup, err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.NoMetaData,
	}

	result, err := client.Get(ctx, id.AccountName, id.TableName, input)
	if err != nil {
		return fmt.Errorf("Error retrieving Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, *resourceGroup, err)
	}

	d.Set("storage_account_name", id.AccountName)
	d.Set("table_name", id.TableName)
	d.Set("partition_key", id.PartitionKey)
	d.Set("row_key", id.RowKey)
	if err := d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("Error setting `entity` for Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, *resourceGroup, err)
	}

	return nil
}

func resourceArmStorageTableEntityDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	storageClient := meta.(*ArmClient).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup, err := storageClient.FindResourceGroup(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error locating Resource Group for Storage Table %q (Account %s): %s", id.TableName, id.AccountName, err)
	}
	if resourceGroup == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed already", id.TableName, id.AccountName)
		return nil
	}

	client, err := storageClient.TableEntityClient(ctx, *resourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error building Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, *resourceGroup, err)
	}

	input := entities.DeleteEntityInput{
		PartitionKey: id.PartitionKey,
		RowKey:       id.RowKey,
	}

	if _, err := client.Delete(ctx, id.AccountName, id.TableName, input); err != nil {
		return fmt.Errorf("Error deleting Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, *resourceGroup, err)
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
