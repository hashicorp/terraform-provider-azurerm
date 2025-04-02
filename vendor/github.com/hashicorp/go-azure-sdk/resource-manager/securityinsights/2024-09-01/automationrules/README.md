
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules` Documentation

The `automationrules` SDK allows for interaction with Azure Resource Manager `securityinsights` (API Version `2024-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
```


### Client Initialization

```go
client := automationrules.NewAutomationRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutomationRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := automationrules.NewAutomationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "automationRuleId")

payload := automationrules.AutomationRule{
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


### Example Usage: `AutomationRulesClient.Delete`

```go
ctx := context.TODO()
id := automationrules.NewAutomationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "automationRuleId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutomationRulesClient.Get`

```go
ctx := context.TODO()
id := automationrules.NewAutomationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "automationRuleId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutomationRulesClient.List`

```go
ctx := context.TODO()
id := automationrules.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
