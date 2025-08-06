
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/staticcidrs` Documentation

The `staticcidrs` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/staticcidrs"
```


### Client Initialization

```go
client := staticcidrs.NewStaticCidrsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StaticCidrsClient.Create`

```go
ctx := context.TODO()
id := staticcidrs.NewStaticCidrID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName", "staticCidrName")

payload := staticcidrs.StaticCidr{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticCidrsClient.Delete`

```go
ctx := context.TODO()
id := staticcidrs.NewStaticCidrID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName", "staticCidrName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StaticCidrsClient.Get`

```go
ctx := context.TODO()
id := staticcidrs.NewStaticCidrID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName", "staticCidrName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StaticCidrsClient.List`

```go
ctx := context.TODO()
id := staticcidrs.NewIPamPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "ipamPoolName")

// alternatively `client.List(ctx, id, staticcidrs.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, staticcidrs.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
