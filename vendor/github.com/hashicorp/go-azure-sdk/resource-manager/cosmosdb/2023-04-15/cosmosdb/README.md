
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/cosmosdb` Documentation

The `cosmosdb` SDK allows for interaction with the Azure Resource Manager Service `cosmosdb` (API Version `2023-04-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/cosmosdb"
```


### Client Initialization

```go
client := cosmosdb.NewCosmosDBClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CosmosDBClient.CassandraResourcesCreateUpdateCassandraKeyspace`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

payload := cosmosdb.CassandraKeyspaceCreateUpdateParameters{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraKeyspaceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesCreateUpdateCassandraTable`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

payload := cosmosdb.CassandraTableCreateUpdateParameters{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesDeleteCassandraKeyspace`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

if err := client.CassandraResourcesDeleteCassandraKeyspaceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesDeleteCassandraTable`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

if err := client.CassandraResourcesDeleteCassandraTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesGetCassandraKeyspace`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

read, err := client.CassandraResourcesGetCassandraKeyspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesGetCassandraKeyspaceThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

read, err := client.CassandraResourcesGetCassandraKeyspaceThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesGetCassandraTable`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

read, err := client.CassandraResourcesGetCassandraTable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesGetCassandraTableThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

read, err := client.CassandraResourcesGetCassandraTableThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesListCassandraKeyspaces`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.CassandraResourcesListCassandraKeyspaces(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesListCassandraTables`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

read, err := client.CassandraResourcesListCassandraTables(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesMigrateCassandraKeyspaceToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

if err := client.CassandraResourcesMigrateCassandraKeyspaceToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesMigrateCassandraKeyspaceToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

if err := client.CassandraResourcesMigrateCassandraKeyspaceToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesMigrateCassandraTableToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

if err := client.CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesMigrateCassandraTableToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

if err := client.CassandraResourcesMigrateCassandraTableToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesUpdateCassandraKeyspaceThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.CassandraResourcesUpdateCassandraKeyspaceThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CassandraResourcesUpdateCassandraTableThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "cassandraKeyspaceValue", "tableValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.CassandraResourcesUpdateCassandraTableThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.CollectionListMetricDefinitions`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue")

read, err := client.CollectionListMetricDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue")

read, err := client.CollectionListMetrics(ctx, id, cosmosdb.DefaultCollectionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionListUsages`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue")

read, err := client.CollectionListUsages(ctx, id, cosmosdb.DefaultCollectionListUsagesOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionPartitionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue")

read, err := client.CollectionPartitionListMetrics(ctx, id, cosmosdb.DefaultCollectionPartitionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionPartitionListUsages`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue")

read, err := client.CollectionPartitionListUsages(ctx, id, cosmosdb.DefaultCollectionPartitionListUsagesOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionPartitionRegionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "regionValue", "databaseValue", "collectionValue")

read, err := client.CollectionPartitionRegionListMetrics(ctx, id, cosmosdb.DefaultCollectionPartitionRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.CollectionRegionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "regionValue", "databaseValue", "collectionValue")

read, err := client.CollectionRegionListMetrics(ctx, id, cosmosdb.DefaultCollectionRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountRegionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "regionValue")

read, err := client.DatabaseAccountRegionListMetrics(ctx, id, cosmosdb.DefaultDatabaseAccountRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsCheckNameExists`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountNameID("databaseAccountValue")

read, err := client.DatabaseAccountsCheckNameExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsCreateOrUpdate`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.DatabaseAccountCreateUpdateParameters{
	// ...
}


if err := client.DatabaseAccountsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsDelete`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

if err := client.DatabaseAccountsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsFailoverPriorityChange`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.FailoverPolicies{
	// ...
}


if err := client.DatabaseAccountsFailoverPriorityChangeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsGet`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsGetReadOnlyKeys`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsGetReadOnlyKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsList`

```go
ctx := context.TODO()
id := cosmosdb.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.DatabaseAccountsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListByResourceGroup`

```go
ctx := context.TODO()
id := cosmosdb.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.DatabaseAccountsListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListConnectionStrings`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListConnectionStrings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListKeys`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListMetricDefinitions`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListMetricDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListMetrics(ctx, id, cosmosdb.DefaultDatabaseAccountsListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListReadOnlyKeys`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListReadOnlyKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsListUsages`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.DatabaseAccountsListUsages(ctx, id, cosmosdb.DefaultDatabaseAccountsListUsagesOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsOfflineRegion`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.RegionForOnlineOffline{
	// ...
}


if err := client.DatabaseAccountsOfflineRegionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsOnlineRegion`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.RegionForOnlineOffline{
	// ...
}


if err := client.DatabaseAccountsOnlineRegionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsRegenerateKey`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.DatabaseAccountRegenerateKeyParameters{
	// ...
}


if err := client.DatabaseAccountsRegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseAccountsUpdate`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

payload := cosmosdb.DatabaseAccountUpdateParameters{
	// ...
}


if err := client.DatabaseAccountsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.DatabaseListMetricDefinitions`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue")

read, err := client.DatabaseListMetricDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue")

read, err := client.DatabaseListMetrics(ctx, id, cosmosdb.DefaultDatabaseListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.DatabaseListUsages`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue")

read, err := client.DatabaseListUsages(ctx, id, cosmosdb.DefaultDatabaseListUsagesOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesCreateUpdateGremlinDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

payload := cosmosdb.GremlinDatabaseCreateUpdateParameters{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesCreateUpdateGremlinGraph`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

payload := cosmosdb.GremlinGraphCreateUpdateParameters{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinGraphThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesDeleteGremlinDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

if err := client.GremlinResourcesDeleteGremlinDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesDeleteGremlinGraph`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

if err := client.GremlinResourcesDeleteGremlinGraphThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesGetGremlinDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

read, err := client.GremlinResourcesGetGremlinDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesGetGremlinDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

read, err := client.GremlinResourcesGetGremlinDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesGetGremlinGraph`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

read, err := client.GremlinResourcesGetGremlinGraph(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesGetGremlinGraphThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

read, err := client.GremlinResourcesGetGremlinGraphThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesListGremlinDatabases`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.GremlinResourcesListGremlinDatabases(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesListGremlinGraphs`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

read, err := client.GremlinResourcesListGremlinGraphs(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesMigrateGremlinDatabaseToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

if err := client.GremlinResourcesMigrateGremlinDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesMigrateGremlinDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

if err := client.GremlinResourcesMigrateGremlinDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesMigrateGremlinGraphToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

if err := client.GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesMigrateGremlinGraphToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

if err := client.GremlinResourcesMigrateGremlinGraphToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesUpdateGremlinDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.GremlinResourcesUpdateGremlinGraphThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "gremlinDatabaseValue", "graphValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.GremlinResourcesUpdateGremlinGraphThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.LocationsGet`

```go
ctx := context.TODO()
id := cosmosdb.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.LocationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.LocationsList`

```go
ctx := context.TODO()
id := cosmosdb.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.LocationsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesCreateUpdateMongoDBCollection`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

payload := cosmosdb.MongoDBCollectionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesCreateUpdateMongoDBDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

payload := cosmosdb.MongoDBDatabaseCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoDBDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesDeleteMongoDBCollection`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

if err := client.MongoDBResourcesDeleteMongoDBCollectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesDeleteMongoDBDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

if err := client.MongoDBResourcesDeleteMongoDBDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesGetMongoDBCollection`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

read, err := client.MongoDBResourcesGetMongoDBCollection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesGetMongoDBCollectionThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

read, err := client.MongoDBResourcesGetMongoDBCollectionThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesGetMongoDBDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

read, err := client.MongoDBResourcesGetMongoDBDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesGetMongoDBDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

read, err := client.MongoDBResourcesGetMongoDBDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesListMongoDBCollections`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

read, err := client.MongoDBResourcesListMongoDBCollections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesListMongoDBDatabases`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.MongoDBResourcesListMongoDBDatabases(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesMigrateMongoDBCollectionToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

if err := client.MongoDBResourcesMigrateMongoDBCollectionToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesMigrateMongoDBCollectionToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

if err := client.MongoDBResourcesMigrateMongoDBCollectionToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesMigrateMongoDBDatabaseToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

if err := client.MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

if err := client.MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesUpdateMongoDBCollectionThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue", "collectionValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesUpdateMongoDBCollectionThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.MongoDBResourcesUpdateMongoDBDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "mongodbDatabaseValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesUpdateMongoDBDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.PartitionKeyRangeIdListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewPartitionKeyRangeIdID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "databaseValue", "collectionValue", "partitionKeyRangeIdValue")

read, err := client.PartitionKeyRangeIdListMetrics(ctx, id, cosmosdb.DefaultPartitionKeyRangeIdListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.PartitionKeyRangeIdRegionListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewCollectionPartitionKeyRangeIdID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "regionValue", "databaseValue", "collectionValue", "partitionKeyRangeIdValue")

read, err := client.PartitionKeyRangeIdRegionListMetrics(ctx, id, cosmosdb.DefaultPartitionKeyRangeIdRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.PercentileListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.PercentileListMetrics(ctx, id, cosmosdb.DefaultPercentileListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.PercentileSourceTargetListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewSourceRegionTargetRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sourceRegionValue", "targetRegionValue")

read, err := client.PercentileSourceTargetListMetrics(ctx, id, cosmosdb.DefaultPercentileSourceTargetListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.PercentileTargetListMetrics`

```go
ctx := context.TODO()
id := cosmosdb.NewTargetRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "targetRegionValue")

read, err := client.PercentileTargetListMetrics(ctx, id, cosmosdb.DefaultPercentileTargetListMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateClientEncryptionKey`

```go
ctx := context.TODO()
id := cosmosdb.NewClientEncryptionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "clientEncryptionKeyValue")

payload := cosmosdb.ClientEncryptionKeyCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateClientEncryptionKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateSqlContainer`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

payload := cosmosdb.SqlContainerCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlContainerThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateSqlDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

payload := cosmosdb.SqlDatabaseCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateSqlStoredProcedure`

```go
ctx := context.TODO()
id := cosmosdb.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "storedProcedureValue")

payload := cosmosdb.SqlStoredProcedureCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateSqlTrigger`

```go
ctx := context.TODO()
id := cosmosdb.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "triggerValue")

payload := cosmosdb.SqlTriggerCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlTriggerThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesCreateUpdateSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := cosmosdb.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "userDefinedFunctionValue")

payload := cosmosdb.SqlUserDefinedFunctionCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesDeleteSqlContainer`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

if err := client.SqlResourcesDeleteSqlContainerThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesDeleteSqlDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

if err := client.SqlResourcesDeleteSqlDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesDeleteSqlStoredProcedure`

```go
ctx := context.TODO()
id := cosmosdb.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "storedProcedureValue")

if err := client.SqlResourcesDeleteSqlStoredProcedureThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesDeleteSqlTrigger`

```go
ctx := context.TODO()
id := cosmosdb.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "triggerValue")

if err := client.SqlResourcesDeleteSqlTriggerThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesDeleteSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := cosmosdb.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "userDefinedFunctionValue")

if err := client.SqlResourcesDeleteSqlUserDefinedFunctionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetClientEncryptionKey`

```go
ctx := context.TODO()
id := cosmosdb.NewClientEncryptionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "clientEncryptionKeyValue")

read, err := client.SqlResourcesGetClientEncryptionKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlContainer`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

read, err := client.SqlResourcesGetSqlContainer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlContainerThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

read, err := client.SqlResourcesGetSqlContainerThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlDatabase`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

read, err := client.SqlResourcesGetSqlDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

read, err := client.SqlResourcesGetSqlDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlStoredProcedure`

```go
ctx := context.TODO()
id := cosmosdb.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "storedProcedureValue")

read, err := client.SqlResourcesGetSqlStoredProcedure(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlTrigger`

```go
ctx := context.TODO()
id := cosmosdb.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "triggerValue")

read, err := client.SqlResourcesGetSqlTrigger(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesGetSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := cosmosdb.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue", "userDefinedFunctionValue")

read, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListClientEncryptionKeys`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

read, err := client.SqlResourcesListClientEncryptionKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListSqlContainers`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

read, err := client.SqlResourcesListSqlContainers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListSqlDatabases`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.SqlResourcesListSqlDatabases(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListSqlStoredProcedures`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

read, err := client.SqlResourcesListSqlStoredProcedures(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListSqlTriggers`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

read, err := client.SqlResourcesListSqlTriggers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesListSqlUserDefinedFunctions`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

read, err := client.SqlResourcesListSqlUserDefinedFunctions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.SqlResourcesMigrateSqlContainerToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

if err := client.SqlResourcesMigrateSqlContainerToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesMigrateSqlContainerToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

if err := client.SqlResourcesMigrateSqlContainerToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesMigrateSqlDatabaseToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

if err := client.SqlResourcesMigrateSqlDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesMigrateSqlDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

if err := client.SqlResourcesMigrateSqlDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesUpdateSqlContainerThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue", "containerValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.SqlResourcesUpdateSqlContainerThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.SqlResourcesUpdateSqlDatabaseThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "sqlDatabaseValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.SqlResourcesUpdateSqlDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.TableResourcesCreateUpdateTable`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

payload := cosmosdb.TableCreateUpdateParameters{
	// ...
}


if err := client.TableResourcesCreateUpdateTableThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.TableResourcesDeleteTable`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

if err := client.TableResourcesDeleteTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.TableResourcesGetTable`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

read, err := client.TableResourcesGetTable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.TableResourcesGetTableThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

read, err := client.TableResourcesGetTableThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.TableResourcesListTables`

```go
ctx := context.TODO()
id := cosmosdb.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue")

read, err := client.TableResourcesListTables(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CosmosDBClient.TableResourcesMigrateTableToAutoscale`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

if err := client.TableResourcesMigrateTableToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.TableResourcesMigrateTableToManualThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

if err := client.TableResourcesMigrateTableToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CosmosDBClient.TableResourcesUpdateTableThroughput`

```go
ctx := context.TODO()
id := cosmosdb.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "tableValue")

payload := cosmosdb.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.TableResourcesUpdateTableThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
