
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces` Documentation

The `azuremonitorworkspaces` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2023-04-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
```


### Client Initialization

```go
client := azuremonitorworkspaces.NewAzureMonitorWorkspacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AzureMonitorWorkspacesClient.Create`

```go
ctx := context.TODO()
id := azuremonitorworkspaces.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := azuremonitorworkspaces.AzureMonitorWorkspaceResource{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AzureMonitorWorkspacesClient.Delete`

```go
ctx := context.TODO()
id := azuremonitorworkspaces.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AzureMonitorWorkspacesClient.Get`

```go
ctx := context.TODO()
id := azuremonitorworkspaces.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AzureMonitorWorkspacesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureMonitorWorkspacesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureMonitorWorkspacesClient.Update`

```go
ctx := context.TODO()
id := azuremonitorworkspaces.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := azuremonitorworkspaces.AzureMonitorWorkspaceResourceForUpdate{
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
