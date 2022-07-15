
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps` Documentation

The `apps` SDK allows for interaction with the Azure Resource Manager Service `iotcentral` (API Version `2021-11-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
```


### Client Initialization

```go
client := apps.NewAppsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := apps.NewSubscriptionID()

payload := apps.OperationInputs{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppsClient.CheckSubdomainAvailability`

```go
ctx := context.TODO()
id := apps.NewSubscriptionID()

payload := apps.OperationInputs{
	// ...
}


read, err := client.CheckSubdomainAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apps.NewIotAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := apps.App{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppsClient.Delete`

```go
ctx := context.TODO()
id := apps.NewIotAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AppsClient.Get`

```go
ctx := context.TODO()
id := apps.NewIotAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := apps.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppsClient.ListBySubscription`

```go
ctx := context.TODO()
id := apps.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppsClient.ListTemplates`

```go
ctx := context.TODO()
id := apps.NewSubscriptionID()

// alternatively `client.ListTemplates(ctx, id)` can be used to do batched pagination
items, err := client.ListTemplatesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppsClient.Update`

```go
ctx := context.TODO()
id := apps.NewIotAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := apps.AppPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
