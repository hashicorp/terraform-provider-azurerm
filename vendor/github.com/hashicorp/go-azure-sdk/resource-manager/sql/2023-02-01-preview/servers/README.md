
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers` Documentation

The `servers` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-02-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"
```


### Client Initialization

```go
client := servers.NewServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServersClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := servers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := servers.CheckNameAvailabilityRequest{
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


### Example Usage: `ServersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.Server{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Delete`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Get`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.Get(ctx, id, servers.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ImportDatabase`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.ImportNewDatabaseDefinition{
	// ...
}


if err := client.ImportDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.List`

```go
ctx := context.TODO()
id := servers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, servers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, servers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := servers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, servers.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, servers.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServersClient.RefreshStatus`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

if err := client.RefreshStatusThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Update`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.ServerUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
