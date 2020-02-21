package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	storage "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/entities"
)

func dataSourceArmStorageTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageTableRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key", "query"},
				ExactlyOneOf:  []string{"resource_id", "key", "query"},
			},

			"key": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"resource_id", "query"},
				ExactlyOneOf:  []string{"key", "query", "resource_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"partition_key": {
							Type:     schema.TypeString,
							Required: true,
						},

						"row_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"query": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"resource_id", "key"},
				ExactlyOneOf:  []string{"query", "resource_id", "key"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},

						"top": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
			},

			"entities": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

type Key struct {
	storageAccount string
	table          string
	partitionKey   string
	rowKey         string
}

type Query struct {
	storageAccount string
	table          string
	filter         string
	top            *int
}

func dataSourceArmStorageTableRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	storageClient := meta.(*clients.Client).Storage

	var results []map[string]interface{}

	resourceId := d.Get("resource_id").(string)
	if len(resourceId) > 0 {
		result, err := fetchById(resourceId, storageClient, ctx)
		if err != nil {
			return fmt.Errorf("Error retrieving Entity (ID %q): %s", resourceId, err)
		}

		d.SetId(resourceId)
		results = append(results, result.Entity) // wrap result Entity into slice
	}

	keyConfig := d.Get("key").(*schema.Set).List()
	if len(keyConfig) > 0 {
		key := Key{}
		for _, config := range keyConfig {
			cm := config.(map[string]interface{})

			cmSan := cm["storage_account_name"].(string)
			if cmSan != "" {
				key.storageAccount = cmSan
			}

			cmTn := cm["table_name"].(string)
			if cmTn != "" {
				key.table = cmTn
			}

			cmPk := cm["partition_key"].(string)
			if cmPk != "" {
				key.partitionKey = cmPk
			}

			cmRk := cm["row_key"].(string)
			if cmRk != "" {
				key.rowKey = cmRk
			}
		}

		resourceId := fmt.Sprintf("https://%s.table.core.windows.net/%s(PartitionKey='%s',RowKey='%s')", key.storageAccount, key.table, key.partitionKey, key.rowKey)
		result, err := fetchById(resourceId, storageClient, ctx)
		if err != nil {
			return fmt.Errorf("Error retrieving Entity (Storage Account %q, Table %q, partition key %q, row key %q): %s", key.storageAccount, key.table, key.partitionKey, key.rowKey, err)
		}

		d.SetId(resourceId)
		results = append(results, result.Entity) // wrap result Entity into slice
	}

	queryConfig := d.Get("query").(*schema.Set).List()
	if len(queryConfig) > 0 {
		query := Query{}
		for _, config := range queryConfig {
			cm := config.(map[string]interface{})

			cmSan := cm["storage_account_name"].(string)
			if cmSan != "" {
				query.storageAccount = cmSan
			}

			cmTn := cm["table_name"].(string)
			if cmTn != "" {
				query.table = cmTn
			}

			cmF := cm["filter"].(string)
			if cmF != "" {
				query.filter = cmF
			}

			cmT := cm["top"].(int)
			if cmT > -1 {
				query.top = &cmT
			}
		}

		result, err := fetchByQuery(query, storageClient, ctx)
		if err != nil {
			return fmt.Errorf("Error querying Entities (Storage Account %q, Table %q, query %+v): %s", query.storageAccount, query.table, query, err)
		}

		d.SetId(fmt.Sprintf("https://%s.table.core.windows.net/%s$filter=%s&$top=%d)", query.storageAccount, query.table, query.filter, query.top))
		results = result.Entities
	}

	if err := d.Set("entities", flattenEntities(results)); err != nil {
		return fmt.Errorf("Error setting `entity`: %s", err)
	}

	return nil
}

func getEntityClient(storageAccount string, storageClient *storage.Client, ctx context.Context) (*entities.Client, error) {
	account, err := storageClient.FindAccount(ctx, storageAccount)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Storage Account %q: %s", storageAccount, err)
	}
	if account == nil {
		return nil, fmt.Errorf("[WARN] Unable to determine Resource Group (Account %s) - assuming removed & removing from state", storageAccount)
	}

	return storageClient.TableEntityClient(ctx, *account)
}

func fetchById(resourceId string, storageClient *storage.Client, ctx context.Context) (*entities.GetEntityResult, error) {
	id, err := entities.ParseResourceID(resourceId)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Resource ID %q: %s", resourceId, err)
	}

	entityClient, err := getEntityClient(id.AccountName, storageClient, ctx)
	if err != nil {
		return nil, fmt.Errorf("Error building Table Entity entityClient for Storage Account %q: %s", id.AccountName, err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.NoMetaData,
	}
	result, err := entityClient.Get(ctx, id.AccountName, id.TableName, input)

	return &result, err
}

func fetchByQuery(query Query, storageClient *storage.Client, ctx context.Context) (*entities.QueryEntitiesResult, error) {
	entityClient, err := getEntityClient(query.storageAccount, storageClient, ctx)
	if err != nil {
		return nil, fmt.Errorf("Error building Table Entity entityClient for Storage Account %q: %s", query.storageAccount, err)
	}

	input := entities.QueryEntitiesInput{
		Filter:        &query.filter,
		Top:           query.top,
		MetaDataLevel: entities.NoMetaData,
	}
	result, err := entityClient.Query(ctx, query.storageAccount, query.table, input)

	return &result, err
}

func flattenEntities(entities []map[string]interface{}) []map[string]interface{} {
	for _, entity := range entities {
		delete(entity, "PartitionKey")
		delete(entity, "RowKey")
		delete(entity, "Timestamp")
	}

	return entities
}
