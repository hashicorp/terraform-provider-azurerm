
## `github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/geographichierarchies` Documentation

The `geographichierarchies` SDK allows for interaction with the Azure Resource Manager Service `trafficmanager` (API Version `2018-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/geographichierarchies"
```


### Client Initialization

```go
client := geographichierarchies.NewGeographicHierarchiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GeographicHierarchiesClient.GetDefault`

```go
ctx := context.TODO()


read, err := client.GetDefault(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
