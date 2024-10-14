
## `github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-06-01-preview/redisenterprise` Documentation

The `redisenterprise` SDK allows for interaction with Azure Resource Manager `redisenterprise` (API Version `2024-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-06-01-preview/redisenterprise"
```


### Client Initialization

```go
client := redisenterprise.NewRedisEnterpriseClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisEnterpriseClient.Create`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

if err := client.DatabasesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesExport`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

payload := redisenterprise.ExportClusterParameters{
	// ...
}


if err := client.DatabasesExportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesFlush`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

payload := redisenterprise.FlushParameters{
	// ...
}


if err := client.DatabasesFlushThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesForceLinkToReplicationGroup`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

payload := redisenterprise.ForceLinkParameters{
	// ...
}


if err := client.DatabasesForceLinkToReplicationGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesForceUnlink`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

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
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

payload := redisenterprise.DatabaseUpdate{
	// ...
}


if err := client.DatabasesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.DatabasesUpgradeDBRedisVersion`

```go
ctx := context.TODO()
id := redisenterprise.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName", "databaseName")

if err := client.DatabasesUpgradeDBRedisVersionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.Delete`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisEnterpriseClient.Get`

```go
ctx := context.TODO()
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := redisenterprise.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

payload := redisenterprise.ClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
