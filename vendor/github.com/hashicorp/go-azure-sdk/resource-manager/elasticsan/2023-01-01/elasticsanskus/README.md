
## `github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsanskus` Documentation

The `elasticsanskus` SDK allows for interaction with Azure Resource Manager `elasticsan` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsanskus"
```


### Client Initialization

```go
client := elasticsanskus.NewElasticSanSkusClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ElasticSanSkusClient.SkusList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.SkusList(ctx, id, elasticsanskus.DefaultSkusListOperationOptions())` can be used to do batched pagination
items, err := client.SkusListComplete(ctx, id, elasticsanskus.DefaultSkusListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
