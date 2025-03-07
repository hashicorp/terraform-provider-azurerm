
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/networkstatus` Documentation

The `networkstatus` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/networkstatus"
```


### Client Initialization

```go
client := networkstatus.NewNetworkStatusClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkStatusClient.ListByLocation`

```go
ctx := context.TODO()
id := networkstatus.NewLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "locationName")

read, err := client.ListByLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkStatusClient.ListByService`

```go
ctx := context.TODO()
id := networkstatus.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.ListByService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
