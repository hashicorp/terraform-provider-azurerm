
## `github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers` Documentation

The `servers` SDK allows for interaction with the Azure Resource Manager Service `mariadb` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
```


### Client Initialization

```go
client := servers.NewServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServersClient.Create`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.ServerForCreate{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Delete`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Get`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.Update`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.ServerUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
