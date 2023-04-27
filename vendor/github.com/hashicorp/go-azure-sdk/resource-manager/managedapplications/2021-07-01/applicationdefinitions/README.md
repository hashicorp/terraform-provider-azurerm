
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions` Documentation

The `applicationdefinitions` SDK allows for interaction with the Azure Resource Manager Service `managedapplications` (API Version `2021-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
```


### Client Initialization

```go
client := applicationdefinitions.NewApplicationDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := applicationdefinitions.NewApplicationDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationDefinitionValue")

payload := applicationdefinitions.ApplicationDefinition{
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


### Example Usage: `ApplicationDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := applicationdefinitions.NewApplicationDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationDefinitionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationDefinitionsClient.Get`

```go
ctx := context.TODO()
id := applicationdefinitions.NewApplicationDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationDefinitionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationDefinitionsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := applicationdefinitions.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationDefinitionsClient.ListBySubscription`

```go
ctx := context.TODO()
id := applicationdefinitions.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationDefinitionsClient.Update`

```go
ctx := context.TODO()
id := applicationdefinitions.NewApplicationDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationDefinitionValue")

payload := applicationdefinitions.ApplicationDefinitionPatchable{
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
