
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/desktop` Documentation

The `desktop` SDK allows for interaction with Azure Resource Manager `desktopvirtualization` (API Version `2022-02-10-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/desktop"
```


### Client Initialization

```go
client := desktop.NewDesktopClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DesktopClient.Get`

```go
ctx := context.TODO()
id := desktop.NewDesktopID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGroupName", "desktopName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DesktopClient.List`

```go
ctx := context.TODO()
id := desktop.NewApplicationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGroupName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DesktopClient.Update`

```go
ctx := context.TODO()
id := desktop.NewDesktopID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGroupName", "desktopName")

payload := desktop.DesktopPatch{
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
