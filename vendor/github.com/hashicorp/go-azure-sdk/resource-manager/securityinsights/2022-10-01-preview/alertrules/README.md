
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules` Documentation

The `alertrules` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2022-10-01-preview`).

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


### Example Usage: `AlertRulesClient.AlertRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "ruleIdValue")

payload := alertrules.AlertRule{
	// ...
}


read, err := client.AlertRulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRulesClient.AlertRulesDelete`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "ruleIdValue")

read, err := client.AlertRulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRulesClient.AlertRulesGet`

```go
ctx := context.TODO()
id := alertrules.NewAlertRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "ruleIdValue")

read, err := client.AlertRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRulesClient.AlertRulesList`

```go
ctx := context.TODO()
id := alertrules.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.AlertRulesList(ctx, id)` can be used to do batched pagination
items, err := client.AlertRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
