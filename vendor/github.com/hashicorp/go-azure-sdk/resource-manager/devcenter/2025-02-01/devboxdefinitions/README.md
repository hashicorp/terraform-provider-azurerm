
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions` Documentation

The `devboxdefinitions` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
```


### Client Initialization

```go
client := devboxdefinitions.NewDevBoxDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DevBoxDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevCenterDevBoxDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "devBoxDefinitionName")

payload := devboxdefinitions.DevBoxDefinition{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DevBoxDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevCenterDevBoxDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "devBoxDefinitionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevBoxDefinitionsClient.Get`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevCenterDevBoxDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "devBoxDefinitionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevBoxDefinitionsClient.GetByProject`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevBoxDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "devBoxDefinitionName")

read, err := client.GetByProject(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevBoxDefinitionsClient.ListByDevCenter`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.ListByDevCenter(ctx, id)` can be used to do batched pagination
items, err := client.ListByDevCenterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevBoxDefinitionsClient.ListByProject`

```go
ctx := context.TODO()
id := devboxdefinitions.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.ListByProject(ctx, id)` can be used to do batched pagination
items, err := client.ListByProjectComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevBoxDefinitionsClient.Update`

```go
ctx := context.TODO()
id := devboxdefinitions.NewDevCenterDevBoxDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "devBoxDefinitionName")

payload := devboxdefinitions.DevBoxDefinitionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
