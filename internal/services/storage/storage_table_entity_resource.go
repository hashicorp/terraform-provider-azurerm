// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func resourceStorageTableEntity() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageTableEntityCreate,
		Read:   resourceStorageTableEntityRead,
		Update: resourceStorageTableEntityUpdate,
		Delete: resourceStorageTableEntityDelete,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := entities.ParseEntityID(id, storageDomainSuffix)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
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

	return resource
}

func resourceStorageTableEntityCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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
	}

	if storageTableId == nil {
		return fmt.Errorf("determining storage table ID")
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageTableId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %v", storageTableId.AccountId.AccountName, storageTableId.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q for Table %q", storageTableId.AccountId.AccountName, storageTableId.TableName)
	}

	accountId, err := accounts.ParseAccountID(storageTableId.ID(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %s", err)
	}

	id := entities.NewEntityID(*accountId, storageTableId.TableName, partitionKey, rowKey)

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Entity Client: %v", err)
	}

	getEntityInput := entities.GetEntityInput{
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		MetaDataLevel: entities.NoMetaData,
	}
	existing, err := client.Get(ctx, storageTableId.TableName, getEntityInput)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) && !response.WasForbidden(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_table_entity", id.ID())
	}

	input := entities.InsertOrMergeEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity:       d.Get("entity").(map[string]interface{}),
	}

	if _, err = client.InsertOrMerge(ctx, storageTableId.TableName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageTableEntityRead(d, meta)
}

func resourceStorageTableEntityUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	storageTableId := tables.NewTableID(id.AccountId, id.TableName)

	account, err := storageClient.FindAccount(ctx, subscriptionId, storageTableId.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %v", storageTableId.AccountId.AccountName, storageTableId.TableName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Storage Account %q for Table %q - assuming removed & removing from state", storageTableId.AccountId.AccountName, storageTableId.TableName)
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
		Entity:       d.Get("entity").(map[string]interface{}),
	}

	if _, err = client.InsertOrMerge(ctx, storageTableId.TableName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageTableEntityRead(d, meta)
}

func resourceStorageTableEntityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Table %q: %s", id.AccountId.AccountName, id.TableName, err)
	}
	if account == nil {
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountId.AccountName)
		d.SetId("")
		return nil
	}

	storageTableId := tables.NewTableID(id.AccountId, id.TableName)

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Table Entity Client for %s: %+v", account.StorageAccountId, err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.FullMetaData,
	}

	result, err := client.Get(ctx, id.TableName, input)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("storage_table_id", storageTableId.ID())
	d.Set("partition_key", id.PartitionKey)
	d.Set("row_key", id.RowKey)

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	return nil
}

func resourceStorageTableEntityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %s", id.AccountId.AccountName, id.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Entity Client for %s: %+v", account.StorageAccountId, err)
	}

	input := entities.DeleteEntityInput{
		PartitionKey: id.PartitionKey,
		RowKey:       id.RowKey,
	}

	if _, err = client.Delete(ctx, id.TableName, input); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
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
			case "Edm.Int32", "Edm.Int64":
				// `v` returned as string for int 64
				result[k] = fmt.Sprint(v)
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
					// fmt.Sprintf("%f", v) will return `123.123000` for `123.123`, have to use fmt.Sprint
					result[k] = fmt.Sprint(v)
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
