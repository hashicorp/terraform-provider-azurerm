
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule` Documentation

The `namespacesauthorizationrule` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
```


### Client Initialization

```go
client := namespacesauthorizationrule.NewNamespacesAuthorizationRuleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

payload := namespacesauthorizationrule.SBAuthorizationRule{
	// ...
}


read, err := client.NamespacesCreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesGetAuthorizationRule`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesListAuthorizationRules`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.NamespacesListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesListKeys`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesAuthorizationRuleClient.NamespacesRegenerateKeys`

```go
ctx := context.TODO()
id := namespacesauthorizationrule.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

payload := namespacesauthorizationrule.RegenerateAccessKeyParameters{
	// ...
}


read, err := client.NamespacesRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
