
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts` Documentation

The `dedicatedhosts` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
```


### Client Initialization

```go
client := dedicatedhosts.NewDedicatedHostsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DedicatedHostsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dedicatedhosts.NewHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupValue", "hostValue")

payload := dedicatedhosts.DedicatedHost{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DedicatedHostsClient.Delete`

```go
ctx := context.TODO()
id := dedicatedhosts.NewHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupValue", "hostValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DedicatedHostsClient.Get`

```go
ctx := context.TODO()
id := dedicatedhosts.NewHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupValue", "hostValue")

read, err := client.Get(ctx, id, dedicatedhosts.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DedicatedHostsClient.Update`

```go
ctx := context.TODO()
id := dedicatedhosts.NewHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostGroupValue", "hostValue")

payload := dedicatedhosts.DedicatedHostUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
