package entities

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/tables"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageTableEntity = Client{}

func TestEntitiesLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	tableName := fmt.Sprintf("table%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	tablesClient, err := tables.NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "table", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(tablesClient.Client, testData, auth.SharedKeyTable); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Creating Table..")
	if _, err := tablesClient.Create(ctx, tableName); err != nil {
		t.Fatalf("Error creating Table %q: %s", tableName, err)
	}
	defer tablesClient.Delete(ctx, tableName)

	entitiesClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "table", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(entitiesClient.Client, testData, auth.SharedKeyTable); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	partitionKey := "hello"
	rowKey := "there"

	t.Logf("[DEBUG] Inserting..")
	insertInput := InsertEntityInput{
		MetaDataLevel: NoMetaData,
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
		Entity: map[string]interface{}{
			"hello": "world",
		},
	}
	if _, err := entitiesClient.Insert(ctx, tableName, insertInput); err != nil {
		t.Logf("Error retrieving: %s", err)
	}

	t.Logf("[DEBUG] Insert or Merging..")
	insertOrMergeInput := InsertOrMergeEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity: map[string]interface{}{
			"hello": "ther88e",
		},
	}
	if _, err := entitiesClient.InsertOrMerge(ctx, tableName, insertOrMergeInput); err != nil {
		t.Logf("Error insert/merging: %s", err)
	}

	t.Logf("[DEBUG] Insert or Replacing..")
	insertOrReplaceInput := InsertOrReplaceEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Entity: map[string]interface{}{
			"hello": "pandas",
		},
	}
	if _, err := entitiesClient.InsertOrReplace(ctx, tableName, insertOrReplaceInput); err != nil {
		t.Logf("Error inserting/replacing: %s", err)
	}

	t.Logf("[DEBUG] Querying..")
	queryInput := QueryEntitiesInput{
		MetaDataLevel: NoMetaData,
	}
	results, err := entitiesClient.Query(ctx, tableName, queryInput)
	if err != nil {
		t.Logf("Error querying: %s", err)
	}

	if len(results.Entities) != 1 {
		t.Fatalf("Expected 1 item but got %d", len(results.Entities))
	}

	for _, v := range results.Entities {
		thisPartitionKey := v["PartitionKey"].(string)
		thisRowKey := v["RowKey"].(string)
		if partitionKey != thisPartitionKey {
			t.Fatalf("Expected Partition Key to be %q but got %q", partitionKey, thisPartitionKey)
		}
		if rowKey != thisRowKey {
			t.Fatalf("Expected Partition Key to be %q but got %q", rowKey, thisRowKey)
		}
	}

	t.Logf("[DEBUG] Retrieving..")
	getInput := GetEntityInput{
		MetaDataLevel: MinimalMetaData,
		PartitionKey:  partitionKey,
		RowKey:        rowKey,
	}
	getResults, err := entitiesClient.Get(ctx, tableName, getInput)
	if err != nil {
		t.Logf("Error querying: %s", err)
	}

	partitionKey2 := getResults.Entity["PartitionKey"].(string)
	rowKey2 := getResults.Entity["RowKey"].(string)
	if partitionKey2 != partitionKey {
		t.Fatalf("Expected Partition Key to be %q but got %q", partitionKey, partitionKey2)
	}
	if rowKey2 != rowKey {
		t.Fatalf("Expected Row Key to be %q but got %q", rowKey, rowKey2)
	}

	t.Logf("[DEBUG] Deleting..")
	deleteInput := DeleteEntityInput{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
	}
	if _, err := entitiesClient.Delete(ctx, tableName, deleteInput); err != nil {
		t.Logf("Error deleting: %s", err)
	}
}
