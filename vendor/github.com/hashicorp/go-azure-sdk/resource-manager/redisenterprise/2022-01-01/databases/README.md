
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
if err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Create`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.Database{
	// ...
}

future, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")
future, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Export`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.ExportClusterParameters{
	// ...
}

future, err := client.Export(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.ForceUnlink`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.ForceUnlinkParameters{
	// ...
}

future, err := client.ForceUnlink(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")
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
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.ImportClusterParameters{
	// ...
}

future, err := client.Import(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.ListByCluster`

```go
ctx := context.TODO()
id := databases.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")
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
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")
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
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.RegenerateKeyParameters{
	// ...
}

future, err := client.RegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Update`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := databases.DatabaseUpdate{
	// ...
}

future, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```
