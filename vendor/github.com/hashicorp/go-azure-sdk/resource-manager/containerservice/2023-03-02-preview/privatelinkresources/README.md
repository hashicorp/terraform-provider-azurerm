
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/privatelinkresources` Documentation

The `privatelinkresources` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2023-03-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/privatelinkresources"
```


### Client Initialization

```go
client := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkResourcesClient.List`

```go
ctx := context.TODO()
id := privatelinkresources.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
