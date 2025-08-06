
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/sessionhost` Documentation

The `sessionhost` SDK allows for interaction with Azure Resource Manager `desktopvirtualization` (API Version `2024-04-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/sessionhost"
```


### Client Initialization

```go
client := sessionhost.NewSessionHostClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SessionHostClient.Delete`

```go
ctx := context.TODO()
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "sessionHostName")

read, err := client.Delete(ctx, id, sessionhost.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SessionHostClient.Get`

```go
ctx := context.TODO()
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "sessionHostName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SessionHostClient.List`

```go
ctx := context.TODO()
id := sessionhost.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName")

// alternatively `client.List(ctx, id, sessionhost.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, sessionhost.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SessionHostClient.Update`

```go
ctx := context.TODO()
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName", "sessionHostName")

payload := sessionhost.SessionHostPatch{
	// ...
}


read, err := client.Update(ctx, id, payload, sessionhost.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
