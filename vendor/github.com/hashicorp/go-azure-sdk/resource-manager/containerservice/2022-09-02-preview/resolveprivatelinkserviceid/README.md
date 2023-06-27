
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/resolveprivatelinkserviceid` Documentation

The `resolveprivatelinkserviceid` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2022-09-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/resolveprivatelinkserviceid"
```


### Client Initialization

```go
client := resolveprivatelinkserviceid.NewResolvePrivateLinkServiceIdClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResolvePrivateLinkServiceIdClient.POST`

```go
ctx := context.TODO()
id := resolveprivatelinkserviceid.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

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
