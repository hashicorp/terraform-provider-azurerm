
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/staticmembers` Documentation

The `staticmembers` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/staticmembers"
```


### Client Initialization

```go
client := staticmembers.NewStaticMembersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StaticMembersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := staticmembers.NewStaticMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "networkGroupName", "staticMemberName")

payload := staticmembers.StaticMember{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticMembersClient.Delete`

```go
ctx := context.TODO()
id := staticmembers.NewStaticMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "networkGroupName", "staticMemberName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticMembersClient.Get`

```go
ctx := context.TODO()
id := staticmembers.NewStaticMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "networkGroupName", "staticMemberName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticMembersClient.List`

```go
ctx := context.TODO()
id := staticmembers.NewNetworkGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "networkGroupName")

// alternatively `client.List(ctx, id, staticmembers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, staticmembers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
