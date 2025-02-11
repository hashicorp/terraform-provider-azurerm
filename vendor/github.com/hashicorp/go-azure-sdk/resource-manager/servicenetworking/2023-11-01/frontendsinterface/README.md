
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/frontendsinterface` Documentation

The `frontendsinterface` SDK allows for interaction with Azure Resource Manager `servicenetworking` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/frontendsinterface"
```


### Client Initialization

```go
client := frontendsinterface.NewFrontendsInterfaceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FrontendsInterfaceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := frontendsinterface.NewFrontendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "frontendName")

payload := frontendsinterface.Frontend{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FrontendsInterfaceClient.Delete`

```go
ctx := context.TODO()
id := frontendsinterface.NewFrontendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "frontendName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FrontendsInterfaceClient.Get`

```go
ctx := context.TODO()
id := frontendsinterface.NewFrontendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "frontendName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FrontendsInterfaceClient.ListByTrafficController`

```go
ctx := context.TODO()
id := frontendsinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName")

// alternatively `client.ListByTrafficController(ctx, id)` can be used to do batched pagination
items, err := client.ListByTrafficControllerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FrontendsInterfaceClient.Update`

```go
ctx := context.TODO()
id := frontendsinterface.NewFrontendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "frontendName")

payload := frontendsinterface.FrontendUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
