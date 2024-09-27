
## `github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targets` Documentation

The `targets` SDK allows for interaction with Azure Resource Manager `chaosstudio` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targets"
```


### Client Initialization

```go
client := targets.NewTargetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TargetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName")

payload := targets.Target{
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


### Example Usage: `TargetsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TargetsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TargetsClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id, targets.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, targets.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
