
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers` Documentation

The `servers` SDK allows for interaction with Azure Resource Manager `postgresql` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
```


### Client Initialization

```go
client := servers.NewServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServersClient.BackupsLongTermRetentionCheckPrerequisites`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servers.BackupRequestBase{
	// ...
}


read, err := client.BackupsLongTermRetentionCheckPrerequisites(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.BackupsLongTermRetentionStart`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servers.BackupsLongTermRetentionRequest{
	// ...
}


if err := client.BackupsLongTermRetentionStartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.CapabilitiesByServerList`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.CapabilitiesByServerList(ctx, id)` can be used to do batched pagination
items, err := client.CapabilitiesByServerListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServersClient.CapturedLogsListByServer`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.CapturedLogsListByServer(ctx, id)` can be used to do batched pagination
items, err := client.CapturedLogsListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

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
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Get`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ListByResourceGroup`

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


### Example Usage: `ServersClient.ListBySubscription`

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


### Example Usage: `ServersClient.MigrationsCheckNameAvailability`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servers.MigrationNameAvailability{
	// ...
}


read, err := client.MigrationsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ReplicasListByServer`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ReplicasListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ReplicasListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServersClient.Restart`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servers.RestartParameter{
	// ...
}


if err := client.RestartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Start`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Stop`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Update`

```go
ctx := context.TODO()
id := servers.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servers.ServerForPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
