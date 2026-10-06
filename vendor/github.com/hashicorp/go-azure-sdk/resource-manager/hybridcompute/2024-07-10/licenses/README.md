
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenses` Documentation

The `licenses` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2024-07-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenses"
```


### Client Initialization

```go
client := licenses.NewLicensesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LicensesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := licenses.NewLicenseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "licenseName")

payload := licenses.License{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LicensesClient.Delete`

```go
ctx := context.TODO()
id := licenses.NewLicenseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "licenseName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LicensesClient.Get`

```go
ctx := context.TODO()
id := licenses.NewLicenseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "licenseName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LicensesClient.ListByResourceGroup`

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


### Example Usage: `LicensesClient.ListBySubscription`

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


### Example Usage: `LicensesClient.Update`

```go
ctx := context.TODO()
id := licenses.NewLicenseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "licenseName")

payload := licenses.LicenseUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LicensesClient.ValidateLicense`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := licenses.License{
	// ...
}


if err := client.ValidateLicenseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
