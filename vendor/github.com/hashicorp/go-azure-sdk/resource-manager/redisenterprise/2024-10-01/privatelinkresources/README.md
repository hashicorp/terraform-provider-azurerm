
## `github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-10-01/privatelinkresources` Documentation

The `privatelinkresources` SDK allows for interaction with Azure Resource Manager `redisenterprise` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-10-01/privatelinkresources"
```


### Client Initialization

```go
client := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkResourcesClient.ListByCluster`

```go
ctx := context.TODO()
id := privatelinkresources.NewRedisEnterpriseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisEnterpriseName")

read, err := client.ListByCluster(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
