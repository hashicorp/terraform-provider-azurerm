
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases` Documentation

The `databases` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-02-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases"
```


### Client Initialization

```go
client := databases.NewDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasesClient.CreateOrUpdate`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

payload := databases.Database{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Export`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

payload := databases.ExportDatabaseDefinition{
	// ...
}


if err := client.ExportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Failover`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

if err := client.FailoverThenPoll(ctx, id, databases.DefaultFailoverOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

read, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.Import`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

payload := databases.ImportExistingDatabaseDefinition{
	// ...
}


if err := client.ImportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.ListByElasticPool`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "elasticPoolValue")
=======
id := databases.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "elasticPoolValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

// alternatively `client.ListByElasticPool(ctx, id)` can be used to do batched pagination
items, err := client.ListByElasticPoolComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatabasesClient.ListByServer`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
=======
id := databases.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatabasesClient.ListInaccessibleByServer`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
=======
id := databases.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

// alternatively `client.ListInaccessibleByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListInaccessibleByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatabasesClient.Pause`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

if err := client.PauseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Rename`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

payload := databases.ResourceMoveDefinition{
	// ...
}


read, err := client.Rename(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.Resume`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

if err := client.ResumeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Update`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

payload := databases.DatabaseUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.UpgradeDataWarehouse`

```go
ctx := context.TODO()
<<<<<<< HEAD
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
=======
id := databases.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c

if err := client.UpgradeDataWarehouseThenPoll(ctx, id); err != nil {
	// handle the error
}
```
