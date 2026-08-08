
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/roledefinitions` Documentation

The `roledefinitions` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/roledefinitions"
```


### Client Initialization

```go
client := roledefinitions.NewRoleDefinitionsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("https://endpoint-url.example.com", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := roledefinitions.RoleDefinitionCreateParameters{
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


### Example Usage: `RoleDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("https://endpoint-url.example.com", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleDefinitionsClient.Get`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("https://endpoint-url.example.com", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleDefinitionsClient.List`

```go
ctx := context.TODO()
id := roledefinitions.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id, roledefinitions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, roledefinitions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
