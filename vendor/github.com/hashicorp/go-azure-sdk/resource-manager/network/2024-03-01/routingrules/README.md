
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/routingrules` Documentation

The `routingrules` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/routingrules"
```


### Client Initialization

```go
client := routingrules.NewRoutingRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoutingRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routingrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName", "ruleName")

payload := routingrules.RoutingRule{
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


### Example Usage: `RoutingRulesClient.Delete`

```go
ctx := context.TODO()
id := routingrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName", "ruleName")

if err := client.DeleteThenPoll(ctx, id, routingrules.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `RoutingRulesClient.Get`

```go
ctx := context.TODO()
id := routingrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName", "ruleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoutingRulesClient.List`

```go
ctx := context.TODO()
id := routingrules.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName", "ruleCollectionName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
