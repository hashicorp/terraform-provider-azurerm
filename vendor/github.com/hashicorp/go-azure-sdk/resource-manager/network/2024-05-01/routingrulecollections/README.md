
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrulecollections` Documentation

The `routingrulecollections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrulecollections"
```


### Client Initialization

```go
client := routingrulecollections.NewRoutingRuleCollectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoutingRuleCollectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routingrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName")

payload := routingrulecollections.RoutingRuleCollection{
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


### Example Usage: `RoutingRuleCollectionsClient.Delete`

```go
ctx := context.TODO()
id := routingrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName")

if err := client.DeleteThenPoll(ctx, id, routingrulecollections.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `RoutingRuleCollectionsClient.Get`

```go
ctx := context.TODO()
id := routingrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoutingRuleCollectionsClient.List`

```go
ctx := context.TODO()
id := routingrulecollections.NewRoutingConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName")

// alternatively `client.List(ctx, id, routingrulecollections.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, routingrulecollections.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
