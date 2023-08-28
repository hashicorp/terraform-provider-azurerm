
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkgroups` Documentation

The `networkgroups` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkgroups"
```


### Client Initialization

```go
client := networkgroups.NewNetworkGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkgroups.NewNetworkGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "networkGroupValue")

payload := networkgroups.NetworkGroup{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, networkgroups.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkGroupsClient.Delete`

```go
ctx := context.TODO()
id := networkgroups.NewNetworkGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "networkGroupValue")

if err := client.DeleteThenPoll(ctx, id, networkgroups.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGroupsClient.Get`

```go
ctx := context.TODO()
id := networkgroups.NewNetworkGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "networkGroupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkGroupsClient.List`

```go
ctx := context.TODO()
id := networkgroups.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

// alternatively `client.List(ctx, id, networkgroups.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networkgroups.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
