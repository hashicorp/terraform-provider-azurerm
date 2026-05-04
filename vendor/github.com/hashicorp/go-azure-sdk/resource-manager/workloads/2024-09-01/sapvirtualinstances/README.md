
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances` Documentation

The `sapvirtualinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2024-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
```


### Client Initialization

```go
client := sapvirtualinstances.NewSAPVirtualInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPVirtualInstancesClient.Create`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

payload := sapvirtualinstances.SAPVirtualInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPVirtualInstancesClient.Delete`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPVirtualInstancesClient.Get`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPVirtualInstancesClient.InvokeAvailabilityZoneDetails`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapvirtualinstances.SAPAvailabilityZoneDetailsRequest{
	// ...
}


read, err := client.InvokeAvailabilityZoneDetails(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPVirtualInstancesClient.InvokeDiskConfigurations`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapvirtualinstances.SAPDiskConfigurationsRequest{
	// ...
}


read, err := client.InvokeDiskConfigurations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPVirtualInstancesClient.InvokeSapSupportedSku`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapvirtualinstances.SAPSupportedSkusRequest{
	// ...
}


read, err := client.InvokeSapSupportedSku(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPVirtualInstancesClient.InvokeSizingRecommendations`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapvirtualinstances.SAPSizingRecommendationRequest{
	// ...
}


read, err := client.InvokeSizingRecommendations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPVirtualInstancesClient.ListByResourceGroup`

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


### Example Usage: `SAPVirtualInstancesClient.ListBySubscription`

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


### Example Usage: `SAPVirtualInstancesClient.Start`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

payload := sapvirtualinstances.StartRequest{
	// ...
}


if err := client.StartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPVirtualInstancesClient.Stop`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

payload := sapvirtualinstances.StopRequest{
	// ...
}


if err := client.StopThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPVirtualInstancesClient.Update`

```go
ctx := context.TODO()
id := sapvirtualinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

payload := sapvirtualinstances.UpdateSAPVirtualInstanceRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
