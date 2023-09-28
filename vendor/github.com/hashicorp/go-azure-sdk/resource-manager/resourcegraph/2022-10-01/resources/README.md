
## `github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/resources` Documentation

The `resources` SDK allows for interaction with the Azure Resource Manager Service `resourcegraph` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/resourcegraph/2022-10-01/resources"
```


### Client Initialization

```go
client := resources.NewResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourcesClient.Resources`

```go
ctx := context.TODO()

payload := resources.QueryRequest{
	// ...
}


read, err := client.Resources(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
