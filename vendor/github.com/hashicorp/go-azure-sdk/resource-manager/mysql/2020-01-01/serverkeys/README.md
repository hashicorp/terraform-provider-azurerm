
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverkeys` Documentation

The `serverkeys` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2020-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverkeys"
```


### Client Initialization

```go
client := serverkeys.NewServerKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerKeysClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := serverkeys.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "keyValue")

payload := serverkeys.ServerKey{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerKeysClient.Delete`

```go
ctx := context.TODO()
id := serverkeys.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "keyValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServerKeysClient.Get`

```go
ctx := context.TODO()
id := serverkeys.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "keyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerKeysClient.List`

```go
ctx := context.TODO()
id := serverkeys.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
