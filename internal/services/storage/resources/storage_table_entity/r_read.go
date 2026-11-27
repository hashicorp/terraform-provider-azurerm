// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table_entity

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
)

func resourceStorageTableEntityRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := entities.ParseEntityID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	var tableName string
	var accountName string
	var storageTableIdFmtd string
	var account *client.AccountDetails

	tableIdRaw, ok := d.GetOk("storage_table_id")
	storageTableIdRaw := ""
	if ok {
		storageTableIdRaw = tableIdRaw.(string)
	}

	if storageTableIdRaw == "" {
		accountName = id.AccountId.AccountName
		tableName = id.TableName
		account, err = storageClient.FindAccount(ctx, subscriptionId, accountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Table %q: %v", accountName, tableName, err)
		}
		if account != nil {
			storageTableId := tables.NewTableID(subscriptionId, account.StorageAccountId.ResourceGroupName, accountName, tableName)
			storageTableIdFmtd = storageTableId.ID()
		}
	} else {
		storageTableId, err := tables.ParseTableID(storageTableIdRaw)
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
		log.Printf("[WARN] Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", tableName, accountName)
		d.SetId("")
		return nil
	}

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
		if response.WasNotFound(result.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("storage_table_id", storageTableIdFmtd)
	d.Set("partition_key", id.PartitionKey)
	d.Set("row_key", id.RowKey)

	if err = d.Set("entity", flattenEntity(result.Entity)); err != nil {
		return fmt.Errorf("setting `entity` for %s: %v", id, err)
	}

	return nil
}

// The api returns extra information that we already have. We'll remove it here before setting it in state.
func flattenEntity(entity map[string]any) map[string]any {
	delete(entity, "PartitionKey")
	delete(entity, "RowKey")
	delete(entity, "Timestamp")

	result := map[string]any{}
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
			// https://learn.microsoft.com/rest/api/storageservices/payload-format-for-table-service-operations#property-types-in-a-json-feed
			switch c := v.(type) {
			case bool:
				result[k] = fmt.Sprint(v)
				result[k+"@odata.type"] = "Edm.Boolean"
			case float64:
				f64 := v.(float64)
				if v == float64(int64(f64)) {
					result[k] = strconv.FormatInt(int64(f64), 10)
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
