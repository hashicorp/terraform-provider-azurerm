
## `github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource` Documentation

The `grafanaresource` SDK allows for interaction with Azure Resource Manager `dashboard` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource"
```


### Client Initialization

```go
client := grafanaresource.NewGrafanaResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GrafanaResourceClient.GrafanaCheckEnterpriseDetails`

```go
ctx := context.TODO()
id := grafanaresource.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

read, err := client.GrafanaCheckEnterpriseDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GrafanaResourceClient.GrafanaCreate`

```go
ctx := context.TODO()
id := grafanaresource.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

payload := grafanaresource.ManagedGrafana{
	// ...
}


if err := client.GrafanaCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GrafanaResourceClient.GrafanaDelete`

```go
ctx := context.TODO()
id := grafanaresource.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

if err := client.GrafanaDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GrafanaResourceClient.GrafanaGet`

```go
ctx := context.TODO()
id := grafanaresource.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

read, err := client.GrafanaGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GrafanaResourceClient.GrafanaList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.GrafanaList(ctx, id)` can be used to do batched pagination
items, err := client.GrafanaListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GrafanaResourceClient.GrafanaListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.GrafanaListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.GrafanaListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GrafanaResourceClient.GrafanaUpdate`

```go
ctx := context.TODO()
id := grafanaresource.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

payload := grafanaresource.ManagedGrafanaUpdateParameters{
	// ...
}


read, err := client.GrafanaUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
