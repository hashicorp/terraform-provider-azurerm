// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/entities"
)

func resourceStorageTableEntity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageTableEntityCreateUpdate,
		Read:   resourceStorageTableEntityRead,
		Update: resourceStorageTableEntityCreateUpdate,
		Delete: resourceStorageTableEntityDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := entities.ParseResourceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"table_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageTableName,
			},
			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},
			"partition_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"row_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"entity": {
				Type:     pluginsdk.TypeMap,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceStorageTableEntityCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving Account %q for Table %q: %s", accountName, tableName, err)
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
		return fmt.Errorf("building Entity Client: %s", err)
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
				return fmt.Errorf("checking for presence of existing Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", partitionKey, rowKey, tableName, accountName, account.ResourceGroup, err)
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
		return fmt.Errorf("creating Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %+v", partitionKey, rowKey, tableName, accountName, account.ResourceGroup, err)
	}

	resourceID := client.GetResourceID(accountName, tableName, partitionKey, rowKey)
	d.SetId(resourceID)

	return resourceStorageTableEntityRead(d, meta)
}

func resourceStorageTableEntityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", id.AccountName, id.TableName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Table Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.FullMetaData,
	}

	result, err := client.Get(ctx, id.AccountName, id.TableName, input)
	if err != nil {
		return fmt.Errorf("retrieving Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("storage_account_name", id.AccountName)
	d.Set("table_name", id.TableName)
	d.Set("partition_key", id.PartitionKey)
	d.Set("row_key", id.RowKey)
	if err := d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func resourceStorageTableEntityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	id, err := entities.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", id.AccountName, id.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("Storage Account %q was not found!", id.AccountName)
	}

	client, err := storageClient.TableEntityClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("building Entity Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	input := entities.DeleteEntityInput{
		PartitionKey: id.PartitionKey,
		RowKey:       id.RowKey,
	}

	if _, err := client.Delete(ctx, id.AccountName, id.TableName, input); err != nil {
		return fmt.Errorf("deleting Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %s", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

// The api returns extra information that we already have. We'll remove it here before setting it in state.
func flattenEntity(entity map[string]interface{}) map[string]interface{} {
	delete(entity, "PartitionKey")
	delete(entity, "RowKey")
	delete(entity, "Timestamp")

	result := map[string]interface{}{}
	for k, v := range entity {
		// skip ODATA annotation returned with fullmetadata
		if strings.HasPrefix(k, "odata.") || strings.HasSuffix(k, "@odata.type") {
			continue
		}
		if dtype, ok := entity[k+"@odata.type"]; ok {
			switch dtype {
			case "Edm.Boolean":
				result[k] = fmt.Sprint(v)
			case "Edm.Double":
				result[k] = fmt.Sprintf("%f", v)
			case "Edm.Int32":
				fallthrough
			case "Edm.Int64":
				result[k] = fmt.Sprintf("%d", int64(v.(float64)))
			case "Edm.String":
				result[k] = v
			default:
				log.Printf("[WARN] key %q with unexpected @odata.type %q", k, dtype)
				continue
			}

			result[k+"@odata.type"] = dtype
		} else {
			// special handling for property types that do not require the annotation to be present
			// https://docs.microsoft.com/en-us/rest/api/storageservices/payload-format-for-table-service-operations#property-types-in-a-json-feed
			switch c := v.(type) {
			case bool:
				result[k] = fmt.Sprint(v)
				result[k+"@odata.type"] = "Edm.Boolean"
			case float64:
				f64 := v.(float64)
				if v == float64(int64(f64)) {
					result[k] = fmt.Sprintf("%d", int64(f64))
					result[k+"@odata.type"] = "Edm.Int32"
				} else {
					result[k] = fmt.Sprintf("%f", v)
					result[k+"@odata.type"] = "Edm.Double"
				}
			case string:
				result[k] = v
			default:
				log.Printf("[WARN] key %q with unexpected type %T", k, c)
			}
		}
	}

	return result
}
