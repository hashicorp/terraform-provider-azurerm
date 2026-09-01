
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeteraccessrules` Documentation

The `networksecurityperimeteraccessrules` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeteraccessrules"
```


### Client Initialization

```go
client := networksecurityperimeteraccessrules.NewNetworkSecurityPerimeterAccessRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterAccessRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecurityperimeteraccessrules.NewAccessRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName", "accessRuleName")

payload := networksecurityperimeteraccessrules.NspAccessRule{
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


### Example Usage: `NetworkSecurityPerimeterAccessRulesClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeteraccessrules.NewAccessRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName", "accessRuleName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterAccessRulesClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeteraccessrules.NewAccessRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName", "accessRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterAccessRulesClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeteraccessrules.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName")

// alternatively `client.List(ctx, id, networksecurityperimeteraccessrules.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeteraccessrules.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkSecurityPerimeterAccessRulesClient.Reconcile`

```go
ctx := context.TODO()
id := networksecurityperimeteraccessrules.NewAccessRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName", "accessRuleName")
var payload interface{}

read, err := client.Reconcile(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
