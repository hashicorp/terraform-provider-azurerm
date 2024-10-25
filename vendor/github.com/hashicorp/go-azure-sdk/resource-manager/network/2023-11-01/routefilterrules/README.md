
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routefilterrules` Documentation

The `routefilterrules` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routefilterrules"
```


### Client Initialization

```go
client := routefilterrules.NewRouteFilterRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RouteFilterRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routefilterrules.NewRouteFilterRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterName", "routeFilterRuleName")

payload := routefilterrules.RouteFilterRule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RouteFilterRulesClient.Delete`

```go
ctx := context.TODO()
id := routefilterrules.NewRouteFilterRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterName", "routeFilterRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RouteFilterRulesClient.Get`

```go
ctx := context.TODO()
id := routefilterrules.NewRouteFilterRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterName", "routeFilterRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RouteFilterRulesClient.ListByRouteFilter`

```go
ctx := context.TODO()
id := routefilterrules.NewRouteFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeFilterName")

// alternatively `client.ListByRouteFilter(ctx, id)` can be used to do batched pagination
items, err := client.ListByRouteFilterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
