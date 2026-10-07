
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/restorables` Documentation

The `restorables` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2024-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/restorables"
```


### Client Initialization

```go
client := restorables.NewRestorablesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RestorablesClient.GremlinResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := restorables.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

payload := restorables.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.GremlinResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RestorablesClient.MongoDBResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := restorables.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

payload := restorables.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.MongoDBResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RestorablesClient.RestorableDatabaseAccountsGetByLocation`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableDatabaseAccountsGetByLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableDatabaseAccountsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.RestorableDatabaseAccountsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableDatabaseAccountsListByLocation`

```go
ctx := context.TODO()
id := restorables.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.RestorableDatabaseAccountsListByLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableGremlinDatabasesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableGremlinDatabasesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableGremlinGraphsList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableGremlinGraphsList(ctx, id, restorables.DefaultRestorableGremlinGraphsListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableGremlinResourcesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableGremlinResourcesList(ctx, id, restorables.DefaultRestorableGremlinResourcesListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableMongodbCollectionsList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableMongodbCollectionsList(ctx, id, restorables.DefaultRestorableMongodbCollectionsListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableMongodbDatabasesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableMongodbDatabasesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableMongodbResourcesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableMongodbResourcesList(ctx, id, restorables.DefaultRestorableMongodbResourcesListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableSqlContainersList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableSqlContainersList(ctx, id, restorables.DefaultRestorableSqlContainersListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableSqlDatabasesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableSqlDatabasesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableSqlResourcesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableSqlResourcesList(ctx, id, restorables.DefaultRestorableSqlResourcesListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableTableResourcesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableTableResourcesList(ctx, id, restorables.DefaultRestorableTableResourcesListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.RestorableTablesList`

```go
ctx := context.TODO()
id := restorables.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableTablesList(ctx, id, restorables.DefaultRestorableTablesListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorablesClient.SqlResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := restorables.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

payload := restorables.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.SqlResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RestorablesClient.TableResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := restorables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

payload := restorables.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.TableResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
