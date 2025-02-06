// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type storageTableEntitiesDataSource struct{}

var _ sdk.DataSource = storageTableEntitiesDataSource{}

type TableEntitiesDataSourceModel struct {
	StorageTableId     string                       `tfschema:"storage_table_id"`
	TableName          string                       `tfschema:"table_name,removedInNextMajorVersion"`
	StorageAccountName string                       `tfschema:"storage_account_name,removedInNextMajorVersion"`
	Filter             string                       `tfschema:"filter"`
	Select             []string                     `tfschema:"select"`
	Items              []TableEntityDataSourceModel `tfschema:"items"`
}

type TableEntityDataSourceModel struct {
	PartitionKey string                 `tfschema:"partition_key"`
	RowKey       string                 `tfschema:"row_key"`
	Properties   map[string]interface{} `tfschema:"properties"`
}

func (k storageTableEntitiesDataSource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"storage_table_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageTableDataPlaneID,
		},

		"filter": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"select": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}

	return s
}

func (k storageTableEntitiesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"items": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"partition_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"row_key": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"properties": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (k storageTableEntitiesDataSource) ModelObject() interface{} {
	return &TableEntitiesDataSourceModel{}
}

func (k storageTableEntitiesDataSource) ResourceType() string {
	return "azurerm_storage_table_entities"
}

func (k storageTableEntitiesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model TableEntitiesDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			storageClient := metadata.Client.Storage
			subscriptionId := metadata.Client.Account.SubscriptionId

			var storageTableId *tables.TableId
			var err error
			if model.StorageTableId != "" {
				storageTableId, err = tables.ParseTableID(model.StorageTableId, storageClient.StorageDomainSuffix)
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
				return fmt.Errorf("the parent Storage Account %s was not found", storageTableId.AccountId.AccountName)
			}

			client, err := storageClient.TableEntityDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Table Entity Client for %s: %+v", account.StorageAccountId, err)
			}

			input := entities.QueryEntitiesInput{
				Filter:        &model.Filter,
				MetaDataLevel: entities.MinimalMetaData,
			}

			if model.Select != nil {
				model.Select = append(model.Select, "RowKey", "PartitionKey")
				input.PropertyNamesToSelect = &model.Select
			}

			id := parse.NewStorageTableEntitiesId(storageTableId.AccountId.AccountName, storageClient.StorageDomainSuffix, storageTableId.TableName, model.Filter)

			result, err := client.Query(ctx, storageTableId.TableName, input)
			if err != nil {
				return fmt.Errorf("retrieving Entities (Filter %q) (Table %q in %s): %+v", model.Filter, storageTableId.TableName, account.StorageAccountId, err)
			}

			var flattenedEntities []TableEntityDataSourceModel
			for _, entity := range result.Entities {
				flattenedEntity := flattenEntityWithMetadata(entity)
				if len(flattenedEntity.Properties) == 0 {
					// if we use selector, we get empty objects back, skip them
					continue
				}
				flattenedEntities = append(flattenedEntities, flattenedEntity)
			}
			model.Items = flattenedEntities
			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

// The api returns extra information that we already have. We'll remove it here before setting it in state.
func flattenEntityWithMetadata(entity map[string]interface{}) TableEntityDataSourceModel {
	delete(entity, "Timestamp")

	result := TableEntityDataSourceModel{}

	properties := map[string]interface{}{}
	for k, v := range entity {
		if k == "PartitionKey" {
			result.PartitionKey = v.(string)
			continue
		}

		if k == "RowKey" {
			result.RowKey = v.(string)
			continue
		}
		// skip ODATA annotation returned with fullmetadata
		if strings.HasPrefix(k, "odata.") || strings.HasSuffix(k, "@odata.type") {
			continue
		}
		if dtype, ok := entity[k+"@odata.type"]; ok {
			switch dtype {
			case "Edm.Boolean":
				properties[k] = fmt.Sprint(v)
			case "Edm.Double":
				properties[k] = fmt.Sprintf("%f", v)
			case "Edm.Int32", "Edm.Int64":
				// `v` returned as string for int 64
				properties[k] = fmt.Sprint(v)
			case "Edm.String":
				properties[k] = v
			default:
				log.Printf("[WARN] key %q with unexpected @odata.type %q", k, dtype)
				continue
			}

			properties[k+"@odata.type"] = dtype
		} else {
			// special handling for property types that do not require the annotation to be present
			// https://docs.microsoft.com/en-us/rest/api/storageservices/payload-format-for-table-service-operations#property-types-in-a-json-feed
			switch c := v.(type) {
			case bool:
				properties[k] = fmt.Sprint(v)
				properties[k+"@odata.type"] = "Edm.Boolean"
			case float64:
				f64 := v.(float64)
				if v == float64(int64(f64)) {
					properties[k] = fmt.Sprintf("%d", int64(f64))
					properties[k+"@odata.type"] = "Edm.Int32"
				} else {
					// fmt.Sprintf("%f", v) will return `123.123000` for `123.123`, have to use fmt.Sprint
					properties[k] = fmt.Sprint(v)
					properties[k+"@odata.type"] = "Edm.Double"
				}
			case string:
				properties[k] = v
			default:
				log.Printf("[WARN] key %q with unexpected type %T", k, c)
			}
		}
	}
	result.Properties = properties

	return result
}
