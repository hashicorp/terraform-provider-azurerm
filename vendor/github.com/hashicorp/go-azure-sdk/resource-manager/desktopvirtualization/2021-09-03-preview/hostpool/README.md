
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/hostpool` Documentation

The `hostpool` SDK allows for interaction with the Azure Resource Manager Service `desktopvirtualization` (API Version `2021-09-03-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/hostpool"
```


### Client Initialization

```go
client := hostpool.NewHostPoolClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HostPoolClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := hostpool.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

payload := hostpool.HostPool{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HostPoolClient.Delete`

```go
ctx := context.TODO()
id := hostpool.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

read, err := client.Delete(ctx, id, hostpool.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HostPoolClient.Get`

```go
ctx := context.TODO()
id := hostpool.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HostPoolClient.List`

```go
ctx := context.TODO()
id := hostpool.NewSubscriptionID()

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HostPoolClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := hostpool.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HostPoolClient.RetrieveRegistrationToken`

```go
ctx := context.TODO()
id := hostpool.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

read, err := client.RetrieveRegistrationToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HostPoolClient.Update`

```go
ctx := context.TODO()
id := hostpool.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

payload := hostpool.HostPoolPatch{
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
