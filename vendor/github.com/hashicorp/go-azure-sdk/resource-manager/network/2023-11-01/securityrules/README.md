
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/securityrules` Documentation

The `securityrules` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/securityrules"
```


### Client Initialization

```go
client := securityrules.NewSecurityRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securityrules.NewSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName", "securityRuleName")

payload := securityrules.SecurityRule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityRulesClient.DefaultSecurityRulesGet`

```go
ctx := context.TODO()
id := securityrules.NewDefaultSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName", "defaultSecurityRuleName")

read, err := client.DefaultSecurityRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityRulesClient.DefaultSecurityRulesList`

```go
ctx := context.TODO()
id := securityrules.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

// alternatively `client.DefaultSecurityRulesList(ctx, id)` can be used to do batched pagination
items, err := client.DefaultSecurityRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecurityRulesClient.Delete`

```go
ctx := context.TODO()
id := securityrules.NewSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName", "securityRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityRulesClient.Get`

```go
ctx := context.TODO()
id := securityrules.NewSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName", "securityRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityRulesClient.List`

```go
ctx := context.TODO()
id := securityrules.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
