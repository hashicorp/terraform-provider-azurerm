
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces` Documentation

The `authorizationrulesnamespaces` SDK allows for interaction with the Azure Resource Manager Service `eventhub` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
```


### Client Initialization

```go
client := authorizationrulesnamespaces.NewAuthorizationRulesNamespacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "authorizationRuleValue")

payload := authorizationrulesnamespaces.AuthorizationRule{
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


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "authorizationRuleValue")

read, err := client.NamespacesDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesGetAuthorizationRule`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "authorizationRuleValue")

read, err := client.NamespacesGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesListAuthorizationRules`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue")

// alternatively `client.NamespacesListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesListKeys`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "authorizationRuleValue")

read, err := client.NamespacesListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationRulesNamespacesClient.NamespacesRegenerateKeys`

```go
ctx := context.TODO()
id := authorizationrulesnamespaces.NewAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "authorizationRuleValue")

payload := authorizationrulesnamespaces.RegenerateAccessKeyParameters{
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
