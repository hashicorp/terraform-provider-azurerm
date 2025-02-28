
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/privateendpointconnections` Documentation

The `privateendpointconnections` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/privateendpointconnections"
```


### Client Initialization

```go
client := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionCreateOrUpdate`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "privateEndpointConnectionName")

payload := privateendpointconnections.PrivateEndpointConnectionRequest{
	// ...
}


if err := client.PrivateEndpointConnectionCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionDelete`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "privateEndpointConnectionName")

if err := client.PrivateEndpointConnectionDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionGetByName`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "privateEndpointConnectionName")

read, err := client.PrivateEndpointConnectionGetByName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionGetPrivateLinkResource`

```go
ctx := context.TODO()
id := privateendpointconnections.NewPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "privateLinkResourceName")

read, err := client.PrivateEndpointConnectionGetPrivateLinkResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionListByService`

```go
ctx := context.TODO()
id := privateendpointconnections.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.PrivateEndpointConnectionListByService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointConnectionsClient.PrivateEndpointConnectionListPrivateLinkResources`

```go
ctx := context.TODO()
id := privateendpointconnections.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.PrivateEndpointConnectionListPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
