
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/resolveprivatelinkserviceid` Documentation

The `resolveprivatelinkserviceid` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2023-03-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/resolveprivatelinkserviceid"
```


### Client Initialization

```go
client := resolveprivatelinkserviceid.NewResolvePrivateLinkServiceIdClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResolvePrivateLinkServiceIdClient.POST`

```go
ctx := context.TODO()
id := commonids.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName")

payload := resolveprivatelinkserviceid.PrivateLinkResource{
	// ...
}


read, err := client.POST(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
