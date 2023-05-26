
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2021-09-01-preview/alertruletemplates` Documentation

The `alertruletemplates` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2021-09-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2021-09-01-preview/alertruletemplates"
```


### Client Initialization

```go
client := alertruletemplates.NewAlertRuleTemplatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AlertRuleTemplatesClient.AlertRuleTemplatesGet`

```go
ctx := context.TODO()
id := alertruletemplates.NewAlertRuleTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "alertRuleTemplateIdValue")

read, err := client.AlertRuleTemplatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertRuleTemplatesClient.AlertRuleTemplatesList`

```go
ctx := context.TODO()
id := alertruletemplates.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.AlertRuleTemplatesList(ctx, id)` can be used to do batched pagination
items, err := client.AlertRuleTemplatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
