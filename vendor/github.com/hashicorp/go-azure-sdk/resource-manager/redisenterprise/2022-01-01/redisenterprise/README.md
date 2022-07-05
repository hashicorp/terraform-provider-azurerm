
## `github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/redisenterprise` Documentation

The `redisenterprise` SDK allows for interaction with the Azure Resource Manager Service `redisenterprise` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/redisenterprise"
```


### Client Initialization

```go
client := redisenterprise.NewRedisEnterpriseClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisEnterpriseClient.Create`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := redisenterprise.Cluster{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesCreate`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.Database{
	// ...
}


if err := client.DatabasesCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesDelete`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

if err := client.DatabasesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesExport`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.ExportClusterParameters{
	// ...
}


if err := client.DatabasesExportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesForceUnlink`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.ForceUnlinkParameters{
	// ...
}


if err := client.DatabasesForceUnlinkThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesGet`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.DatabasesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesImport`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.ImportClusterParameters{
	// ...
}


if err := client.DatabasesImportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesListByCluster`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.DatabasesListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.DatabasesListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesListKeys`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.DatabasesListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesRegenerateKey`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.RegenerateKeyParameters{
	// ...
}


if err := client.DatabasesRegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesUpdate`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := redisenterprise.DatabaseUpdate{
	// ...
}


if err := client.DatabasesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.Delete`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.Get`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisEnterpriseClient.List`

```go
ctx := context.TODO()
id := redisenterprise.NewSubscriptionID()

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisEnterpriseClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := redisenterprise.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisEnterpriseClient.Update`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := redisenterprise.ClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
