
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user` Documentation

The `user` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user"
```


### Client Initialization

```go
client := user.NewUserClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UserClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

payload := user.UserCreateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, user.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UserClient.Delete`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

read, err := client.Delete(ctx, id, user.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UserClient.Get`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UserClient.GetEntityTag`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UserClient.ListByService`

```go
ctx := context.TODO()
id := user.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, user.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, user.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UserClient.Update`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

payload := user.UserUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, user.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
