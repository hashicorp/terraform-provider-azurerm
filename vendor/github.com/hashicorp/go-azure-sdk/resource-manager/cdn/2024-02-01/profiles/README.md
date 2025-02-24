
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles` Documentation

The `profiles` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
```


### Client Initialization

```go
client := profiles.NewProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProfilesClient.CanMigrate`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := profiles.CanMigrateParameters{
	// ...
}


if err := client.CanMigrateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.Create`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.Profile{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.Delete`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.GenerateSsoUri`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.GenerateSsoUri(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.Get`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProfilesClient.ListByResourceGroup`

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


### Example Usage: `ProfilesClient.ListResourceUsage`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.ListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProfilesClient.ListSupportedOptimizationTypes`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.ListSupportedOptimizationTypes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.Migrate`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := profiles.MigrationParameters{
	// ...
}


if err := client.MigrateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.MigrationCommit`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

if err := client.MigrationCommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.Update`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.ProfileUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
