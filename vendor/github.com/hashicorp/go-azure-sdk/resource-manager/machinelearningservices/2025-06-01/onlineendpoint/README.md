
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/onlineendpoint` Documentation

The `onlineendpoint` SDK allows for interaction with Azure Resource Manager `machinelearningservices` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/onlineendpoint"
```


### Client Initialization

```go
client := onlineendpoint.NewOnlineEndpointClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OnlineEndpointClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

payload := onlineendpoint.OnlineEndpointTrackedResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OnlineEndpointClient.Delete`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OnlineEndpointClient.Get`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OnlineEndpointClient.GetToken`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

read, err := client.GetToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OnlineEndpointClient.List`

```go
ctx := context.TODO()
id := onlineendpoint.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.List(ctx, id, onlineendpoint.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, onlineendpoint.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OnlineEndpointClient.ListKeys`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OnlineEndpointClient.RegenerateKeys`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

payload := onlineendpoint.RegenerateEndpointKeysRequest{
	// ...
}


if err := client.RegenerateKeysThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OnlineEndpointClient.Update`

```go
ctx := context.TODO()
id := onlineendpoint.NewOnlineEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "onlineEndpointName")

payload := onlineendpoint.PartialMinimalTrackedResourceWithIdentity{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
