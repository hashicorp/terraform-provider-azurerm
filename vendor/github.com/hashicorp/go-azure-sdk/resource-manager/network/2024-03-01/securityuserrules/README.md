
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/securityuserrules` Documentation

The `securityuserrules` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/securityuserrules"
```


### Client Initialization

```go
client := securityuserrules.NewSecurityUserRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityUserRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securityuserrules.NewRuleCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName", "ruleName")

payload := securityuserrules.SecurityUserRule{
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


### Example Usage: `SecurityUserRulesClient.Delete`

```go
ctx := context.TODO()
id := securityuserrules.NewRuleCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName", "ruleName")

if err := client.DeleteThenPoll(ctx, id, securityuserrules.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityUserRulesClient.Get`

```go
ctx := context.TODO()
id := securityuserrules.NewRuleCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName", "ruleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityUserRulesClient.List`

```go
ctx := context.TODO()
id := securityuserrules.NewSecurityUserConfigurationRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
