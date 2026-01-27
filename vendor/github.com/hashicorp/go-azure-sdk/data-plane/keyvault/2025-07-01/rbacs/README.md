
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/rbacs` Documentation

The `rbacs` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/rbacs"
```


### Client Initialization

```go
client := rbacs.NewRbacsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `RbacsClient.RoleAssignmentsCreate`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleAssignmentID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := rbacs.RoleAssignmentCreateParameters{
	// ...
}


read, err := client.RoleAssignmentsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleAssignmentsDelete`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleAssignmentID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.RoleAssignmentsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleAssignmentsGet`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleAssignmentID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.RoleAssignmentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleAssignmentsListForScope`

```go
ctx := context.TODO()
id := rbacs.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.RoleAssignmentsListForScope(ctx, id, rbacs.DefaultRoleAssignmentsListForScopeOperationOptions())` can be used to do batched pagination
items, err := client.RoleAssignmentsListForScopeComplete(ctx, id, rbacs.DefaultRoleAssignmentsListForScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RbacsClient.RoleDefinitionsCreateOrUpdate`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleDefinitionID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := rbacs.RoleDefinitionCreateParameters{
	// ...
}


read, err := client.RoleDefinitionsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleDefinitionsDelete`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleDefinitionID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.RoleDefinitionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleDefinitionsGet`

```go
ctx := context.TODO()
id := rbacs.NewScopedRoleDefinitionID("https://endpoint_url", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.RoleDefinitionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RbacsClient.RoleDefinitionsList`

```go
ctx := context.TODO()
id := rbacs.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.RoleDefinitionsList(ctx, id, rbacs.DefaultRoleDefinitionsListOperationOptions())` can be used to do batched pagination
items, err := client.RoleDefinitionsListComplete(ctx, id, rbacs.DefaultRoleDefinitionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
