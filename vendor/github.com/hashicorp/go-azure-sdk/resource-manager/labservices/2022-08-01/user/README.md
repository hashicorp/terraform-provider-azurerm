
## `github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user` Documentation

The `user` SDK allows for interaction with the Azure Resource Manager Service `labservices` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
```


### Client Initialization

```go
client := user.NewUserClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UserClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "userValue")

payload := user.User{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `UserClient.Delete`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "userValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UserClient.Get`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "userValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UserClient.Invite`

```go
ctx := context.TODO()
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "userValue")

payload := user.InviteBody{
	// ...
}


if err := client.InviteThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `UserClient.ListByLab`

```go
ctx := context.TODO()
id := user.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

// alternatively `client.ListByLab(ctx, id)` can be used to do batched pagination
items, err := client.ListByLabComplete(ctx, id)
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
id := user.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "userValue")

payload := user.UserUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
