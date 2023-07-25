
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs` Documentation

The `labs` SDK allows for interaction with the Azure Resource Manager Service `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs"
```


### Client Initialization

```go
client := labs.NewLabsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LabsClient.ClaimAnyVM`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

if err := client.ClaimAnyVMThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.CreateEnvironment`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.LabVirtualMachineCreationParameter{
	// ...
}


if err := client.CreateEnvironmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.Lab{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.Delete`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.ExportResourceUsage`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.ExportResourceUsageParameters{
	// ...
}


if err := client.ExportResourceUsageThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.GenerateUploadUri`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.GenerateUploadUriParameter{
	// ...
}


read, err := client.GenerateUploadUri(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LabsClient.Get`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

read, err := client.Get(ctx, id, labs.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LabsClient.ImportVirtualMachine`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.ImportLabVirtualMachineRequest{
	// ...
}


if err := client.ImportVirtualMachineThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LabsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := labs.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, labs.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, labs.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LabsClient.ListBySubscription`

```go
ctx := context.TODO()
id := labs.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, labs.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, labs.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LabsClient.ListVhds`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

// alternatively `client.ListVhds(ctx, id)` can be used to do batched pagination
items, err := client.ListVhdsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LabsClient.Update`

```go
ctx := context.TODO()
id := labs.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

payload := labs.UpdateResource{
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
