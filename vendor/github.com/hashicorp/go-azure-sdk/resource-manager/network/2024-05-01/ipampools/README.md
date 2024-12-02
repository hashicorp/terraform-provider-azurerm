
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools` Documentation

The `ipampools` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/ipampools"
```


### Client Initialization

```go
client := ipampools.NewIPamPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IPamPoolsClient.Create`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

payload := ipampools.IPamPool{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IPamPoolsClient.Delete`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IPamPoolsClient.Get`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IPamPoolsClient.GetPoolUsage`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

read, err := client.GetPoolUsage(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IPamPoolsClient.List`

```go
ctx := context.TODO()
id := ipampools.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, ipampools.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, ipampools.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IPamPoolsClient.ListAssociatedResources`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

// alternatively `client.ListAssociatedResources(ctx, id)` can be used to do batched pagination
items, err := client.ListAssociatedResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IPamPoolsClient.Update`

```go
ctx := context.TODO()
id := ipampools.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

payload := ipampools.IPamPoolUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
