
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancekeys` Documentation

The `managedinstancekeys` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancekeys"
```


### Client Initialization

```go
client := managedinstancekeys.NewManagedInstanceKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedInstanceKeysClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedinstancekeys.NewManagedInstanceKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceValue", "keyValue")

payload := managedinstancekeys.ManagedInstanceKey{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedInstanceKeysClient.Delete`

```go
ctx := context.TODO()
id := managedinstancekeys.NewManagedInstanceKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceValue", "keyValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedInstanceKeysClient.Get`

```go
ctx := context.TODO()
id := managedinstancekeys.NewManagedInstanceKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceValue", "keyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedInstanceKeysClient.ListByInstance`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceValue")

// alternatively `client.ListByInstance(ctx, id, managedinstancekeys.DefaultListByInstanceOperationOptions())` can be used to do batched pagination
items, err := client.ListByInstanceComplete(ctx, id, managedinstancekeys.DefaultListByInstanceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
