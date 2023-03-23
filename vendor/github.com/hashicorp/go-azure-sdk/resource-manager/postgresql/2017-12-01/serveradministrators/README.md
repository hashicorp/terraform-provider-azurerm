
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serveradministrators` Documentation

The `serveradministrators` SDK allows for interaction with the Azure Resource Manager Service `postgresql` (API Version `2017-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serveradministrators"
```


### Client Initialization

```go
client := serveradministrators.NewServerAdministratorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerAdministratorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := serveradministrators.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := serveradministrators.ServerAdministratorResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerAdministratorsClient.Delete`

```go
ctx := context.TODO()
id := serveradministrators.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServerAdministratorsClient.Get`

```go
ctx := context.TODO()
id := serveradministrators.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerAdministratorsClient.List`

```go
ctx := context.TODO()
id := serveradministrators.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
