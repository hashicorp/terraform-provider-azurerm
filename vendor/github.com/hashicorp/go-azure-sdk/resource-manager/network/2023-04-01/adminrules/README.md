
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/adminrules` Documentation

The `adminrules` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/adminrules"
```


### Client Initialization

```go
client := adminrules.NewAdminRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdminRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := adminrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue", "ruleCollectionValue", "ruleValue")

payload := adminrules.BaseAdminRule{
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


### Example Usage: `AdminRulesClient.Delete`

```go
ctx := context.TODO()
id := adminrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue", "ruleCollectionValue", "ruleValue")

if err := client.DeleteThenPoll(ctx, id, adminrules.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AdminRulesClient.Get`

```go
ctx := context.TODO()
id := adminrules.NewRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue", "ruleCollectionValue", "ruleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdminRulesClient.List`

```go
ctx := context.TODO()
id := adminrules.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue", "ruleCollectionValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
