
## `github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces` Documentation

The `playwrightworkspaces` SDK allows for interaction with Azure Resource Manager `loadtestservice` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
```


### Client Initialization

```go
client := playwrightworkspaces.NewPlaywrightWorkspacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PlaywrightWorkspacesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := playwrightworkspaces.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PlaywrightWorkspacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := playwrightworkspaces.NewPlaywrightWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "playwrightWorkspaceName")

payload := playwrightworkspaces.PlaywrightWorkspace{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PlaywrightWorkspacesClient.Delete`

```go
ctx := context.TODO()
id := playwrightworkspaces.NewPlaywrightWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "playwrightWorkspaceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PlaywrightWorkspacesClient.Get`

```go
ctx := context.TODO()
id := playwrightworkspaces.NewPlaywrightWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "playwrightWorkspaceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PlaywrightWorkspacesClient.ListByResourceGroup`

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


### Example Usage: `PlaywrightWorkspacesClient.ListBySubscription`

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


### Example Usage: `PlaywrightWorkspacesClient.Update`

```go
ctx := context.TODO()
id := playwrightworkspaces.NewPlaywrightWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "playwrightWorkspaceName")

payload := playwrightworkspaces.PlaywrightWorkspaceUpdate{
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
