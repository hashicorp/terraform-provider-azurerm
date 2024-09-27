
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments` Documentation

The `integrationserviceenvironments` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments"
```


### Client Initialization

```go
client := integrationserviceenvironments.NewIntegrationServiceEnvironmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationServiceEnvironmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationserviceenvironments.NewIntegrationServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationServiceEnvironmentName")

payload := integrationserviceenvironments.IntegrationServiceEnvironment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationServiceEnvironmentsClient.Delete`

```go
ctx := context.TODO()
id := integrationserviceenvironments.NewIntegrationServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationServiceEnvironmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationServiceEnvironmentsClient.Get`

```go
ctx := context.TODO()
id := integrationserviceenvironments.NewIntegrationServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationServiceEnvironmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationServiceEnvironmentsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, integrationserviceenvironments.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, integrationserviceenvironments.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationServiceEnvironmentsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, integrationserviceenvironments.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, integrationserviceenvironments.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationServiceEnvironmentsClient.Update`

```go
ctx := context.TODO()
id := integrationserviceenvironments.NewIntegrationServiceEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationServiceEnvironmentName")

payload := integrationserviceenvironments.IntegrationServiceEnvironment{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
