
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker` Documentation

The `servicelinker` SDK allows for interaction with Azure Resource Manager `servicelinker` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
```


### Client Initialization

```go
client := servicelinker.NewServicelinkerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServicelinkerClient.ConnectorCreateOrUpdate`

```go
ctx := context.TODO()
id := servicelinker.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "connectorName")

payload := servicelinker.LinkerResource{
	// ...
}


if err := client.ConnectorCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServicelinkerClient.ConnectorGet`

```go
ctx := context.TODO()
id := servicelinker.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "connectorName")

read, err := client.ConnectorGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServicelinkerClient.ConnectorList`

```go
ctx := context.TODO()
id := servicelinker.NewLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.ConnectorList(ctx, id)` can be used to do batched pagination
items, err := client.ConnectorListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServicelinkerClient.LinkerCreateOrUpdate`

```go
ctx := context.TODO()
id := servicelinker.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

payload := servicelinker.LinkerResource{
	// ...
}


if err := client.LinkerCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServicelinkerClient.LinkerGet`

```go
ctx := context.TODO()
id := servicelinker.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

read, err := client.LinkerGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServicelinkerClient.LinkerList`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.LinkerList(ctx, id)` can be used to do batched pagination
items, err := client.LinkerListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
