
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/containerservices` Documentation

The `containerservices` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2019-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/containerservices"
```


### Client Initialization

```go
client := containerservices.NewContainerServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerServicesClient.ListOrchestrators`

```go
ctx := context.TODO()
id := containerservices.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.ListOrchestrators(ctx, id, containerservices.DefaultListOrchestratorsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
