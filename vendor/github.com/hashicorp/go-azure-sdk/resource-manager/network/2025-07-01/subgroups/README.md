
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subgroups` Documentation

The `subgroups` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subgroups"
```


### Client Initialization

```go
client := subgroups.NewSubgroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubgroupsClient.Get`

```go
ctx := context.TODO()
id := subgroups.NewSubgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "interconnectGroupName", "subgroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubgroupsClient.List`

```go
ctx := context.TODO()
id := subgroups.NewInterconnectGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "interconnectGroupName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
