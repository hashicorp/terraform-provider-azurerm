
## `github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/users` Documentation

The `users` SDK allows for interaction with Azure Resource Manager `mongocluster` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/users"
```


### Client Initialization

```go
client := users.NewUsersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UsersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := users.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName", "userName")

payload := users.User{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `UsersClient.Delete`

```go
ctx := context.TODO()
id := users.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName", "userName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UsersClient.Get`

```go
ctx := context.TODO()
id := users.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName", "userName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UsersClient.ListByMongoCluster`

```go
ctx := context.TODO()
id := users.NewMongoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mongoClusterName")

// alternatively `client.ListByMongoCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByMongoClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
