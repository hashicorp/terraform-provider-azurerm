
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools` Documentation

The `elasticpools` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
```


### Client Initialization

```go
client := elasticpools.NewElasticPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ElasticPoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "elasticPoolName")

payload := elasticpools.ElasticPool{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticPoolsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "elasticPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticPoolsClient.Failover`

```go
ctx := context.TODO()
id := commonids.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "elasticPoolName")

if err := client.FailoverThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticPoolsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "elasticPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticPoolsClient.ListByServer`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

// alternatively `client.ListByServer(ctx, id, elasticpools.DefaultListByServerOperationOptions())` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id, elasticpools.DefaultListByServerOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticPoolsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewSqlElasticPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "elasticPoolName")

payload := elasticpools.ElasticPoolUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
