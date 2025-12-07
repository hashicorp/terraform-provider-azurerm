
## `github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas` Documentation

The `managedgrafanas` SDK allows for interaction with Azure Resource Manager `dashboard` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
```


### Client Initialization

```go
client := managedgrafanas.NewManagedGrafanasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedGrafanasClient.GrafanaCheckEnterpriseDetails`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

read, err := client.GrafanaCheckEnterpriseDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedGrafanasClient.GrafanaCreate`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

payload := managedgrafanas.ManagedGrafana{
	// ...
}


if err := client.GrafanaCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedGrafanasClient.GrafanaDelete`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

if err := client.GrafanaDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedGrafanasClient.GrafanaFetchAvailablePlugins`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

// alternatively `client.GrafanaFetchAvailablePlugins(ctx, id)` can be used to do batched pagination
items, err := client.GrafanaFetchAvailablePluginsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedGrafanasClient.GrafanaGet`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

read, err := client.GrafanaGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedGrafanasClient.GrafanaList`

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


### Example Usage: `ManagedGrafanasClient.GrafanaListByResourceGroup`

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


### Example Usage: `ManagedGrafanasClient.GrafanaUpdate`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

payload := managedgrafanas.ManagedGrafanaUpdateParameters{
	// ...
}


if err := client.GrafanaUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedGrafanasClient.ManagedPrivateEndpointsRefresh`

```go
ctx := context.TODO()
id := managedgrafanas.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

if err := client.ManagedPrivateEndpointsRefreshThenPoll(ctx, id); err != nil {
	// handle the error
}
```
