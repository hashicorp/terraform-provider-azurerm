
## `github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments` Documentation

The `environments` SDK allows for interaction with the Azure Resource Manager Service `timeseriesinsights` (API Version `2020-05-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
```


### Client Initialization

```go
client := environments.NewEnvironmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnvironmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := environments.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

payload := environments.EnvironmentCreateOrUpdateParameters{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EnvironmentsClient.Delete`

```go
ctx := context.TODO()
id := environments.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentsClient.Get`

```go
ctx := context.TODO()
id := environments.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

read, err := client.Get(ctx, id, environments.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := environments.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentsClient.ListBySubscription`

```go
ctx := context.TODO()
id := environments.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentsClient.Update`

```go
ctx := context.TODO()
id := environments.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

payload := environments.EnvironmentUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
