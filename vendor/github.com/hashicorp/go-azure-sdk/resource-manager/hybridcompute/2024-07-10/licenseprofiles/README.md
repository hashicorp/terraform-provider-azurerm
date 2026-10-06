
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenseprofiles` Documentation

The `licenseprofiles` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2024-07-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenseprofiles"
```


### Client Initialization

```go
client := licenseprofiles.NewLicenseProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LicenseProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := licenseprofiles.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

payload := licenseprofiles.LicenseProfile{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LicenseProfilesClient.Delete`

```go
ctx := context.TODO()
id := licenseprofiles.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LicenseProfilesClient.Get`

```go
ctx := context.TODO()
id := licenseprofiles.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LicenseProfilesClient.List`

```go
ctx := context.TODO()
id := licenseprofiles.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LicenseProfilesClient.Update`

```go
ctx := context.TODO()
id := licenseprofiles.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

payload := licenseprofiles.LicenseProfileUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
