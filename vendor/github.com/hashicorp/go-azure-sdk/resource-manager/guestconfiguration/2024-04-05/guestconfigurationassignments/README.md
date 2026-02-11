
## `github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments` Documentation

The `guestconfigurationassignments` SDK allows for interaction with Azure Resource Manager `guestconfiguration` (API Version `2024-04-05`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
```


### Client Initialization

```go
client := guestconfigurationassignments.NewGuestConfigurationAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GuestConfigurationAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "guestConfigurationAssignmentName")

payload := guestconfigurationassignments.GuestConfigurationAssignment{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "guestConfigurationAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.Get`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "guestConfigurationAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.GuestConfigurationAssignmentReportsGet`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewProviders2GuestConfigurationAssignmentReportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "guestConfigurationAssignmentName", "reportId")

read, err := client.GuestConfigurationAssignmentReportsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.GuestConfigurationAssignmentReportsList`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "guestConfigurationAssignmentName")

// alternatively `client.GuestConfigurationAssignmentReportsList(ctx, id)` can be used to do batched pagination
items, err := client.GuestConfigurationAssignmentReportsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.List`

```go
ctx := context.TODO()
id := guestconfigurationassignments.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GuestConfigurationAssignmentsClient.SubscriptionList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.SubscriptionList(ctx, id)` can be used to do batched pagination
items, err := client.SubscriptionListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
