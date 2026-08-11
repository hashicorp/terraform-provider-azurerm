
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles` Documentation

The `profiles` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
```


### Client Initialization

```go
client := profiles.NewProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProfilesClient.AFDProfilesCheckEndpointNameAvailability`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.CheckEndpointNameAvailabilityInput{
	// ...
}


read, err := client.AFDProfilesCheckEndpointNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.AFDProfilesCheckHostNameAvailability`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.CheckHostNameAvailabilityInput{
	// ...
}


read, err := client.AFDProfilesCheckHostNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.AFDProfilesListResourceUsage`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.AFDProfilesListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.AFDProfilesListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProfilesClient.AFDProfilesUpgrade`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.ProfileUpgradeParameters{
	// ...
}


if err := client.AFDProfilesUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.AFDProfilesValidateSecret`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.ValidateSecretInput{
	// ...
}


read, err := client.AFDProfilesValidateSecret(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.CdnCanMigrateToAfd`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

if err := client.CdnCanMigrateToAfdThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProfilesClient.CdnMigrateToAfd`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

payload := profiles.CdnMigrationToAfdParameters{
	// ...
}


if err := client.CdnMigrateToAfdThenPoll(ctx, id, payload); err != nil {
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


### Example Usage: `ProfilesClient.LogAnalyticsGetLogAnalyticsLocations`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetLogAnalyticsLocations(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.LogAnalyticsGetLogAnalyticsMetrics`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetLogAnalyticsMetrics(ctx, id, profiles.DefaultLogAnalyticsGetLogAnalyticsMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.LogAnalyticsGetLogAnalyticsRankings`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetLogAnalyticsRankings(ctx, id, profiles.DefaultLogAnalyticsGetLogAnalyticsRankingsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.LogAnalyticsGetLogAnalyticsResources`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetLogAnalyticsResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.LogAnalyticsGetWafLogAnalyticsMetrics`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetWafLogAnalyticsMetrics(ctx, id, profiles.DefaultLogAnalyticsGetWafLogAnalyticsMetricsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.LogAnalyticsGetWafLogAnalyticsRankings`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

read, err := client.LogAnalyticsGetWafLogAnalyticsRankings(ctx, id, profiles.DefaultLogAnalyticsGetWafLogAnalyticsRankingsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.MigrationAbort`

```go
ctx := context.TODO()
id := profiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

if err := client.MigrationAbortThenPoll(ctx, id); err != nil {
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
