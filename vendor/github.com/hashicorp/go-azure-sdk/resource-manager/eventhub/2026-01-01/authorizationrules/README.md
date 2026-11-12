
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/authorizationrules` Documentation

The `authorizationrules` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/authorizationrules"
```


### Client Initialization

```go
client := authorizationrules.NewAuthorizationRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationRulesClient.NamespacesCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrules.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

payload := authorizationrules.AuthorizationRule{
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


### Example Usage: `AuthorizationRulesClient.NamespacesDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrules.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesClient.NamespacesGetAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrules.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesClient.NamespacesListAuthorizationRules`

```go
ctx := context.TODO()
id := authorizationrules.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.NamespacesListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AuthorizationRulesClient.NamespacesListKeys`

```go
ctx := context.TODO()
id := authorizationrules.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

read, err := client.NamespacesListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesClient.NamespacesRegenerateKeys`

```go
ctx := context.TODO()
id := authorizationrules.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "authorizationRuleName")

payload := authorizationrules.RegenerateAccessKeyParameters{
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
