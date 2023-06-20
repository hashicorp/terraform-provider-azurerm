
## `github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/databases` Documentation

The `databases` SDK allows for interaction with the Azure Resource Manager Service `redisenterprise` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/databases"
```


### Client Initialization

```go
client := databases.NewDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasesClient.Create`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.Database{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Export`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.ExportClusterParameters{
	// ...
}


if err := client.ExportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.ForceUnlink`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.ForceUnlinkParameters{
	// ...
}


if err := client.ForceUnlinkThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

read, err := client.Get(ctx, id)
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
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.ImportClusterParameters{
	// ...
}


if err := client.ImportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.ListByCluster`

```go
ctx := context.TODO()
id := databases.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatabasesClient.ListKeys`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.RegenerateKey`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.RegenerateKeyParameters{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Update`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseValue", "databaseValue")

payload := databases.DatabaseUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
