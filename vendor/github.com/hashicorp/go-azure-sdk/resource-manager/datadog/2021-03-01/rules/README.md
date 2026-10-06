
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules` Documentation

The `rules` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
```


### Client Initialization

```go
client := rules.NewRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RulesClient.TagRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := rules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

payload := rules.MonitoringTagRules{
	// ...
}


read, err := client.TagRulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RulesClient.TagRulesGet`

```go
ctx := context.TODO()
id := rules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "tagRuleName")

read, err := client.TagRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RulesClient.TagRulesList`

```go
ctx := context.TODO()
id := rules.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.TagRulesList(ctx, id)` can be used to do batched pagination
items, err := client.TagRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
