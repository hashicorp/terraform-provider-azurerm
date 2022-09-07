
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/servicelinker` Documentation

The `servicelinker` SDK allows for interaction with the Azure Resource Manager Service `servicelinker` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/servicelinker"
```


### Client Initialization

```go
client := servicelinker.NewServiceLinkerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceLinkerClient.LinkerCreateOrUpdate`

```go
ctx := context.TODO()
id := servicelinker.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerValue")

payload := servicelinker.LinkerResource{
	// ...
}


if err := client.LinkerCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceLinkerClient.LinkerGet`

```go
ctx := context.TODO()
id := servicelinker.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerValue")

read, err := client.LinkerGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServiceLinkerClient.LinkerList`

```go
ctx := context.TODO()
id := servicelinker.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.LinkerList(ctx, id)` can be used to do batched pagination
items, err := client.LinkerListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
