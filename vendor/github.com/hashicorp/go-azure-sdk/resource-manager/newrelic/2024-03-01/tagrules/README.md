
## `github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules` Documentation

The `tagrules` SDK allows for interaction with Azure Resource Manager `newrelic` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
```


### Client Initialization

```go
client := tagrules.NewTagRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

payload := tagrules.TagRule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TagRulesClient.Delete`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TagRulesClient.Get`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.ListByNewRelicMonitorResource`

```go
ctx := context.TODO()
id := tagrules.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.ListByNewRelicMonitorResource(ctx, id)` can be used to do batched pagination
items, err := client.ListByNewRelicMonitorResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagRulesClient.Update`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

payload := tagrules.TagRuleUpdate{
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
