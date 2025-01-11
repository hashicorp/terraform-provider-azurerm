
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules` Documentation

The `alertrules` SDK allows for interaction with Azure Resource Manager `securityinsights` (API Version `2022-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
```


### Client Initialization

```go
client := alertrules.NewAlertRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AlertRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "ruleId")

payload := alertrules.AlertRule{
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


### Example Usage: `AlertRulesClient.Delete`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "ruleId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRulesClient.Get`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "ruleId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRulesClient.List`

```go
ctx := context.TODO()
id := alertrules.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
