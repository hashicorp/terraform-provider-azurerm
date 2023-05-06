
## `github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules` Documentation

The `tagrules` SDK allows for interaction with the Azure Resource Manager Service `logz` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
```


### Client Initialization

```go
client := tagrules.NewTagRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "tagRuleValue")

payload := tagrules.MonitoringTagRules{
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


### Example Usage: `TagRulesClient.Delete`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "tagRuleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.Get`

```go
ctx := context.TODO()
id := tagrules.NewTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "tagRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.List`

```go
ctx := context.TODO()
id := tagrules.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagRulesClient.SubAccountTagRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := tagrules.NewAccountTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue", "tagRuleValue")

payload := tagrules.MonitoringTagRules{
	// ...
}


read, err := client.SubAccountTagRulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.SubAccountTagRulesDelete`

```go
ctx := context.TODO()
id := tagrules.NewAccountTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue", "tagRuleValue")

read, err := client.SubAccountTagRulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.SubAccountTagRulesGet`

```go
ctx := context.TODO()
id := tagrules.NewAccountTagRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue", "tagRuleValue")

read, err := client.SubAccountTagRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagRulesClient.SubAccountTagRulesList`

```go
ctx := context.TODO()
id := tagrules.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

// alternatively `client.SubAccountTagRulesList(ctx, id)` can be used to do batched pagination
items, err := client.SubAccountTagRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
