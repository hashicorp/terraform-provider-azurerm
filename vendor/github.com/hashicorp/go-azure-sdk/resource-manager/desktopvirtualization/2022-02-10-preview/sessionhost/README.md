
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/sessionhost` Documentation

The `sessionhost` SDK allows for interaction with the Azure Resource Manager Service `desktopvirtualization` (API Version `2022-02-10-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/sessionhost"
```


### Client Initialization

```go
client := sessionhost.NewSessionHostClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SessionHostClient.Delete`

```go
ctx := context.TODO()
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue", "sessionHostValue")

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
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue", "sessionHostValue")

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
id := sessionhost.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
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
id := sessionhost.NewSessionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolValue", "sessionHostValue")

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
