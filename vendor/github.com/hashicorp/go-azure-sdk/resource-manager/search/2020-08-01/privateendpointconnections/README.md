
## `github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-08-01/privateendpointconnections` Documentation

The `privateendpointconnections` SDK allows for interaction with the Azure Resource Manager Service `search` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-08-01/privateendpointconnections"
```


### Client Initialization

```go
client := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateEndpointConnectionsClient.Delete`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue", "privateEndpointConnectionValue")

read, err := client.Delete(ctx, id, privateendpointconnections.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointConnectionsClient.Get`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue", "privateEndpointConnectionValue")

read, err := client.Get(ctx, id, privateendpointconnections.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointConnectionsClient.ListByService`

```go
ctx := context.TODO()
id := privateendpointconnections.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

// alternatively `client.ListByService(ctx, id, privateendpointconnections.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, privateendpointconnections.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateEndpointConnectionsClient.Update`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue", "privateEndpointConnectionValue")

payload := privateendpointconnections.PrivateEndpointConnection{
	// ...
}


read, err := client.Update(ctx, id, payload, privateendpointconnections.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
