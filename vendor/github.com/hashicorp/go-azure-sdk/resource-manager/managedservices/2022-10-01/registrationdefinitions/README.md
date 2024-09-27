
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions` Documentation

The `registrationdefinitions` SDK allows for interaction with Azure Resource Manager `managedservices` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
```


### Client Initialization

```go
client := registrationdefinitions.NewRegistrationDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistrationDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := registrationdefinitions.NewScopedRegistrationDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationDefinitionId")

payload := registrationdefinitions.RegistrationDefinition{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistrationDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := registrationdefinitions.NewScopedRegistrationDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationDefinitionId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistrationDefinitionsClient.Get`

```go
ctx := context.TODO()
id := registrationdefinitions.NewScopedRegistrationDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationDefinitionId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistrationDefinitionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id, registrationdefinitions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, registrationdefinitions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
