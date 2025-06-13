
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles` Documentation

The `autoupgradeprofiles` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
```


### Client Initialization

```go
client := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutoUpgradeProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autoupgradeprofiles.NewAutoUpgradeProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "autoUpgradeProfileName")

payload := autoupgradeprofiles.AutoUpgradeProfile{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, autoupgradeprofiles.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AutoUpgradeProfilesClient.Delete`

```go
ctx := context.TODO()
id := autoupgradeprofiles.NewAutoUpgradeProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "autoUpgradeProfileName")

if err := client.DeleteThenPoll(ctx, id, autoupgradeprofiles.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AutoUpgradeProfilesClient.Get`

```go
ctx := context.TODO()
id := autoupgradeprofiles.NewAutoUpgradeProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "autoUpgradeProfileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutoUpgradeProfilesClient.ListByFleet`

```go
ctx := context.TODO()
id := autoupgradeprofiles.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

// alternatively `client.ListByFleet(ctx, id)` can be used to do batched pagination
items, err := client.ListByFleetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
