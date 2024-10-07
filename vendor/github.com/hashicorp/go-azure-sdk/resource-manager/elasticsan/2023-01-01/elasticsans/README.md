
## `github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans` Documentation

The `elasticsans` SDK allows for interaction with Azure Resource Manager `elasticsan` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
```


### Client Initialization

```go
client := elasticsans.NewElasticSansClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ElasticSansClient.Create`

```go
ctx := context.TODO()
id := elasticsans.NewElasticSanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName")

payload := elasticsans.ElasticSan{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticSansClient.Delete`

```go
ctx := context.TODO()
id := elasticsans.NewElasticSanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticSansClient.Get`

```go
ctx := context.TODO()
id := elasticsans.NewElasticSanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticSansClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticSansClient.Update`

```go
ctx := context.TODO()
id := elasticsans.NewElasticSanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName")

payload := elasticsans.ElasticSanUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
