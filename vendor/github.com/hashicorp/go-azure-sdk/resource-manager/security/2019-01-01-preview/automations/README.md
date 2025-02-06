
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations` Documentation

The `automations` SDK allows for interaction with Azure Resource Manager `security` (API Version `2019-01-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
```


### Client Initialization

```go
client := automations.NewAutomationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutomationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := automations.NewAutomationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationName")

payload := automations.Automation{
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


### Example Usage: `AutomationsClient.Delete`

```go
ctx := context.TODO()
id := automations.NewAutomationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutomationsClient.Get`

```go
ctx := context.TODO()
id := automations.NewAutomationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutomationsClient.List`

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


### Example Usage: `AutomationsClient.ListByResourceGroup`

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


### Example Usage: `AutomationsClient.Validate`

```go
ctx := context.TODO()
id := automations.NewAutomationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationName")

payload := automations.Automation{
	// ...
}


read, err := client.Validate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
