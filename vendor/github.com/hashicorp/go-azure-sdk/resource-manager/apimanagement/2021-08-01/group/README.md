
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/group` Documentation

The `group` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/group"
```


### Client Initialization

```go
client := group.NewGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GroupClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

payload := group.GroupCreateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, group.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.Delete`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.Delete(ctx, id, group.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.Get`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.GetEntityTag`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.ListByService`

```go
ctx := context.TODO()
id := group.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, group.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, group.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GroupClient.Update`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

payload := group.GroupUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, group.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
